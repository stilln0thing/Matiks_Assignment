package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/stilln0thing/matiks_leaderboard/internal/config"
	"github.com/stilln0thing/matiks_leaderboard/internal/database"
	"github.com/stilln0thing/matiks_leaderboard/internal/handler"
	"github.com/stilln0thing/matiks_leaderboard/internal/repository"
	"github.com/stilln0thing/matiks_leaderboard/internal/service"
	"github.com/stilln0thing/matiks_leaderboard/internal/simulator"
	"github.com/stilln0thing/matiks_leaderboard/internal/worker"
)

func main() {
	// 1. Load config
	cfg := config.Load()
	// 2. Connect to PostgreSQL
	log.Println("Connecting to PostgreSQL...")
	db, err := database.NewPostgres(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("PostgreSQL failed: %v", err)
	}
	defer db.Close()
	log.Println("PostgreSQL connected")
	// 3. Connect to Redis
	log.Println("Connecting to Redis...")
	redisClient, err := database.NewRedis(cfg.RedisAddr)
	if err != nil {
		log.Fatalf("Redis failed: %v", err)
	}
	defer redisClient.Close()
	log.Println("Redis connected")
	// 4. Initialize repositories
	userRepo := repository.NewUserRepository(db)
	cacheRepo := repository.NewCacheRepository(redisClient)
	// 5. Initialize DB writer worker
	dbWriter := worker.NewDBWriter(userRepo, 10000, 500, 250*time.Millisecond)
	// 6. Initialize service
	leaderboardService := service.NewLeaderboardService(userRepo, cacheRepo, dbWriter.Queue())
	// 7. Warm cache from DB
	ctx := context.Background()
	if err := leaderboardService.WarmCache(ctx); err != nil {
		log.Printf("Warning: Failed to warm cache: %v", err)
	}
	// 8. Initialize handler
	leaderboardHandler := handler.NewLeaderboardHandler(leaderboardService)
	// 9. Initialize simulator (optional - for demo)
	scoreUpdater := simulator.NewScoreUpdater(userRepo, leaderboardService, 1*time.Second, 10)
	// 10. Setup Gin router
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(gin.Logger())
	// CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy"})
	})
	// API routes
	api := r.Group("/api")
	leaderboardHandler.RegisterRoutes(api)
	// Create shutdown context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	// Start background workers
	go dbWriter.Start(ctx)
	go scoreUpdater.Start(ctx)
	// Start HTTP server
	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: r,
	}
	go func() {
		log.Printf("Server starting on port %s", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()
	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down...")
	cancel()
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()
	srv.Shutdown(shutdownCtx)
	log.Println("Server exited")
}

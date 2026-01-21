package simulator

import (
	"context"
	"log"
	"math/rand"
	"time"

	"github.com/stilln0thing/matiks_leaderboard/internal/repository"
	"github.com/stilln0thing/matiks_leaderboard/internal/service"
)

type ScoreUpdater struct {
	userRepo    *repository.UserRepository
	leaderboard *service.LeaderboardService
	interval    time.Duration
	batchSize   int
}

func NewScoreUpdater(
	userRepo *repository.UserRepository,
	leaderboard *service.LeaderboardService,
	interval time.Duration,
	batchSize int,
) *ScoreUpdater {
	return &ScoreUpdater{
		userRepo:    userRepo,
		leaderboard: leaderboard,
		interval:    interval,
		batchSize:   batchSize,
	}
}
func (s *ScoreUpdater) Start(ctx context.Context) {
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()
	log.Printf("[Simulator] Started - updating %d users every %v", s.batchSize, s.interval)
	for {
		select {
		case <-ticker.C:
			s.updateRandomUsers(ctx)
		case <-ctx.Done():
			log.Println("[Simulator] Stopped")
			return
		}
	}
}
func (s *ScoreUpdater) updateRandomUsers(ctx context.Context) {
	userIDs, err := s.userRepo.GetRandomUserIDs(ctx, s.batchSize)
	if err != nil {
		return
	}
	for _, userID := range userIDs {
		user, err := s.userRepo.GetUserByID(ctx, userID)
		if err != nil {
			continue
		}
		// Random change: -100 to +100
		delta := rand.Intn(201) - 100
		newRating := user.Rating + delta
		// Clamp to valid range
		if newRating < 100 {
			newRating = 100
		}
		if newRating > 5000 {
			newRating = 5000
		}
		s.leaderboard.UpdateRating(ctx, userID, newRating)
	}
}

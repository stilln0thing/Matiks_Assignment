package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/stilln0thing/matiks_leaderboard/internal/service"
)

type LeaderboardHandler struct {
	service *service.LeaderboardService
}

func NewLeaderboardHandler(service *service.LeaderboardService) *LeaderboardHandler {
	return &LeaderboardHandler{service: service}
}

func (h *LeaderboardHandler) RegisterRoutes(r *gin.RouterGroup) {
	r.GET("/leaderboard", h.GetLeaderboard)
	r.GET("/search", h.SearchUsers)
	r.GET("/user/:id/rank", h.GetUserRank)
	r.POST("/rating", h.UpdateRating)
}

// GET /api/leaderboard?limit=50&offset=0
func (h *LeaderboardHandler) GetLeaderboard(c *gin.Context) {
	limit, _ := strconv.ParseInt(c.DefaultQuery("limit", "50"), 10, 64)
	offset, _ := strconv.ParseInt(c.DefaultQuery("offset", "0"), 10, 64)
	// Clamp limit
	if limit > 100 {
		limit = 100
	}
	if limit < 1 {
		limit = 50
	}
	users, total, err := h.service.GetLeaderboard(c.Request.Context(), limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"users":  users,
		"total":  total,
		"limit":  limit,
		"offset": offset,
	})
}

// GET /api/search?q=john
func (h *LeaderboardHandler) SearchUsers(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "query 'q' is required"})
		return
	}
	users, err := h.service.SearchUsers(c.Request.Context(), query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"users": users,
		"count": len(users),
	})
}

// GET /api/user/:id/rank
func (h *LeaderboardHandler) GetUserRank(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}
	user, err := h.service.GetUserRank(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}
	c.JSON(http.StatusOK, user)
}

type UpdateRatingRequest struct {
	UserID int64 `json:"user_id" binding:"required"`
	Rating int   `json:"rating" binding:"required,min=100,max=5000"`
}

// POST /api/rating
func (h *LeaderboardHandler) UpdateRating(c *gin.Context) {
	var req UpdateRatingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := h.service.UpdateRating(c.Request.Context(), req.UserID, req.Rating)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "updated"})
}

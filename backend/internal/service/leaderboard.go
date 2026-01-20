package service

import (
	"context"
	"log"
	"time"

	"github.com/stilln0thing/matiks_leaderboard/internal/models"
	"github.com/stilln0thing/matiks_leaderboard/internal/repository"
)

type LeaderboardService struct {
	userRepo    *repository.UserRepository
	cacheRepo   *repository.CacheRepository
	updateQueue chan<- models.RatingUpdate // Write-only channel
}

func NewLeaderboardService(
	userRepo *repository.UserRepository,
	cacheRepo *repository.CacheRepository,
	updateQueue chan<- models.RatingUpdate,
) *LeaderboardService {
	return &LeaderboardService{
		userRepo:    userRepo,
		cacheRepo:   cacheRepo,
		updateQueue: updateQueue,
	}
}

// GetLeaderboard - Returns paginated leaderboard from Redis
func (s *LeaderboardService) GetLeaderboard(ctx context.Context, limit, offset int64) ([]models.RankedUser, int64, error) {
	users, err := s.cacheRepo.GetLeaderboard(ctx, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	total, err := s.cacheRepo.GetTotalUsers(ctx)
	if err != nil {
		return nil, 0, err
	}
	return users, total, nil
}

// SearchUsers - Search + get live ranks
func (s *LeaderboardService) SearchUsers(ctx context.Context, query string) ([]models.RankedUser, error) {
	// Search in PostgreSQL
	users, err := s.userRepo.SearchUsers(ctx, query)
	if err != nil {
		return nil, err
	}
	// Get live ranks from Redis
	rankedUsers := make([]models.RankedUser, 0, len(users))
	for _, u := range users {
		rank, rating, err := s.cacheRepo.GetRank(ctx, u.ID)
		if err != nil {
			// Fallback to DB rating
			rank = 0
			rating = u.Rating
		}
		rankedUsers = append(rankedUsers, models.RankedUser{
			Rank:     rank,
			ID:       u.ID,
			Username: u.Username,
			Rating:   rating,
		})
	}
	return rankedUsers, nil
}

// GetUserRank - Get single user's rank
func (s *LeaderboardService) GetUserRank(ctx context.Context, userID int64) (*models.RankedUser, error) {
	user, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	rank, rating, err := s.cacheRepo.GetRank(ctx, userID)
	if err != nil {
		rating = user.Rating
		rank = 0
	}
	return &models.RankedUser{
		Rank:     rank,
		ID:       user.ID,
		Username: user.Username,
		Rating:   rating,
	}, nil
}

// UpdateRating - Redis first, then async DB!
func (s *LeaderboardService) UpdateRating(ctx context.Context, userID int64, newRating int) error {
	version := time.Now().UnixNano() // Use timestamp as version
	// 1. Update Redis FIRST (fast path, instant feedback)
	err := s.cacheRepo.UpdateRating(ctx, userID, newRating, version)
	if err != nil {
		return err
	}
	// 2. Queue for async DB write (non-blocking)
	select {
	case s.updateQueue <- models.RatingUpdate{
		UserID:  userID,
		Rating:  newRating,
		Version: version,
	}:
	default:
		log.Printf("[Service] Warning: update queue full!")
	}
	return nil
}

// WarmCache - Load all users from DB to Redis at startup
func (s *LeaderboardService) WarmCache(ctx context.Context) error {
	log.Println("[Service] Warming cache...")
	start := time.Now()
	users, err := s.userRepo.GetAllUsers(ctx)
	if err != nil {
		return err
	}
	err = s.cacheRepo.WarmCache(ctx, users)
	if err != nil {
		return err
	}
	log.Printf("[Service] Cache warmed with %d users in %v", len(users), time.Since(start))
	return nil
}

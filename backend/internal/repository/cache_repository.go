package repository

import (
	"context"
	"fmt"
	"strconv"

	"github.com/redis/go-redis/v9"
	"github.com/stilln0thing/matiks_leaderboard/internal/models"
)

// Lua script for ATOMIC rating update with version check

const updateRatingScript = `
local oldVersion = redis.call('HGET', KEYS[2], 'version')
if oldVersion and tonumber(ARGV[3]) <= tonumber(oldVersion) then
    return 0  -- Stale update, ignore!
end
redis.call('ZADD', KEYS[1], ARGV[2], ARGV[1])
redis.call('HSET', KEYS[2], 'rating', ARGV[2], 'version', ARGV[3])
return 1
`
const (
	LeaderboardKey = "leaderboard:zset" // Sorted set for rankings
	UserHashPrefix = "user:hash:"       // Hash for user metadata
)

type CacheRepository struct {
	client       *redis.Client
	updateScript *redis.Script
}

func NewCacheRepository(client *redis.Client) *CacheRepository {
	return &CacheRepository{
		client:       client,
		updateScript: redis.NewScript(updateRatingScript),
	}
}

// UpdateRating - Atomic update using Lua script
func (r *CacheRepository) UpdateRating(ctx context.Context, userID int64, rating int, version int64) error {
	userIDStr := strconv.FormatInt(userID, 10)
	hashKey := UserHashPrefix + userIDStr
	_, err := r.updateScript.Run(ctx, r.client,
		[]string{LeaderboardKey, hashKey},
		userIDStr, rating, version,
	).Int()
	if err != nil {
		return err
	}
	// result == 0 means stale update was ignored (version conflict)
	return nil
}

// GetRank - Get user's rank using ZCOUNT
// This is O(log N) which scales well!
func (r *CacheRepository) GetRank(ctx context.Context, userID int64) (int64, int, error) {
	userIDStr := strconv.FormatInt(userID, 10)
	hashKey := UserHashPrefix + userIDStr
	// Get user's current rating
	ratingStr, err := r.client.HGet(ctx, hashKey, "rating").Result()
	if err == redis.Nil {
		return 0, 0, fmt.Errorf("user not found in cache")
	}
	if err != nil {
		return 0, 0, err
	}
	rating, _ := strconv.Atoi(ratingStr)
	// Count users with HIGHER rating
	// rank = count of users with rating > this user's rating + 1
	count, err := r.client.ZCount(ctx, LeaderboardKey,
		fmt.Sprintf("(%d", rating), "+inf",
	).Result()
	if err != nil {
		return 0, 0, err
	}
	return count + 1, rating, nil // +1 for 1-indexed rank
}

// GetLeaderboard - Paginated leaderboard with tie-aware ranking
func (r *CacheRepository) GetLeaderboard(ctx context.Context, limit, offset int64) ([]models.RankedUser, error) {
	// Get user IDs and scores from sorted set (descending order)
	results, err := r.client.ZRevRangeWithScores(ctx, LeaderboardKey, offset, offset+limit-1).Result()
	if err != nil {
		return nil, err
	}
	if len(results) == 0 {
		return []models.RankedUser{}, nil
	}
	// Pipeline to get usernames efficiently
	pipe := r.client.Pipeline()
	cmds := make([]*redis.MapStringStringCmd, len(results))
	for i, z := range results {
		userID := z.Member.(string)
		cmds[i] = pipe.HGetAll(ctx, UserHashPrefix+userID)
	}
	_, err = pipe.Exec(ctx)
	if err != nil {
		return nil, err
	}
	// Build result with tie-aware ranking
	users := make([]models.RankedUser, 0, len(results))
	var currentRank int64 = offset + 1
	var prevRating float64 = -1
	for i, z := range results {
		userID, _ := strconv.ParseInt(z.Member.(string), 10, 64)
		rating := int(z.Score)
		data, _ := cmds[i].Result()
		username := data["username"]
		// Handle ties - same rating = same rank
		if float64(rating) != prevRating {
			currentRank = offset + int64(i) + 1
			prevRating = float64(rating)
		}
		users = append(users, models.RankedUser{
			Rank:     currentRank,
			ID:       userID,
			Username: username,
			Rating:   rating,
		})
	}
	return users, nil
}

// SetUser - Store user in Redis (used during cache warming)
func (r *CacheRepository) SetUser(ctx context.Context, user models.User) error {
	userIDStr := strconv.FormatInt(user.ID, 10)
	hashKey := UserHashPrefix + userIDStr
	pipe := r.client.Pipeline()

	// Add to sorted set
	pipe.ZAdd(ctx, LeaderboardKey, redis.Z{
		Score:  float64(user.Rating),
		Member: userIDStr,
	})

	// Store metadata in hash
	pipe.HSet(ctx, hashKey, map[string]interface{}{
		"username": user.Username,
		"rating":   user.Rating,
		"version":  user.Version,
	})
	_, err := pipe.Exec(ctx)
	return err
}

// WarmCache - Load all users from slice into Redis
func (r *CacheRepository) WarmCache(ctx context.Context, users []models.User) error {
	pipe := r.client.Pipeline()
	for _, u := range users {
		userIDStr := strconv.FormatInt(u.ID, 10)
		hashKey := UserHashPrefix + userIDStr
		pipe.ZAdd(ctx, LeaderboardKey, redis.Z{
			Score:  float64(u.Rating),
			Member: userIDStr,
		})
		pipe.HSet(ctx, hashKey, map[string]interface{}{
			"username": u.Username,
			"rating":   u.Rating,
			"version":  u.Version,
		})
	}
	_, err := pipe.Exec(ctx)
	return err
}

// GetTotalUsers - Count in leaderboard
func (r *CacheRepository) GetTotalUsers(ctx context.Context) (int64, error) {
	return r.client.ZCard(ctx, LeaderboardKey).Result()
}

// FlushAll - Clear Redis (for testing)
func (r *CacheRepository) FlushAll(ctx context.Context) error {
	return r.client.FlushAll(ctx).Err()
}

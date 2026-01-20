package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/stilln0thing/matiks_leaderboard/internal/models"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}
}

// GetAllUsers - Used for cache warming
func (r *UserRepository) GetAllUsers(ctx context.Context) ([]models.User, error) {
	var users []models.User
	err := r.db.SelectContext(ctx, &users,
		"SELECT id, username, rating, version FROM users ORDER BY id")
	return users, err
}

// GetUserByID - Get single user
func (r *UserRepository) GetUserByID(ctx context.Context, id int64) (*models.User, error) {
	var user models.User
	err := r.db.GetContext(ctx, &user,
		"SELECT id, username, rating, version FROM users WHERE id = $1", id)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// SearchUsers - Fuzzy search by username
func (r *UserRepository) SearchUsers(ctx context.Context, query string) ([]models.User, error) {
	var users []models.User
	searchQuery := "%" + query + "%"
	err := r.db.SelectContext(ctx, &users,
		"SELECT id, username, rating, version FROM users WHERE LOWER(username) LIKE LOWER($1) ORDER BY rating DESC LIMIT 100",
		searchQuery)
	return users, err
}

// BatchUpdateRatings - Efficient batch update with version check
func (r *UserRepository) BatchUpdateRatings(ctx context.Context, updates []models.RatingUpdate) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	stmt, err := tx.PrepareContext(ctx,
		"UPDATE users SET rating = $1, version = $2, updated_at = NOW() WHERE id = $3 AND version < $2")
	if err != nil {
		return err
	}
	defer stmt.Close()
	for _, u := range updates {
		_, err = stmt.ExecContext(ctx, u.Rating, u.Version, u.UserID)
		if err != nil {
			return err
		}
	}
	return tx.Commit()
}

// CreateUser - Create new user
func (r *UserRepository) CreateUser(ctx context.Context, username string, rating int) (*models.User, error) {
	var user models.User
	err := r.db.QueryRowContext(ctx,
		"INSERT INTO users (username, rating, version) VALUES ($1, $2, 0) RETURNING id, username, rating, version",
		username, rating).Scan(&user.ID, &user.Username, &user.Rating, &user.Version)
	return &user, err
}

// GetRandomUserIDs - For simulator
func (r *UserRepository) GetRandomUserIDs(ctx context.Context, limit int) ([]int64, error) {
	var ids []int64
	err := r.db.SelectContext(ctx, &ids,
		"SELECT id FROM users ORDER BY RANDOM() LIMIT $1", limit)
	return ids, err
}

// GetUserCount - Total users
func (r *UserRepository) GetUserCount(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.GetContext(ctx, &count, "SELECT COUNT(*) FROM users")
	return count, err
}

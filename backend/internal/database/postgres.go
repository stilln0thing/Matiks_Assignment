package database

import (
    "time"
    "github.com/jmoiron/sqlx"
    _ "github.com/lib/pq"
)

func NewPostgres(databaseURL string) (*sqlx.DB, error) {
    db, err := sqlx.Connect("postgres", databaseURL)
    if err != nil {
        return nil, err
    }

    db.SetMaxOpenConns(25)
    db.SetMaxIdleConns(5)
    db.SetConnMaxLifetime(5 * time.Minute)
  
    schema := `
    CREATE TABLE IF NOT EXISTS users (
        id SERIAL PRIMARY KEY,
        username VARCHAR(255) UNIQUE NOT NULL,
        rating INTEGER NOT NULL DEFAULT 1000,
        version BIGINT NOT NULL DEFAULT 0,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );
    
    CREATE INDEX IF NOT EXISTS idx_users_rating ON users(rating DESC);
    CREATE INDEX IF NOT EXISTS idx_users_username ON users(username);
    `
    
    _, err = db.Exec(schema)
    return db, err
}
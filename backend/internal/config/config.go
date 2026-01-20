package config

import "os"

type Config struct {
    Port        string
    DatabaseURL string
    RedisAddr   string
}

func Load() *Config {
    return &Config{
        Port:        getEnv("PORT", "8080"),
        DatabaseURL: getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/leaderboard?sslmode=disable"),
        RedisAddr:   getEnv("REDIS_ADDR", "localhost:6379"),
    }
}

func getEnv(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}
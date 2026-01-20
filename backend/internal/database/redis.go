package database

import (
    "context"
    "github.com/redis/go-redis/v9"
)

func NewRedis(addr string) (*redis.Client, error) {
    client := redis.NewClient(&redis.Options{
        Addr:     addr,
        Password: "",
        DB:       0,
        PoolSize: 100,  // Handle concurrent requests
    })
    // Verify connection
    if err := client.Ping(context.Background()).Err(); err != nil {
        return nil, err
    }
    return client, nil
}
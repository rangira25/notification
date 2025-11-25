package storage

import (
	"context"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

func NewRedis(addr, password string) *redis.Client {
	opt := &redis.Options{
		Addr:     addr,
		Password: password,
	}
	client := redis.NewClient(opt)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := client.Ping(ctx).Err(); err != nil {
		log.Fatalf("failed to connect redis: %v", err)
	}
	// Optional: set connection pool options
	client.Options().PoolSize = 10
	client.Options().MinIdleConns = 2
	client.Options().DialTimeout = 5 * time.Second
	return client
}

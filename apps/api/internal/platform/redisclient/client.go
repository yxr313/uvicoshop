package redisclient

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	"uvico-api/internal/platform/config"
)

func New(ctx context.Context, cfg config.RedisConfig) (*redis.Client, error) {
	if cfg.URL == "" {
		return nil, fmt.Errorf("redis: REDIS_URL vacío")
	}

	opts, err := redis.ParseURL(cfg.URL)
	if err != nil {
		return nil, fmt.Errorf("redis: parse url: %w", err)
	}
	if cfg.Password != "" {
		opts.Password = cfg.Password
	}
	opts.DB = cfg.DB
	opts.PoolSize = cfg.PoolSize
	opts.MinIdleConns = cfg.MinIdleConns
	opts.DialTimeout = cfg.DialTimeout
	opts.ReadTimeout = cfg.ReadTimeout
	opts.WriteTimeout = cfg.WriteTimeout

	client := redis.NewClient(opts)

	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := client.Ping(pingCtx).Err(); err != nil {
		_ = client.Close()
		return nil, fmt.Errorf("redis: ping: %w", err)
	}
	return client, nil
}
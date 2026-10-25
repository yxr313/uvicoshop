package cron

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/redis/go-redis/v9"
)

var ErrLockNotHeld = errors.New("cron: lock not held")

type Lock interface {
	Acquire(ctx context.Context, key string, ttl time.Duration) (string, bool, error)
	Release(ctx context.Context, key, token string) error
}

var releaseScript = redis.NewScript(`
	if redis.call("GET", KEYS[1]) == ARGV[1] then
		return redis.call("DEL", KEYS[1])
	else
		return 0
	end
`)

type redisLock struct {
	client *redis.Client
	logger *slog.Logger
}

func NewRedisLock(client *redis.Client, logger *slog.Logger) Lock {
	return &redisLock{client: client, logger: logger}
}

func (l *redisLock) Acquire(ctx context.Context, key string, ttl time.Duration) (string, bool, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", false, err
	}
	token := hex.EncodeToString(b)
	ok, err := l.client.SetNX(ctx, key, token, ttl).Result()
	if err != nil {
		return "", false, fmt.Errorf("redis setnx: %w", err)
	}
	return token, ok, nil
}

func (l *redisLock) Release(ctx context.Context, key, token string) error {
	_, err := releaseScript.Run(ctx, l.client, []string{key}, token).Result()
	return err
}
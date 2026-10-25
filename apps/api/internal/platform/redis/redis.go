package redis

import (
	"context"
	"github.com/redis/go-redis/v9"
	"log/slog"
)

func NewRedisClient(ctx context.Context, redisURL string) (*redis.Client, error) {
	opts, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, err
	}
	client := redis.NewClient(opts)

	if err := client.Ping(ctx).Err(); err != nil {
		slog.Warn("⚠️ No se pudo conectar a Redis (¿Está corriendo?)", "error", err)
	} else {
		slog.Info("✅ Conexión exitosa a Redis")
	}
	return client, nil
}

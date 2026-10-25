package redis

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	goredis "github.com/redis/go-redis/v9"

	storefrontdomain "uvico-api/internal/domain/storefront"
)

const (
	bcvKeyRate      = "bcv:usd:micros"
	bcvKeyUpdatedAt = "bcv:usd:updated_at"
)

type BCVRepo struct {
	client *goredis.Client
}

func NewBCVRepo(client *goredis.Client) *BCVRepo {
	return &BCVRepo{client: client}
}

// GetRate lee la tasa que escribe el cron scraper. Si no existe
// (arranque inicial o scraper caído), retorna error para que el handler
// devuelva 503 — el frontend muestra solo precios en USD.
func (r *BCVRepo) GetRate(ctx context.Context) (*storefrontdomain.BCVRate, error) {
	microsStr, err := r.client.Get(ctx, bcvKeyRate).Result()
	if errors.Is(err, goredis.Nil) {
		return nil, errors.New("bcv: rate no disponible")
	}
	if err != nil {
		return nil, fmt.Errorf("bcv: get rate: %w", err)
	}
	micros, err := strconv.ParseInt(microsStr, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("bcv: parse rate: %w", err)
	}

	updatedAtStr, _ := r.client.Get(ctx, bcvKeyUpdatedAt).Result()
	updatedAt, _ := time.Parse(time.RFC3339, updatedAtStr)

	return &storefrontdomain.BCVRate{
		RateMicros: micros,
		Rate:       formatRate(micros),
		UpdatedAt:  updatedAt,
	}, nil
}

// formatRate convierte micros a string decimal con 4 decimales.
// Ej: 36543200 → "36.5432"
func formatRate(micros int64) string {
	entero := micros / 1_000_000
	decimal := micros % 1_000_000
	return fmt.Sprintf("%d.%04d", entero, decimal/100)
}
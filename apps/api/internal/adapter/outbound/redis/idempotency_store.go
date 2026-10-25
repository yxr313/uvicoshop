package redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	goredis "github.com/redis/go-redis/v9"

	"uvico-api/internal/domain/order"
)

const idempotencyTTL = 24 * time.Hour

type IdempotencyStore struct {
	client *goredis.Client
}

func NewIdempotencyStore(client *goredis.Client) *IdempotencyStore {
	return &IdempotencyStore{client: client}
}

func (s *IdempotencyStore) key(tenantID, idemKey string) string {
	return "checkout:idem:" + tenantID + ":" + idemKey
}

func (s *IdempotencyStore) Get(ctx context.Context, tenantID, key string) (*order.IdempotencyRecord, error) {
	raw, err := s.client.Get(ctx, s.key(tenantID, key)).Bytes()
	if errors.Is(err, goredis.Nil) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("idempotency: get: %w", err)
	}

	var rec order.IdempotencyRecord
	if err := json.Unmarshal(raw, &rec); err != nil {
		return nil, fmt.Errorf("idempotency: unmarshal: %w", err)
	}
	return &rec, nil
}

func (s *IdempotencyStore) Put(ctx context.Context, rec order.IdempotencyRecord) error {
	payload, err := json.Marshal(rec)
	if err != nil {
		return fmt.Errorf("idempotency: marshal: %w", err)
	}

	// SET NX para no sobreescribir un registro existente.
	ok, err := s.client.SetNX(ctx, s.key(rec.TenantID, rec.Key), payload, idempotencyTTL).Result()
	if err != nil {
		return fmt.Errorf("idempotency: setnx: %w", err)
	}
	if !ok {
		// Race: otro proceso ya lo guardó. Verificar que el hash coincida.
		existing, err := s.Get(ctx, rec.TenantID, rec.Key)
		if err != nil {
			return err
		}
		if existing != nil && existing.RequestHash != rec.RequestHash {
			return order.ErrIdempotencyConflict
		}
	}
	return nil
}
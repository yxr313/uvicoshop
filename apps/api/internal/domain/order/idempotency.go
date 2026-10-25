package order

import (
	"context"
	"errors"
)

var (
	ErrIdempotencyConflict = errors.New("order: idempotency key con body distinto")
)

type IdempotencyRecord struct {
	Key         string
	TenantID    string
	OrderID     string
	RequestHash string
}

type IdempotencyStore interface {
	Get(ctx context.Context, tenantID, key string) (*IdempotencyRecord, error)
	Put(ctx context.Context, rec IdempotencyRecord) error
}
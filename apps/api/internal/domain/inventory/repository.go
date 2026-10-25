package inventory

import (
	"context"
	"github.com/jackc/pgx/v5"
	"uvico-api/internal/domain/shared"
)

type Repository interface {
	Reserve(ctx context.Context, tenantID shared.TenantID, productID shared.ProductID, qty int64, referenceID string) error
	Commit(ctx context.Context, tenantID shared.TenantID, productID shared.ProductID, qty int64, referenceID string) error
	Release(ctx context.Context, tenantID shared.TenantID, productID shared.ProductID, qty int64, referenceID string) error
	Adjust(ctx context.Context, tenantID shared.TenantID, productID shared.ProductID, delta int64) error
	GetByProduct(ctx context.Context, tenantID shared.TenantID, productID shared.ProductID) (*Item, error)
}

type TxRepository interface {
	ReleaseTx(ctx context.Context, tx pgx.Tx, tenantID shared.TenantID, productID shared.ProductID, qty int64, referenceID string) error
}
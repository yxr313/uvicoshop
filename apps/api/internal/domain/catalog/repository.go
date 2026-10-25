package catalog

import (
	"context"
	"time"

	"uvico-api/internal/domain/shared"
)

type ListFilter struct {
	TenantID        shared.TenantID
	Status          *Status
	CursorCreatedAt *time.Time
	CursorID        *shared.ProductID
	Limit           int32
}

type Repository interface {
	Create(ctx context.Context, p *Product, initialStock int64) error
	Update(ctx context.Context, p *Product) error
	GetByID(ctx context.Context, tenantID shared.TenantID, id shared.ProductID) (*Product, error)
	GetBySlug(ctx context.Context, tenantID shared.TenantID, slug string) (*Product, error)
	List(ctx context.Context, filter ListFilter) ([]*Product, error)
	SlugExists(ctx context.Context, tenantID shared.TenantID, slug string, excludeID *shared.ProductID) (bool, error)
	SKUExists(ctx context.Context, tenantID shared.TenantID, sku string, excludeID *shared.ProductID) (bool, error)
}
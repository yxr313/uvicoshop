package tenant

import (
	"context"

	"uvico-api/internal/domain/shared"
)

type Repository interface {
	Create(ctx context.Context, t *Tenant) error
	Update(ctx context.Context, t *Tenant) error
	GetByID(ctx context.Context, id shared.TenantID) (*Tenant, error)
	SlugExists(ctx context.Context, slug string) (bool, error)
}
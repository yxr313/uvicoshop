package paymentmethod

import (
	"context"

	"uvico-api/internal/domain/shared"
)

type Repository interface {
	// Replace atomically replaces the tenant's method list. Methods that
	// are not in the input list are set to enabled=false (soft delete).
	Replace(ctx context.Context, tenantID shared.TenantID, methods []*PaymentMethod) error

	// ListEnabled returns only enabled methods (buyer-facing).
	ListEnabled(ctx context.Context, tenantID shared.TenantID) ([]*PaymentMethod, error)

	// ListAll returns all methods including disabled (merchant-facing).
	ListAll(ctx context.Context, tenantID shared.TenantID) ([]*PaymentMethod, error)

	// GetByMethod returns one method or ErrNotFound.
	GetByMethod(ctx context.Context, tenantID shared.TenantID, method Method) (*PaymentMethod, error)
}
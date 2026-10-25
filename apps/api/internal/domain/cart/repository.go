package cart

import (
	"context"

	"uvico-api/internal/domain/shared"
)

type Repository interface {
	// Get retorna el cart o un Cart vacío si no existe.
	Get(ctx context.Context, tenantID shared.TenantID, sessionID string) (*Cart, error)

	// AddItem agrega o incrementa la cantidad si ya existe. Usa Lua
	// para atomicidad (HINCRBY + HSET en un solo paso).
	AddItem(ctx context.Context, tenantID shared.TenantID, sessionID string, item Item) (*Cart, error)

	// SetQuantity reemplaza la cantidad de un item existente.
	SetQuantity(ctx context.Context, tenantID shared.TenantID, sessionID string, productID shared.ProductID, qty int32) (*Cart, error)

	// RemoveItem elimina un item.
	RemoveItem(ctx context.Context, tenantID shared.TenantID, sessionID string, productID shared.ProductID) (*Cart, error)

	// Clear vacía el cart (no lo borra; queda vacío con TTL).
	Clear(ctx context.Context, tenantID shared.TenantID, sessionID string) error

	// Delete borra el cart completo (usado tras checkout exitoso).
	Delete(ctx context.Context, tenantID shared.TenantID, sessionID string) error
}
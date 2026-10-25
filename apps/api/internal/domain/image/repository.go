package image

import (
	"context"

	"uvico-api/internal/domain/shared"
)

type Repository interface {
	// Insert crea la fila en product_images.
	Insert(ctx context.Context, img *Image) error

	// ListByProduct devuelve las imágenes ordenadas por position.
	ListByProduct(ctx context.Context, tenantID shared.TenantID, productID shared.ProductID) ([]*Image, error)

	// GetByID retorna una imagen tenant-scoped.
	GetByID(ctx context.Context, tenantID shared.TenantID, productID shared.ProductID, imageID string) (*Image, error)

	// Delete elimina la fila. Retorna ErrNotFound si no existe.
	Delete(ctx context.Context, tenantID shared.TenantID, productID shared.ProductID, imageID string) error

	// Reorder aplica una nueva permutación de posiciones en una sola tx.
	// `order` es el listado de imageIDs en el nuevo orden; su índice es la
	// nueva position.
	Reorder(ctx context.Context, tenantID shared.TenantID, productID shared.ProductID, order []string) error

	// CountByProduct cuenta las imágenes actuales del producto.
	CountByProduct(ctx context.Context, tenantID shared.TenantID, productID shared.ProductID) (int, error)
}
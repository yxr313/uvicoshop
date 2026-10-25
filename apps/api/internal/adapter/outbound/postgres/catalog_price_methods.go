package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"

	"uvico-api/internal/domain/cart"
	"uvico-api/internal/domain/shared"
)

// GetAnyStatusPrice retorna precio/sku/name/currency/status sin filtrar por
// status='published'. Usado por image.Service para permitir adjuntar
// imágenes a productos en draft antes de publicar.
func (r *CatalogRepo) GetAnyStatusPrice(
	ctx context.Context,
	tenantID shared.TenantID,
	productID shared.ProductID,
) (string, string, int64, string, string, error) {
	var sku, name, currency, status string
	var priceCents int64
	err := r.pool.QueryRow(ctx, `
		SELECT sku, name, price_cents, currency, status
		  FROM products
		 WHERE tenant_id = $1 AND id = $2
	`, tenantID.String(), productID.String()).Scan(&sku, &name, &priceCents, &currency, &status)
	if errors.Is(err, pgx.ErrNoRows) {
		return "", "", 0, "", "", cart.ErrProductNotFound
	}
	if err != nil {
		return "", "", 0, "", "", fmt.Errorf("catalog_repo: get any status price: %w", err)
	}
	return sku, name, priceCents, currency, status, nil
}
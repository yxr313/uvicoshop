package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"uvico-api/internal/domain/catalog"
	"uvico-api/internal/domain/shared"
	"uvico-api/internal/platform/database"
)

type CatalogRepo struct {
	pool *database.Pool
}

func NewCatalogRepo(pool *database.Pool) *CatalogRepo { return &CatalogRepo{pool: pool} }

func (r *CatalogRepo) Create(ctx context.Context, p *catalog.Product, initialStock int64) error {
	return r.pool.WithTenantTx(ctx, p.TenantID.String(), func(ctx context.Context, tx pgx.Tx) error {
		var providerArg any
		if p.ProviderID != nil {
			providerArg = p.ProviderID.String()
		}

		_, err := tx.Exec(ctx, `
			INSERT INTO products (
				id, tenant_id, provider_id, sale_mode, sku, slug, name, description,
				price_cents, currency, status, created_at, updated_at, version
			) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14)
		`, p.ID.String(), p.TenantID.String(), providerArg, string(p.SaleMode),
			p.SKU, p.Slug, p.Name, p.Description,
			p.Price.AmountCents, string(p.Price.Currency), string(p.Status),
			p.CreatedAt, p.UpdatedAt, p.Version)
		if err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && pgErr.Code == "23505" {
				switch pgErr.ConstraintName {
				case "products_tenant_slug_uniq":
					return catalog.ErrSlugTaken
				case "products_tenant_sku_uniq":
					return catalog.ErrSKUTaken
				}
			}
			return fmt.Errorf("catalog_repo: insert: %w", err)
		}

		if p.SaleMode == catalog.SaleOwnInventory {
			if initialStock < 0 {
				return fmt.Errorf("catalog_repo: initialStock no puede ser negativo")
			}
			_, err = tx.Exec(ctx, `
				INSERT INTO inventory_items
					(tenant_id, product_id, on_hand, reserved, version, created_at, updated_at)
				VALUES ($1, $2, $3, 0, 1, NOW(), NOW())
			`, p.TenantID.String(), p.ID.String(), initialStock)
			if err != nil {
				return fmt.Errorf("catalog_repo: init inventory: %w", err)
			}
			if initialStock > 0 {
				_, err = tx.Exec(ctx, `
					INSERT INTO inventory_movements
						(id, tenant_id, product_id, delta, reason, created_at)
					VALUES (gen_random_uuid(), $1, $2, $3, 'adjustment', NOW())
				`, p.TenantID.String(), p.ID.String(), initialStock)
				if err != nil {
					return fmt.Errorf("catalog_repo: init movement: %w", err)
				}
			}
		}

		return insertOutbox(ctx, tx, p.TenantID.String(), "product", p.ID.String(), "product.created", map[string]any{
			"product_id":  p.ID.String(),
			"sale_mode":   string(p.SaleMode),
			"status":      string(p.Status),
			"price_cents": p.Price.AmountCents,
			"currency":    string(p.Price.Currency),
		})
	})
}

func (r *CatalogRepo) Update(ctx context.Context, p *catalog.Product) error {
	return r.pool.WithTenantTx(ctx, p.TenantID.String(), func(ctx context.Context, tx pgx.Tx) error {
		tag, err := tx.Exec(ctx, `
			UPDATE products
			   SET name = $2, description = $3, slug = $4, price_cents = $5,
			       status = $6, published_at = $7, archived_at = $8,
			       updated_at = $9, version = $10
			 WHERE id = $1 AND tenant_id = $11 AND version = $12
		`, p.ID.String(), p.Name, p.Description, p.Slug, p.Price.AmountCents,
			string(p.Status), p.PublishedAt, p.ArchivedAt,
			p.UpdatedAt, p.Version, p.TenantID.String(), p.Version-1)
		if err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && pgErr.Code == "23505" && pgErr.ConstraintName == "products_tenant_slug_uniq" {
				return catalog.ErrSlugTaken
			}
			return fmt.Errorf("catalog_repo: update: %w", err)
		}
		if tag.RowsAffected() == 0 {
			return catalog.ErrVersionConflict
		}
		return insertOutbox(ctx, tx, p.TenantID.String(), "product", p.ID.String(), "product.updated", map[string]any{
			"product_id": p.ID.String(),
			"status":     string(p.Status),
			"version":    p.Version,
		})
	})
}

func (r *CatalogRepo) GetByID(ctx context.Context, tenantID shared.TenantID, id shared.ProductID) (*catalog.Product, error) {
	row := r.pool.QueryRow(ctx, `
		SELECT p.id, p.tenant_id, p.provider_id, p.sale_mode, p.sku, p.slug, p.name, p.description,
		       p.price_cents, p.currency, p.status, p.published_at, p.archived_at,
		       p.created_at, p.updated_at, p.version, COALESCE(i.on_hand, 0) as stock
		  FROM products p
		  LEFT JOIN inventory_items i ON p.id = i.product_id AND p.tenant_id = i.tenant_id
		 WHERE p.id = $1 AND p.tenant_id = $2
	`, id.String(), tenantID.String())
	return scanProduct(row)
}

func (r *CatalogRepo) GetBySlug(ctx context.Context, tenantID shared.TenantID, slug string) (*catalog.Product, error) {
	row := r.pool.QueryRow(ctx, `
		SELECT p.id, p.tenant_id, p.provider_id, p.sale_mode, p.sku, p.slug, p.name, p.description,
		       p.price_cents, p.currency, p.status, p.published_at, p.archived_at,
		       p.created_at, p.updated_at, p.version, COALESCE(i.on_hand, 0) as stock
		  FROM products p
		  LEFT JOIN inventory_items i ON p.id = i.product_id AND p.tenant_id = i.tenant_id
		 WHERE p.tenant_id = $1 AND p.slug = $2
	`, tenantID.String(), slug)
	return scanProduct(row)
}

func (r *CatalogRepo) List(ctx context.Context, f catalog.ListFilter) ([]*catalog.Product, error) {
	limit := f.Limit
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	args := []any{f.TenantID.String(), limit}
	q := `
		SELECT p.id, p.tenant_id, p.provider_id, p.sale_mode, p.sku, p.slug, p.name, p.description,
		       p.price_cents, p.currency, p.status, p.published_at, p.archived_at,
		       p.created_at, p.updated_at, p.version, COALESCE(i.on_hand, 0) as stock
		  FROM products p
		  LEFT JOIN inventory_items i ON p.id = i.product_id AND p.tenant_id = i.tenant_id
		 WHERE p.tenant_id = $1
	`
	if f.Status != nil {
		args = append(args, string(*f.Status))
		q += fmt.Sprintf(" AND p.status = $%d", len(args))
	}
	if f.SaleMode != nil {
		args = append(args, string(*f.SaleMode))
		q += fmt.Sprintf(" AND p.sale_mode = $%d", len(args))
	}
	if f.Search != nil && *f.Search != "" {
		args = append(args, "%"+*f.Search+"%")
		q += fmt.Sprintf(" AND (p.name ILIKE $%d OR p.sku ILIKE $%d)", len(args), len(args))
	}
	if f.CursorCreatedAt != nil && f.CursorID != nil {
		args = append(args, *f.CursorCreatedAt, f.CursorID.String())
		q += fmt.Sprintf(" AND (p.created_at, p.id) < ($%d, $%d)", len(args)-1, len(args))
	}
	q += " ORDER BY p.created_at DESC, p.id DESC LIMIT $2"

	rows, err := r.pool.Query(ctx, q, args...)
	if err != nil {
		return nil, fmt.Errorf("catalog_repo: list: %w", err)
	}
	defer rows.Close()

	var out []*catalog.Product
	for rows.Next() {
		p, err := scanProduct(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, p)
	}
	return out, rows.Err()
}

func (r *CatalogRepo) SlugExists(ctx context.Context, tenantID shared.TenantID, slug string, excludeID *shared.ProductID) (bool, error) {
	var exists bool
	var err error
	if excludeID != nil {
		err = r.pool.QueryRow(ctx, `
			SELECT EXISTS(SELECT 1 FROM products WHERE tenant_id = $1 AND slug = $2 AND id <> $3)
		`, tenantID.String(), slug, excludeID.String()).Scan(&exists)
	} else {
		err = r.pool.QueryRow(ctx, `
			SELECT EXISTS(SELECT 1 FROM products WHERE tenant_id = $1 AND slug = $2)
		`, tenantID.String(), slug).Scan(&exists)
	}
	return exists, err
}

func (r *CatalogRepo) SKUExists(ctx context.Context, tenantID shared.TenantID, sku string, excludeID *shared.ProductID) (bool, error) {
	var exists bool
	var err error
	if excludeID != nil {
		err = r.pool.QueryRow(ctx, `
			SELECT EXISTS(SELECT 1 FROM products WHERE tenant_id = $1 AND sku = $2 AND id <> $3)
		`, tenantID.String(), sku, excludeID.String()).Scan(&exists)
	} else {
		err = r.pool.QueryRow(ctx, `
			SELECT EXISTS(SELECT 1 FROM products WHERE tenant_id = $1 AND sku = $2)
		`, tenantID.String(), sku).Scan(&exists)
	}
	return exists, err
}

func scanProduct(row rowScanner) (*catalog.Product, error) {
	var p catalog.Product
	var id, tenantID, saleMode, currency, status string
	var providerID *string
	var priceCents int64
	var stock int64

	err := row.Scan(
		&id, &tenantID, &providerID, &saleMode, &p.SKU, &p.Slug, &p.Name, &p.Description,
		&priceCents, &currency, &status, &p.PublishedAt, &p.ArchivedAt,
		&p.CreatedAt, &p.UpdatedAt, &p.Version, &stock,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, catalog.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("catalog_repo: scan: %w", err)
	}

	p.ID = shared.ProductID(id)
	p.TenantID = shared.TenantID(tenantID)
	if providerID != nil {
		pid := shared.ProviderID(*providerID)
		p.ProviderID = &pid
	}
	p.SaleMode = catalog.SaleMode(saleMode)
	p.Status = catalog.Status(status)
	p.Price = shared.Money{AmountCents: priceCents, Currency: shared.Currency(currency)}
	p.Stock = stock

	if err := p.Price.Validate(); err != nil {
		return nil, fmt.Errorf("catalog_repo: price corrupto: %w", err)
	}
	return &p, nil
}

func (r *CatalogRepo) GetPublishedPrice(
	ctx context.Context,
	tenantID shared.TenantID,
	productID shared.ProductID,
) (string, string, int64, string, error) {
	var sku, name, currency string
	var priceCents int64
	err := r.pool.QueryRow(ctx, `
		SELECT sku, name, price_cents, currency
		  FROM products
		 WHERE tenant_id = $1 AND id = $2 AND status = 'published'
	`, tenantID.String(), productID.String()).Scan(&sku, &name, &priceCents, &currency)
	if errors.Is(err, pgx.ErrNoRows) {
		return "", "", 0, "", catalog.ErrNotFound // <-- ARREGLADO AQUÍ
	}
	if err != nil {
		return "", "", 0, "", fmt.Errorf("catalog_repo: get price: %w", err)
	}
	return sku, name, priceCents, currency, nil
}

var _ = time.Now
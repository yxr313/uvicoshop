package postgres

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"

	"uvico-api/internal/application/storefront"
	"uvico-api/internal/domain/shared"
	storefrontdomain "uvico-api/internal/domain/storefront"
	"uvico-api/internal/platform/database"
)

type StorefrontRepo struct {
	pool *database.Pool
}

func NewStorefrontRepo(pool *database.Pool) *StorefrontRepo {
	return &StorefrontRepo{pool: pool}
}

func (r *StorefrontRepo) GetTenantBySlug(ctx context.Context, slug string) (*storefront.TenantInfo, error) {
	var id, status, s string
	err := r.pool.QueryRow(ctx, `
		SELECT id, slug, status FROM tenants WHERE slug = $1
	`, slug).Scan(&id, &s, &status)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, storefront.ErrTenantNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("storefront_repo: get tenant: %w", err)
	}
	return &storefront.TenantInfo{
		ID:     shared.TenantID(id),
		Slug:   s,
		Status: status,
	}, nil
}

func (r *StorefrontRepo) GetStoreConfig(ctx context.Context, tenantID shared.TenantID) (*storefrontdomain.StoreConfig, error) {
	var cfg storefrontdomain.StoreConfig
	var tenantSlug, tenantName string
	var puckRaw []byte
	var logo, color, email, phone, wa, insta, tiktok, metaTitle, metaDesc *string

	err := r.pool.QueryRow(ctx, `
		SELECT
			t.id, t.slug, t.name,
			s.logo_url, s.primary_color, s.contact_email, s.contact_phone,
			s.whatsapp, s.currency_display, s.social_instagram, s.social_tiktok,
			s.meta_title, s.meta_description,
			COALESCE(
				(SELECT puck_data FROM store_pages
				  WHERE tenant_id = t.id AND published = true
				  ORDER BY page_slug LIMIT 1),
				'{}'::jsonb
			) AS puck_data
		  FROM tenants t
		  LEFT JOIN store_settings s ON s.tenant_id = t.id
		 WHERE t.id = $1
	`, tenantID.String()).Scan(
		&cfg.TenantID, &tenantSlug, &tenantName,
		&logo, &color, &email, &phone,
		&wa, &cfg.CurrencyDisplay, &insta, &tiktok,
		&metaTitle, &metaDesc, &puckRaw,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, storefront.ErrTenantNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("storefront_repo: get config: %w", err)
	}

	cfg.Slug = tenantSlug
	cfg.Name = tenantName
	cfg.LogoURL = logo
	cfg.PrimaryColor = color
	cfg.ContactEmail = email
	cfg.ContactPhone = phone
	cfg.WhatsApp = wa
	cfg.SocialInstagram = insta
	cfg.SocialTikTok = tiktok
	cfg.MetaTitle = metaTitle
	cfg.MetaDescription = metaDesc
	cfg.PuckData = map[string]any{}
	if len(puckRaw) > 0 {
		_ = json.Unmarshal(puckRaw, &cfg.PuckData)
	}
	return &cfg, nil
}

func (r *StorefrontRepo) ListPublishedProducts(ctx context.Context, tenantID shared.TenantID, limit int) ([]*storefrontdomain.ProductSummary, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT
			p.id, p.sku, p.slug, p.name, p.description, p.price_cents, p.currency, p.sale_mode,
			CASE
				WHEN p.sale_mode = 'dropshipping_verified' THEN NULL
				ELSE COALESCE(i.on_hand - i.reserved, 0)
			END AS available,
			COALESCE(
				(SELECT json_agg(json_build_object('path', pi.storage_path, 'alt', pi.alt_text)
				                 ORDER BY pi.position)
				   FROM product_images pi WHERE pi.product_id = p.id),
				'[]'::json
			) AS images
		  FROM products p
		  LEFT JOIN inventory_items i
		    ON i.tenant_id = p.tenant_id AND i.product_id = p.id
		 WHERE p.tenant_id = $1 AND p.status = 'published'
		 ORDER BY p.created_at DESC, p.id DESC
		 LIMIT $2
	`, tenantID.String(), limit)
	if err != nil {
		return nil, fmt.Errorf("storefront_repo: list products: %w", err)
	}
	defer rows.Close()

	out := make([]*storefrontdomain.ProductSummary, 0)
	for rows.Next() {
		var p storefrontdomain.ProductSummary
		var imagesRaw []byte
		if err := rows.Scan(
			&p.ID, &p.SKU, &p.Slug, &p.Name, &p.Description,
			&p.PriceCents, &p.Currency, &p.SaleMode, &p.Available, &imagesRaw,
		); err != nil {
			return nil, fmt.Errorf("storefront_repo: scan: %w", err)
		}
		p.Images = []storefrontdomain.ProductImageRef{}
		if len(imagesRaw) > 0 {
			_ = json.Unmarshal(imagesRaw, &p.Images)
		}
		out = append(out, &p)
	}
	return out, rows.Err()
}

func (r *StorefrontRepo) GetPublishedProduct(ctx context.Context, tenantID shared.TenantID, productSlug string) (*storefrontdomain.ProductDetail, error) {
	var p storefrontdomain.ProductDetail
	var imagesRaw []byte
	err := r.pool.QueryRow(ctx, `
		SELECT
			p.id, p.sku, p.slug, p.name, p.description, p.price_cents, p.currency, p.sale_mode,
			CASE
				WHEN p.sale_mode = 'dropshipping_verified' THEN NULL
				ELSE COALESCE(i.on_hand - i.reserved, 0)
			END AS available,
			COALESCE(
				(SELECT json_agg(json_build_object('path', pi.storage_path, 'alt', pi.alt_text)
				                 ORDER BY pi.position)
				   FROM product_images pi WHERE pi.product_id = p.id),
				'[]'::json
			) AS images,
			p.updated_at
		  FROM products p
		  LEFT JOIN inventory_items i
		    ON i.tenant_id = p.tenant_id AND i.product_id = p.id
		 WHERE p.tenant_id = $1 AND p.slug = $2 AND p.status = 'published'
	`, tenantID.String(), productSlug).Scan(
		&p.ID, &p.SKU, &p.Slug, &p.Name, &p.Description,
		&p.PriceCents, &p.Currency, &p.SaleMode, &p.Available, &imagesRaw, &p.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, storefront.ErrProductNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("storefront_repo: get product: %w", err)
	}
	p.Images = []storefrontdomain.ProductImageRef{}
	if len(imagesRaw) > 0 {
		_ = json.Unmarshal(imagesRaw, &p.Images)
	}
	return &p, nil
}
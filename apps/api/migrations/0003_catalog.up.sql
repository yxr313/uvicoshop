CREATE TYPE sale_mode AS ENUM ('own_inventory', 'dropshipping_verified');
CREATE TYPE product_status AS ENUM ('draft', 'published', 'archived');

CREATE TABLE products (
    id UUID PRIMARY KEY,
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    provider_id UUID,
    sale_mode sale_mode NOT NULL,
    sku TEXT NOT NULL,
    slug TEXT NOT NULL,
    name TEXT NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    price_cents BIGINT NOT NULL,
    currency VARCHAR(3) NOT NULL DEFAULT 'USD',
    status product_status NOT NULL DEFAULT 'draft',
    published_at TIMESTAMPTZ,
    archived_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    version BIGINT NOT NULL DEFAULT 1,
    CONSTRAINT products_name_len CHECK (char_length(name) BETWEEN 1 AND 120),
    CONSTRAINT products_slug_len CHECK (char_length(slug) BETWEEN 1 AND 80),
    CONSTRAINT products_sku_len CHECK (char_length(sku) BETWEEN 1 AND 64),
    CONSTRAINT products_price_positive CHECK (price_cents > 0),
    CONSTRAINT products_currency_valid CHECK (currency IN ('USD', 'VES')),
    CONSTRAINT products_own_inv_no_provider CHECK (
        (sale_mode = 'own_inventory' AND provider_id IS NULL) OR
        (sale_mode = 'dropshipping_verified' AND provider_id IS NOT NULL)
    )
);

CREATE UNIQUE INDEX products_tenant_slug_uniq ON products (tenant_id, slug);
CREATE UNIQUE INDEX products_tenant_sku_uniq ON products (tenant_id, sku);
CREATE INDEX products_tenant_status_created_idx
    ON products (tenant_id, status, created_at DESC, id DESC);
CREATE INDEX products_tenant_created_idx
    ON products (tenant_id, created_at DESC, id DESC);

ALTER TABLE products ENABLE ROW LEVEL SECURITY;
CREATE POLICY products_tenant_isolation ON products
    USING (tenant_id::text = current_setting('app.tenant_id', true))
    WITH CHECK (tenant_id::text = current_setting('app.tenant_id', true));
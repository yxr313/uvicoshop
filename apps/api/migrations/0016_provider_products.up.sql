CREATE TABLE provider_products (
    provider_id UUID NOT NULL REFERENCES providers(id) ON DELETE CASCADE,
    product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    provider_sku TEXT NOT NULL,
    provider_cost_bps INT NOT NULL DEFAULT 5000,
    lead_time_days INT NOT NULL DEFAULT 3,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (provider_id, product_id),
    CONSTRAINT provider_products_cost_bps CHECK (provider_cost_bps BETWEEN 0 AND 10000),
    CONSTRAINT provider_products_lead_time CHECK (lead_time_days BETWEEN 0 AND 90)
);

CREATE INDEX provider_products_tenant_idx ON provider_products (tenant_id);
CREATE INDEX provider_products_sku_idx ON provider_products (provider_id, provider_sku);
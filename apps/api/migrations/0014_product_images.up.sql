CREATE TABLE product_images (
    id UUID PRIMARY KEY,
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    storage_path TEXT NOT NULL,
    alt_text TEXT NOT NULL DEFAULT '',
    position INT NOT NULL DEFAULT 0,
    width INT,
    height INT,
    size_bytes INT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX product_images_product_idx ON product_images (product_id, position);
CREATE UNIQUE INDEX product_images_path_uniq ON product_images (storage_path);

ALTER TABLE product_images ENABLE ROW LEVEL SECURITY;
CREATE POLICY product_images_tenant_isolation ON product_images
    USING (tenant_id::text = current_setting('app.tenant_id', true))
    WITH CHECK (tenant_id::text = current_setting('app.tenant_id', true));
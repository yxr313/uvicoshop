DO $$
BEGIN
    IF NOT EXISTS (SELECT FROM pg_tables WHERE schemaname = 'public' AND tablename = 'inventory_items') THEN
        CREATE TABLE inventory_items (
            tenant_id  UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
            product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
            on_hand    BIGINT NOT NULL DEFAULT 0,
            reserved   BIGINT NOT NULL DEFAULT 0,
            version    BIGINT NOT NULL DEFAULT 1,
            created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
            updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
            PRIMARY KEY (tenant_id, product_id),
            CONSTRAINT inv_on_hand_nonneg CHECK (on_hand >= 0),
            CONSTRAINT inv_reserved_nonneg CHECK (reserved >= 0),
            CONSTRAINT inv_reserved_le_on_hand CHECK (reserved <= on_hand)
        );
        
        ALTER TABLE inventory_items ENABLE ROW LEVEL SECURITY;
        CREATE POLICY inventory_items_tenant_isolation ON inventory_items
            USING (tenant_id::text = current_setting('app.tenant_id', true))
            WITH CHECK (tenant_id::text = current_setting('app.tenant_id', true));
    END IF;
END
$$;
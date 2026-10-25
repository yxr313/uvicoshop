CREATE TABLE store_settings (
    tenant_id UUID PRIMARY KEY REFERENCES tenants(id) ON DELETE CASCADE,
    logo_url TEXT,
    primary_color VARCHAR(7),
    contact_email TEXT,
    contact_phone TEXT,
    whatsapp TEXT,
    currency_display VARCHAR(3) NOT NULL DEFAULT 'USD',
    social_instagram TEXT,
    social_tiktok TEXT,
    meta_title TEXT,
    meta_description TEXT,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE store_pages (
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    page_slug TEXT NOT NULL DEFAULT '/',
    puck_data JSONB NOT NULL DEFAULT '{}'::jsonb,
    published BOOLEAN NOT NULL DEFAULT false,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (tenant_id, page_slug)
);

CREATE INDEX store_pages_published_idx ON store_pages (tenant_id) WHERE published = true;

ALTER TABLE store_settings ENABLE ROW LEVEL SECURITY;
CREATE POLICY store_settings_tenant_isolation ON store_settings
    USING (tenant_id::text = current_setting('app.tenant_id', true))
    WITH CHECK (tenant_id::text = current_setting('app.tenant_id', true));

ALTER TABLE store_pages ENABLE ROW LEVEL SECURITY;
CREATE POLICY store_pages_tenant_isolation ON store_pages
    USING (tenant_id::text = current_setting('app.tenant_id', true))
    WITH CHECK (tenant_id::text = current_setting('app.tenant_id', true));
ALTER TABLE customers
    ADD COLUMN IF NOT EXISTS marketing_opt_in BOOLEAN NOT NULL DEFAULT false,
    ADD COLUMN IF NOT EXISTS marketing_opt_in_at TIMESTAMPTZ;

CREATE INDEX IF NOT EXISTS customers_marketing_optin_idx
    ON customers (marketing_opt_in_at)
    WHERE marketing_opt_in = true;

CREATE TABLE IF NOT EXISTS order_attribution (
    order_id        UUID PRIMARY KEY REFERENCES orders(id) ON DELETE CASCADE,
    tenant_id       UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    session_id      TEXT,
    utm_source      TEXT,
    utm_medium      TEXT,
    utm_campaign    TEXT,
    utm_term        TEXT,
    utm_content     TEXT,
    fbclid          TEXT,
    gclid           TEXT,
    ttclid          TEXT,
    twclid          TEXT,
    referrer        TEXT,
    landing_page    TEXT,
    user_agent      TEXT,
    accept_language TEXT,
    timezone        TEXT,
    ip_hash         TEXT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS order_attribution_tenant_created_idx ON order_attribution (tenant_id, created_at DESC);
CREATE INDEX IF NOT EXISTS order_attribution_utm_idx ON order_attribution (tenant_id, utm_source, utm_campaign) WHERE utm_source IS NOT NULL;
CREATE INDEX IF NOT EXISTS order_attribution_session_idx ON order_attribution (session_id) WHERE session_id IS NOT NULL;

CREATE TABLE IF NOT EXISTS customer_attribution (
    customer_id     UUID PRIMARY KEY REFERENCES customers(id) ON DELETE CASCADE,
    session_id      TEXT,
    utm_source      TEXT,
    utm_medium      TEXT,
    utm_campaign    TEXT,
    utm_term        TEXT,
    utm_content     TEXT,
    fbclid          TEXT,
    gclid           TEXT,
    ttclid          TEXT,
    twclid          TEXT,
    referrer        TEXT,
    landing_page    TEXT,
    user_agent      TEXT,
    accept_language TEXT,
    timezone        TEXT,
    ip_hash         TEXT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS customer_attribution_utm_idx ON customer_attribution (utm_source, utm_campaign) WHERE utm_source IS NOT NULL;

CREATE TABLE IF NOT EXISTS checkout_idempotency (
    idempotency_key TEXT PRIMARY KEY,
    tenant_id       UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    order_id        UUID NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    request_hash    TEXT NOT NULL,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    expires_at      TIMESTAMPTZ NOT NULL DEFAULT NOW() + INTERVAL '24 hours',
    CONSTRAINT checkout_idempotency_key_len CHECK (char_length(idempotency_key) BETWEEN 8 AND 128)
);

CREATE INDEX IF NOT EXISTS checkout_idempotency_expires_idx ON checkout_idempotency (expires_at);
CREATE INDEX IF NOT EXISTS checkout_idempotency_tenant_idx ON checkout_idempotency (tenant_id, created_at DESC);

ALTER TABLE order_attribution ENABLE ROW LEVEL SECURITY;
ALTER TABLE customer_attribution ENABLE ROW LEVEL SECURITY;
ALTER TABLE checkout_idempotency ENABLE ROW LEVEL SECURITY;

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_policies WHERE tablename = 'order_attribution' AND policyname = 'order_attribution_tenant_isolation') THEN
        CREATE POLICY order_attribution_tenant_isolation ON order_attribution
            USING (tenant_id::text = current_setting('app.tenant_id', true))
            WITH CHECK (tenant_id::text = current_setting('app.tenant_id', true));
    END IF;
    IF NOT EXISTS (SELECT 1 FROM pg_policies WHERE tablename = 'customer_attribution' AND policyname = 'customer_attribution_read') THEN
        CREATE POLICY customer_attribution_read ON customer_attribution FOR SELECT USING (true);
    END IF;
    IF NOT EXISTS (SELECT 1 FROM pg_policies WHERE tablename = 'checkout_idempotency' AND policyname = 'checkout_idempotency_tenant_isolation') THEN
        CREATE POLICY checkout_idempotency_tenant_isolation ON checkout_idempotency
            USING (tenant_id::text = current_setting('app.tenant_id', true))
            WITH CHECK (tenant_id::text = current_setting('app.tenant_id', true));
    END IF;
END $$;
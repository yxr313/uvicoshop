-- 1. Columnas nuevas en orders para tracking del pago manual
ALTER TABLE orders
    ADD COLUMN IF NOT EXISTS payment_method TEXT,
    ADD COLUMN IF NOT EXISTS customer_payment_reference TEXT,
    ADD COLUMN IF NOT EXISTS customer_payment_proof_url TEXT,
    ADD COLUMN IF NOT EXISTS customer_submitted_at TIMESTAMPTZ,
    ADD COLUMN IF NOT EXISTS merchant_confirmed_at TIMESTAMPTZ,
    ADD COLUMN IF NOT EXISTS merchant_confirmed_by TEXT,
    ADD COLUMN IF NOT EXISTS merchant_rejection_reason TEXT,
    ADD COLUMN IF NOT EXISTS access_token UUID;

-- 2. access_token es único (usado por el comprador para subir comprobante)
CREATE UNIQUE INDEX IF NOT EXISTS orders_access_token_uniq
    ON orders (access_token)
    WHERE access_token IS NOT NULL;

-- 3. Dashboard del merchant: órdenes esperando confirmación
CREATE INDEX IF NOT EXISTS orders_awaiting_payment_idx
    ON orders (tenant_id, customer_submitted_at DESC)
    WHERE status = 'awaiting_payment';

-- 4. Métodos de pago por tenant
CREATE TABLE IF NOT EXISTS tenant_payment_methods (
    tenant_id  UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    method     TEXT NOT NULL,
    enabled    BOOLEAN NOT NULL DEFAULT true,
    position   INT NOT NULL DEFAULT 0,
    config     JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (tenant_id, method),
    CONSTRAINT tenant_payment_methods_method_valid CHECK (
        method IN ('pago_movil', 'zelle', 'usdt', 'transfer', 'cash')
    )
);

CREATE INDEX IF NOT EXISTS tenant_payment_methods_tenant_idx
    ON tenant_payment_methods (tenant_id, position)
    WHERE enabled = true;

ALTER TABLE tenant_payment_methods ENABLE ROW LEVEL SECURITY;

DROP POLICY IF EXISTS tenant_payment_methods_tenant_isolation ON tenant_payment_methods;
CREATE POLICY tenant_payment_methods_tenant_isolation ON tenant_payment_methods
    USING (tenant_id::text = current_setting('app.tenant_id', true))
    WITH CHECK (tenant_id::text = current_setting('app.tenant_id', true));
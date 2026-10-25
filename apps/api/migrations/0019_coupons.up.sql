CREATE TYPE coupon_status AS ENUM ('active', 'redeemed', 'expired', 'cancelled');

CREATE TABLE coupons (
    id UUID PRIMARY KEY,
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    code TEXT NOT NULL,
    customer_id UUID NOT NULL REFERENCES customers(id),
    original_order_id UUID NOT NULL REFERENCES orders(id),
    amount_cents BIGINT NOT NULL,
    currency VARCHAR(3) NOT NULL DEFAULT 'USD',
    status coupon_status NOT NULL DEFAULT 'active',
    expires_at TIMESTAMPTZ NOT NULL,
    redeemed_at TIMESTAMPTZ,
    redeemed_order_id UUID REFERENCES orders(id),
    ledger_transaction_id UUID REFERENCES ledger_transactions(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT coupons_amount_positive CHECK (amount_cents > 0)
);

CREATE UNIQUE INDEX coupons_tenant_code_uniq ON coupons (tenant_id, code);
CREATE INDEX coupons_customer_idx ON coupons (customer_id, status);
CREATE INDEX coupons_expiring_idx ON coupons (expires_at) WHERE status = 'active';

ALTER TABLE coupons ENABLE ROW LEVEL SECURITY;
CREATE POLICY coupons_tenant_isolation ON coupons
    USING (tenant_id::text = current_setting('app.tenant_id', true))
    WITH CHECK (tenant_id::text = current_setting('app.tenant_id', true));
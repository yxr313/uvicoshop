CREATE TYPE order_status AS ENUM ('pending_payment', 'paid', 'failed', 'cancelled', 'refunded');
CREATE TYPE order_sale_mode AS ENUM ('own_inventory', 'dropshipping_verified');

CREATE TABLE orders (
    id UUID PRIMARY KEY,
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    customer_id TEXT NOT NULL,
    sale_mode order_sale_mode NOT NULL,
    provider_id UUID,
    status order_status NOT NULL DEFAULT 'pending_payment',
    currency VARCHAR(3) NOT NULL,
    total_cents BIGINT NOT NULL,
    uvico_commission_cents BIGINT NOT NULL,
    provider_payout_cents BIGINT NOT NULL,
    gateway_fee_cents BIGINT NOT NULL,
    merchant_net_cents BIGINT NOT NULL,
    payment_ref TEXT,
    paid_at TIMESTAMPTZ,
    failed_at TIMESTAMPTZ,
    cancelled_at TIMESTAMPTZ,
    refunded_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    version BIGINT NOT NULL DEFAULT 1,
    CONSTRAINT orders_total_positive CHECK (total_cents > 0),
    CONSTRAINT orders_split_invariant CHECK (
        total_cents = uvico_commission_cents + provider_payout_cents + gateway_fee_cents + merchant_net_cents
    )
);

CREATE UNIQUE INDEX orders_payment_ref_uniq ON orders (payment_ref) WHERE payment_ref IS NOT NULL;
CREATE INDEX orders_tenant_created_idx ON orders (tenant_id, created_at DESC);
CREATE INDEX orders_pending_idx ON orders (created_at) WHERE status = 'pending_payment';

CREATE TABLE order_items (
    id UUID PRIMARY KEY,
    order_id UUID NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    product_id UUID NOT NULL,
    sku TEXT NOT NULL,
    name TEXT NOT NULL,
    quantity INT NOT NULL,
    unit_price_cents BIGINT NOT NULL,
    total_price_cents BIGINT NOT NULL,
    CONSTRAINT order_items_qty_positive CHECK (quantity > 0)
);

CREATE INDEX order_items_order_idx ON order_items (order_id);

ALTER TABLE orders ENABLE ROW LEVEL SECURITY;
CREATE POLICY orders_tenant_isolation ON orders
    USING (tenant_id::text = current_setting('app.tenant_id', true))
    WITH CHECK (tenant_id::text = current_setting('app.tenant_id', true));

ALTER TABLE order_items ENABLE ROW LEVEL SECURITY;
CREATE POLICY order_items_tenant_isolation ON order_items
    USING (tenant_id::text = current_setting('app.tenant_id', true));
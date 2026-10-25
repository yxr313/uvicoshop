CREATE TYPE payout_schedule AS ENUM ('weekly', 'fortnightly');
CREATE TYPE payout_cycle_status AS ENUM ('open', 'processing', 'paid', 'failed');
CREATE TYPE payout_status AS ENUM ('pending', 'processing', 'paid', 'failed', 'cancelled');
CREATE TYPE payout_method AS ENUM ('pago_movil', 'zelle', 'usdt', 'transfer');

CREATE TABLE payout_cycles (
    id UUID PRIMARY KEY,
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    schedule payout_schedule NOT NULL,
    period_start TIMESTAMPTZ NOT NULL,
    period_end TIMESTAMPTZ NOT NULL,
    status payout_cycle_status NOT NULL DEFAULT 'open',
    total_cents BIGINT NOT NULL DEFAULT 0,
    currency VARCHAR(3) NOT NULL DEFAULT 'USD',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    processed_at TIMESTAMPTZ,
    CONSTRAINT payout_cycles_total_nonneg CHECK (total_cents >= 0)
);

CREATE INDEX payout_cycles_tenant_created_idx ON payout_cycles (tenant_id, created_at DESC);
CREATE INDEX payout_cycles_status_idx ON payout_cycles (status) WHERE status IN ('open', 'processing');

CREATE TABLE payouts (
    id UUID PRIMARY KEY,
    cycle_id UUID NOT NULL REFERENCES payout_cycles(id) ON DELETE CASCADE,
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    amount_cents BIGINT NOT NULL,
    currency VARCHAR(3) NOT NULL DEFAULT 'USD',
    status payout_status NOT NULL DEFAULT 'pending',
    method payout_method,
    reference TEXT,
    ledger_transaction_id UUID REFERENCES ledger_transactions(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    paid_at TIMESTAMPTZ,
    CONSTRAINT payouts_amount_positive CHECK (amount_cents > 0)
);

CREATE INDEX payouts_tenant_created_idx ON payouts (tenant_id, created_at DESC);
CREATE INDEX payouts_pending_idx ON payouts (created_at) WHERE status = 'pending';

ALTER TABLE payout_cycles ENABLE ROW LEVEL SECURITY;
CREATE POLICY payout_cycles_tenant_isolation ON payout_cycles
    USING (tenant_id::text = current_setting('app.tenant_id', true))
    WITH CHECK (tenant_id::text = current_setting('app.tenant_id', true));

ALTER TABLE payouts ENABLE ROW LEVEL SECURITY;
CREATE POLICY payouts_tenant_isolation ON payouts
    USING (tenant_id::text = current_setting('app.tenant_id', true))
    WITH CHECK (tenant_id::text = current_setting('app.tenant_id', true));
CREATE TYPE ledger_account_kind AS ENUM ('asset', 'liability', 'revenue', 'expense', 'equity');
CREATE TYPE ledger_account_scope AS ENUM ('global', 'tenant');

CREATE TABLE ledger_accounts (
    id UUID PRIMARY KEY,
    tenant_id UUID REFERENCES tenants(id) ON DELETE CASCADE,
    scope ledger_account_scope NOT NULL,
    code TEXT NOT NULL,
    name TEXT NOT NULL,
    kind ledger_account_kind NOT NULL,
    currency VARCHAR(3) NOT NULL DEFAULT 'USD',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT ledger_accounts_scope_consistency CHECK (
        (scope = 'global' AND tenant_id IS NULL) OR
        (scope = 'tenant' AND tenant_id IS NOT NULL)
    )
);

CREATE UNIQUE INDEX ledger_accounts_code_scope_uniq
    ON ledger_accounts (COALESCE(tenant_id, '00000000-0000-0000-0000-000000000000'::uuid), code);
CREATE INDEX ledger_accounts_tenant_idx ON ledger_accounts (tenant_id) WHERE tenant_id IS NOT NULL;
CREATE INDEX ledger_accounts_kind_idx ON ledger_accounts (kind);

CREATE TABLE ledger_transactions (
    id UUID PRIMARY KEY,
    tenant_id UUID REFERENCES tenants(id) ON DELETE CASCADE,
    description TEXT NOT NULL,
    reference_type TEXT NOT NULL,
    reference_id UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT ledger_transactions_reference_type CHECK (
        reference_type IN (
            'order', 'payout', 'refund', 'chargeback', 'subscription',
            'pod_order', 'capital_advance', 'marketplace_fee', 'ads_spend',
            'manual_adjustment', 'gateway_settlement'
        )
    )
);

CREATE INDEX ledger_transactions_tenant_created_idx
    ON ledger_transactions (tenant_id, created_at DESC)
    WHERE tenant_id IS NOT NULL;
CREATE INDEX ledger_transactions_reference_idx
    ON ledger_transactions (reference_type, reference_id)
    WHERE reference_id IS NOT NULL;

CREATE TABLE ledger_entries (
    id UUID NOT NULL,
    transaction_id UUID NOT NULL,
    account_id UUID NOT NULL,
    amount_cents BIGINT NOT NULL,
    currency VARCHAR(3) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (id, created_at)
) PARTITION BY RANGE (created_at);

CREATE INDEX ledger_entries_transaction_idx ON ledger_entries (transaction_id);
CREATE INDEX ledger_entries_account_idx ON ledger_entries (account_id, created_at DESC);

-- Particiones mensuales 2025-10 → 2026-12
CREATE TABLE ledger_entries_2025_10 PARTITION OF ledger_entries FOR VALUES FROM ('2025-10-01') TO ('2025-11-01');
CREATE TABLE ledger_entries_2025_11 PARTITION OF ledger_entries FOR VALUES FROM ('2025-11-01') TO ('2025-12-01');
CREATE TABLE ledger_entries_2025_12 PARTITION OF ledger_entries FOR VALUES FROM ('2025-12-01') TO ('2026-01-01');
CREATE TABLE ledger_entries_2026_01 PARTITION OF ledger_entries FOR VALUES FROM ('2026-01-01') TO ('2026-02-01');
CREATE TABLE ledger_entries_2026_02 PARTITION OF ledger_entries FOR VALUES FROM ('2026-02-01') TO ('2026-03-01');
CREATE TABLE ledger_entries_2026_03 PARTITION OF ledger_entries FOR VALUES FROM ('2026-03-01') TO ('2026-04-01');
CREATE TABLE ledger_entries_2026_04 PARTITION OF ledger_entries FOR VALUES FROM ('2026-04-01') TO ('2026-05-01');
CREATE TABLE ledger_entries_2026_05 PARTITION OF ledger_entries FOR VALUES FROM ('2026-05-01') TO ('2026-06-01');
CREATE TABLE ledger_entries_2026_06 PARTITION OF ledger_entries FOR VALUES FROM ('2026-06-01') TO ('2026-07-01');
CREATE TABLE ledger_entries_2026_07 PARTITION OF ledger_entries FOR VALUES FROM ('2026-07-01') TO ('2026-08-01');
CREATE TABLE ledger_entries_2026_08 PARTITION OF ledger_entries FOR VALUES FROM ('2026-08-01') TO ('2026-09-01');
CREATE TABLE ledger_entries_2026_09 PARTITION OF ledger_entries FOR VALUES FROM ('2026-09-01') TO ('2026-10-01');
CREATE TABLE ledger_entries_2026_10 PARTITION OF ledger_entries FOR VALUES FROM ('2026-10-01') TO ('2026-11-01');
CREATE TABLE ledger_entries_2026_11 PARTITION OF ledger_entries FOR VALUES FROM ('2026-11-01') TO ('2026-12-01');
CREATE TABLE ledger_entries_2026_12 PARTITION OF ledger_entries FOR VALUES FROM ('2026-12-01') TO ('2027-01-01');

-- Invariante: cada ledger_transaction debe balancear (suma de amount_cents = 0).
-- DEFERRABLE INITIALLY DEFERRED permite insertar N entries y validar al commit.
CREATE OR REPLACE FUNCTION verify_ledger_balance()
RETURNS TRIGGER AS $$
DECLARE
    total BIGINT;
BEGIN
    SELECT COALESCE(SUM(amount_cents), 0) INTO total
      FROM ledger_entries
     WHERE transaction_id = NEW.transaction_id;
    IF total != 0 THEN
        RAISE EXCEPTION 'ledger transaction % no balancea: sum=%', NEW.transaction_id, total;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE CONSTRAINT TRIGGER trg_ledger_balance
    AFTER INSERT ON ledger_entries
    DEFERRABLE INITIALLY DEFERRED
    FOR EACH ROW EXECUTE FUNCTION verify_ledger_balance();

-- Cuentas globales (seed). UUIDs fijos para estabilidad entre entornos.
INSERT INTO ledger_accounts (id, tenant_id, scope, code, name, kind, currency) VALUES
    ('00000000-0000-0000-0000-000000001001', NULL, 'global', '1001', 'Uvico Cash',                    'asset',     'USD'),
    ('00000000-0000-0000-0000-000000001002', NULL, 'global', '1002', 'Gateway Settlement Clearing',   'asset',     'USD'),
    ('00000000-0000-0000-0000-000000001003', NULL, 'global', '1003', 'Gateway Fee Receivable',        'asset',     'USD'),
    ('00000000-0000-0000-0000-000000003001', NULL, 'global', '3001', 'Subscription Revenue',          'revenue',   'USD'),
    ('00000000-0000-0000-0000-000000003002', NULL, 'global', '3002', 'Commission Revenue',            'revenue',   'USD'),
    ('00000000-0000-0000-0000-000000003003', NULL, 'global', '3003', 'POD Margin Revenue',            'revenue',   'USD'),
    ('00000000-0000-0000-0000-000000003004', NULL, 'global', '3004', 'Uvico Pay Fee Revenue',         'revenue',   'USD'),
    ('00000000-0000-0000-0000-000000003005', NULL, 'global', '3005', 'Capital Interest Revenue',      'revenue',   'USD'),
    ('00000000-0000-0000-0000-000000003006', NULL, 'global', '3006', 'Marketplace Revenue Share',     'revenue',   'USD'),
    ('00000000-0000-0000-0000-000000003007', NULL, 'global', '3007', 'Provider Ads Revenue',          'revenue',   'USD'),
    ('00000000-0000-0000-0000-000000004001', NULL, 'global', '4001', 'Chargeback Expense',            'expense',   'USD'),
    ('00000000-0000-0000-0000-000000004002', NULL, 'global', '4002', 'Refund Expense',                'expense',   'USD'),
    ('00000000-0000-0000-0000-000000004003', NULL, 'global', '4003', 'Gateway Fees Paid',             'expense',   'USD')
ON CONFLICT DO NOTHING;

-- Trigger: cada tenant nuevo recibe sus cuentas de pasivo automáticamente.
CREATE OR REPLACE FUNCTION create_tenant_ledger_accounts()
RETURNS TRIGGER AS $$
BEGIN
    INSERT INTO ledger_accounts (id, tenant_id, scope, code, name, kind, currency)
    VALUES
        (gen_random_uuid(), NEW.id, 'tenant', '2001', 'Merchant Payable', 'liability', 'USD'),
        (gen_random_uuid(), NEW.id, 'tenant', '2002', 'Provider Payable', 'liability', 'USD');
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_create_tenant_ledger_accounts
    AFTER INSERT ON tenants
    FOR EACH ROW EXECUTE FUNCTION create_tenant_ledger_accounts();

ALTER TABLE ledger_accounts ENABLE ROW LEVEL SECURITY;
CREATE POLICY ledger_accounts_tenant_read ON ledger_accounts FOR SELECT
    USING (
        scope = 'global'
        OR tenant_id::text = current_setting('app.tenant_id', true)
    );

ALTER TABLE ledger_transactions ENABLE ROW LEVEL SECURITY;
CREATE POLICY ledger_transactions_tenant_isolation ON ledger_transactions
    USING (tenant_id IS NULL OR tenant_id::text = current_setting('app.tenant_id', true));

ALTER TABLE ledger_entries ENABLE ROW LEVEL SECURITY;
CREATE POLICY ledger_entries_tenant_isolation ON ledger_entries
    USING (
        EXISTS (
            SELECT 1 FROM ledger_transactions t
             WHERE t.id = ledger_entries.transaction_id
               AND (t.tenant_id IS NULL OR t.tenant_id::text = current_setting('app.tenant_id', true))
        )
    );
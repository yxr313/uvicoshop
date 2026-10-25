CREATE TABLE providers (
    id UUID PRIMARY KEY,
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    cost_bps INT NOT NULL DEFAULT 5000,
    status TEXT NOT NULL DEFAULT 'active',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT providers_name_len CHECK (char_length(name) BETWEEN 1 AND 120),
    CONSTRAINT providers_cost_bps_range CHECK (cost_bps BETWEEN 0 AND 10000),
    CONSTRAINT providers_status_valid CHECK (status IN ('active', 'inactive', 'suspended'))
);

CREATE INDEX providers_tenant_status_idx
    ON providers (tenant_id, status)
    WHERE status = 'active';

ALTER TABLE providers ENABLE ROW LEVEL SECURITY;
CREATE POLICY providers_tenant_isolation ON providers
    USING (tenant_id::text = current_setting('app.tenant_id', true))
    WITH CHECK (tenant_id::text = current_setting('app.tenant_id', true));

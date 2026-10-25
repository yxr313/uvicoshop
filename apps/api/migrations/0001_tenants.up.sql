CREATE TYPE tenant_status AS ENUM (
    'trial', 'active', 'past_due', 'paused', 'pending_deletion', 'deleted'
);

CREATE TABLE tenants (
    id UUID PRIMARY KEY,
    name TEXT NOT NULL,
    slug TEXT NOT NULL,
    owner_user_id UUID NOT NULL,
    status tenant_status NOT NULL DEFAULT 'trial',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    version BIGINT NOT NULL DEFAULT 1,
    CONSTRAINT tenants_name_len CHECK (char_length(name) BETWEEN 1 AND 120),
    CONSTRAINT tenants_slug_len CHECK (char_length(slug) BETWEEN 1 AND 60)
);

CREATE UNIQUE INDEX tenants_slug_unique_idx ON tenants (slug);
CREATE INDEX tenants_owner_idx ON tenants (owner_user_id);
CREATE INDEX tenants_status_idx ON tenants (status);

ALTER TABLE tenants ENABLE ROW LEVEL SECURITY;

CREATE POLICY tenants_self_read ON tenants FOR SELECT
    USING (id::text = current_setting('app.tenant_id', true));

CREATE POLICY tenants_self_update ON tenants FOR UPDATE
    USING (id::text = current_setting('app.tenant_id', true))
    WITH CHECK (id::text = current_setting('app.tenant_id', true));
CREATE TYPE plan_id AS ENUM ('founder', 'essential', 'professional', 'enterprise');
CREATE TYPE billing_interval AS ENUM ('monthly', 'annual');
CREATE TYPE subscription_status AS ENUM (
    'trialing', 'active', 'past_due', 'paused', 'pending_deletion', 'cancelled', 'deleted'
);

CREATE TABLE subscriptions (
    id UUID PRIMARY KEY,
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    plan_id plan_id NOT NULL,
    billing_interval billing_interval NOT NULL,
    status subscription_status NOT NULL,
    trial_ends_at TIMESTAMPTZ,
    current_period_starts_at TIMESTAMPTZ,
    current_period_ends_at TIMESTAMPTZ,
    billing_cycles_completed INT NOT NULL DEFAULT 0,
    past_due_since TIMESTAMPTZ,
    paused_at TIMESTAMPTZ,
    deletion_scheduled_at TIMESTAMPTZ,
    cancelled_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    version BIGINT NOT NULL DEFAULT 1,
    CONSTRAINT subscriptions_tenant_unique UNIQUE (tenant_id),
    CONSTRAINT subscriptions_cycles_nonneg CHECK (billing_cycles_completed >= 0)
);

CREATE INDEX subscriptions_status_idx ON subscriptions (status);
CREATE INDEX subscriptions_period_end_idx
    ON subscriptions (current_period_ends_at)
    WHERE status IN ('active', 'past_due');
CREATE INDEX subscriptions_trial_end_idx
    ON subscriptions (trial_ends_at)
    WHERE status = 'trialing';
CREATE INDEX subscriptions_paused_at_idx
    ON subscriptions (paused_at)
    WHERE status IN ('paused', 'pending_deletion');

ALTER TABLE subscriptions ENABLE ROW LEVEL SECURITY;

CREATE POLICY subscriptions_self_read ON subscriptions FOR SELECT
    USING (tenant_id::text = current_setting('app.tenant_id', true));
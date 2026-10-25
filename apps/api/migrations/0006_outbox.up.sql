CREATE TABLE outbox_events (
    id              UUID PRIMARY KEY,
    tenant_id       UUID NOT NULL,
    aggregate_type  TEXT NOT NULL,
    aggregate_id    UUID NOT NULL,
    event_type      TEXT NOT NULL,
    payload         JSONB NOT NULL,
    headers         JSONB NOT NULL DEFAULT '{}'::jsonb,
    attempts        SMALLINT NOT NULL DEFAULT 0,
    max_attempts    SMALLINT NOT NULL DEFAULT 10,
    next_attempt_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    locked_at       TIMESTAMPTZ,
    locked_by       TEXT,
    last_error      TEXT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    processed_at    TIMESTAMPTZ,
    CONSTRAINT outbox_payload_size CHECK (octet_length(payload::text) <= 262144)
);

CREATE INDEX outbox_pending_idx
    ON outbox_events (next_attempt_at, created_at)
    WHERE processed_at IS NULL AND attempts < max_attempts;
CREATE INDEX outbox_aggregate_idx
    ON outbox_events (tenant_id, aggregate_type, aggregate_id, created_at DESC);
CREATE INDEX outbox_dlq_idx
    ON outbox_events (created_at)
    WHERE processed_at IS NULL AND attempts >= max_attempts;

-- El worker usa service_role. El frontend nunca debe leer esta tabla.
REVOKE ALL ON outbox_events FROM anon, authenticated;
CREATE TABLE processed_webhooks (
    event_id TEXT PRIMARY KEY,
    provider TEXT NOT NULL,
    order_id UUID NOT NULL,
    processed_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX processed_webhooks_order_idx ON processed_webhooks (order_id);
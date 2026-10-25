CREATE TYPE inventory_reason AS ENUM ('reserved', 'committed', 'released', 'adjustment');

CREATE TABLE inventory_movements (
    id           UUID NOT NULL,
    tenant_id    UUID NOT NULL,
    product_id   UUID NOT NULL,
    delta        BIGINT NOT NULL,
    reason       inventory_reason NOT NULL,
    reference_id UUID,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (id, created_at)
) PARTITION BY RANGE (created_at);

CREATE INDEX inv_mov_product_idx
    ON inventory_movements (tenant_id, product_id, created_at DESC);
CREATE INDEX inv_mov_reference_idx
    ON inventory_movements (reference_id)
    WHERE reference_id IS NOT NULL;

-- Particiones de 2025-10 a 2026-12. El job outbox-partitions crea la siguiente cada día 25.
CREATE TABLE inventory_movements_2025_10 PARTITION OF inventory_movements FOR VALUES FROM ('2025-10-01') TO ('2025-11-01');
CREATE TABLE inventory_movements_2025_11 PARTITION OF inventory_movements FOR VALUES FROM ('2025-11-01') TO ('2025-12-01');
CREATE TABLE inventory_movements_2025_12 PARTITION OF inventory_movements FOR VALUES FROM ('2025-12-01') TO ('2026-01-01');
CREATE TABLE inventory_movements_2026_01 PARTITION OF inventory_movements FOR VALUES FROM ('2026-01-01') TO ('2026-02-01');
CREATE TABLE inventory_movements_2026_02 PARTITION OF inventory_movements FOR VALUES FROM ('2026-02-01') TO ('2026-03-01');
CREATE TABLE inventory_movements_2026_03 PARTITION OF inventory_movements FOR VALUES FROM ('2026-03-01') TO ('2026-04-01');
CREATE TABLE inventory_movements_2026_04 PARTITION OF inventory_movements FOR VALUES FROM ('2026-04-01') TO ('2026-05-01');
CREATE TABLE inventory_movements_2026_05 PARTITION OF inventory_movements FOR VALUES FROM ('2026-05-01') TO ('2026-06-01');
CREATE TABLE inventory_movements_2026_06 PARTITION OF inventory_movements FOR VALUES FROM ('2026-06-01') TO ('2026-07-01');
CREATE TABLE inventory_movements_2026_07 PARTITION OF inventory_movements FOR VALUES FROM ('2026-07-01') TO ('2026-08-01');
CREATE TABLE inventory_movements_2026_08 PARTITION OF inventory_movements FOR VALUES FROM ('2026-08-01') TO ('2026-09-01');
CREATE TABLE inventory_movements_2026_09 PARTITION OF inventory_movements FOR VALUES FROM ('2026-09-01') TO ('2026-10-01');
CREATE TABLE inventory_movements_2026_10 PARTITION OF inventory_movements FOR VALUES FROM ('2026-10-01') TO ('2026-11-01');
CREATE TABLE inventory_movements_2026_11 PARTITION OF inventory_movements FOR VALUES FROM ('2026-11-01') TO ('2026-12-01');
CREATE TABLE inventory_movements_2026_12 PARTITION OF inventory_movements FOR VALUES FROM ('2026-12-01') TO ('2027-01-01');

ALTER TABLE inventory_movements ENABLE ROW LEVEL SECURITY;
CREATE POLICY inventory_movements_tenant_isolation ON inventory_movements
    USING (tenant_id::text = current_setting('app.tenant_id', true));
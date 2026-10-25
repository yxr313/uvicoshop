ALTER TABLE orders ADD COLUMN IF NOT EXISTS customer_uuid UUID REFERENCES customers(id) ON DELETE RESTRICT;

INSERT INTO customers (id, email, name, created_at, updated_at)
SELECT DISTINCT gen_random_uuid(), LOWER(o.customer_id || '@legacy.local'),
       'Legacy Customer ' || o.customer_id, NOW(), NOW()
FROM orders o
WHERE o.customer_id IS NOT NULL AND o.customer_id <> '' AND o.customer_uuid IS NULL
ON CONFLICT ((lower(email))) DO NOTHING;

UPDATE orders o SET customer_uuid = c.id
FROM customers c
WHERE LOWER(c.email) = LOWER(o.customer_id || '@legacy.local')
  AND o.customer_id IS NOT NULL AND o.customer_id <> '' AND o.customer_uuid IS NULL;

UPDATE orders o SET customer_uuid = o.customer_id::uuid
WHERE o.customer_uuid IS NULL
  AND o.customer_id ~ '^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$'
  AND EXISTS (SELECT 1 FROM customers c WHERE c.id = o.customer_id::uuid);

INSERT INTO customers (id, email, name, created_at, updated_at)
SELECT DISTINCT gen_random_uuid(), 'orphan-' || o.id || '@legacy.local',
       'Orphan Order ' || o.id, NOW(), NOW()
FROM orders o WHERE o.customer_uuid IS NULL
ON CONFLICT ((lower(email))) DO NOTHING;

UPDATE orders o SET customer_uuid = c.id
FROM customers c
WHERE LOWER(c.email) = LOWER('orphan-' || o.id || '@legacy.local') AND o.customer_uuid IS NULL;

ALTER TABLE orders ALTER COLUMN customer_uuid SET NOT NULL;
ALTER TABLE orders DROP COLUMN IF EXISTS customer_id;
ALTER TABLE orders RENAME COLUMN customer_uuid TO customer_id;

CREATE INDEX IF NOT EXISTS orders_tenant_customer_idx ON orders (tenant_id, customer_id, created_at DESC);

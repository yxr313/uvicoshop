DROP INDEX IF EXISTS orders_tenant_customer_idx;
ALTER TABLE orders ADD COLUMN IF NOT EXISTS customer_id_text TEXT;
UPDATE orders SET customer_id_text = customer_id::text;
ALTER TABLE orders DROP COLUMN IF EXISTS customer_id;
ALTER TABLE orders RENAME COLUMN customer_id_text TO customer_id;

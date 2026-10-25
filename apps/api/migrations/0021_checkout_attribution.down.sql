DROP TABLE IF EXISTS checkout_idempotency;
DROP TABLE IF EXISTS customer_attribution;
DROP TABLE IF EXISTS order_attribution;
DROP INDEX IF EXISTS customers_marketing_optin_idx;
ALTER TABLE customers DROP COLUMN IF EXISTS marketing_opt_in;
ALTER TABLE customers DROP COLUMN IF EXISTS marketing_opt_in_at;
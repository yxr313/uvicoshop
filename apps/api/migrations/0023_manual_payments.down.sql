DROP TABLE IF EXISTS tenant_payment_methods;
DROP INDEX IF EXISTS orders_awaiting_payment_idx;
DROP INDEX IF EXISTS orders_access_token_uniq;
ALTER TABLE orders
    DROP COLUMN IF EXISTS payment_method,
    DROP COLUMN IF EXISTS customer_payment_reference,
    DROP COLUMN IF EXISTS customer_payment_proof_url,
    DROP COLUMN IF EXISTS customer_submitted_at,
    DROP COLUMN IF EXISTS merchant_confirmed_at,
    DROP COLUMN IF EXISTS merchant_confirmed_by,
    DROP COLUMN IF EXISTS merchant_rejection_reason,
    DROP COLUMN IF EXISTS access_token;
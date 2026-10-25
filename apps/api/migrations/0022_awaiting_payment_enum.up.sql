-- ALTER TYPE ... ADD VALUE no puede ejecutarse dentro de una transacción
-- en Postgres < 12. En PG 12+ se permite, pero el nuevo valor no puede usarse
-- en la misma transacción. Por eso separamos este ALTER en su propio archivo.
ALTER TYPE order_status ADD VALUE IF NOT EXISTS 'awaiting_payment' BEFORE 'paid';
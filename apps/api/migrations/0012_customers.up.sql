CREATE TABLE customers (
    id UUID PRIMARY KEY,
    email TEXT NOT NULL,
    phone TEXT,
    name TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX customers_email_uniq ON customers (lower(email));
CREATE INDEX customers_phone_idx ON customers (phone) WHERE phone IS NOT NULL;
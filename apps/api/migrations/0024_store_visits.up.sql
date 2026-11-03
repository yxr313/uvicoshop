CREATE TABLE store_visits (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    session_id TEXT NOT NULL,
    path TEXT NOT NULL,
    device TEXT,
    browser TEXT,
    os TEXT,
    referrer TEXT,
    utm_source TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX store_visits_tenant_created_idx ON store_visits (tenant_id, created_at DESC);

-- RLS para que solo el dueño de la tienda vea sus visitas
ALTER TABLE store_visits ENABLE ROW LEVEL SECURITY;

CREATE POLICY store_visits_tenant_isolation ON store_visits
    USING (tenant_id::text = current_setting('app.tenant_id', true))
    WITH CHECK (tenant_id::text = current_setting('app.tenant_id', true));

-- Permitir al frontend (Realtime y consultas directas si se necesita) leer las visitas de su tienda
CREATE POLICY store_visits_select ON store_visits
    FOR SELECT TO authenticated
    USING (tenant_id IN (SELECT id FROM tenants WHERE owner_user_id = auth.uid()));
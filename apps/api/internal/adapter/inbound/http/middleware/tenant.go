package middleware

import (
	"context"
	"net/http"

	"uvico-api/internal/domain/shared"
)

// WithTenant inyecta el tenant en el contexto. Usada por JWTService.Middleware
// y por TenantResolver. Es la ÚNICA función que escribe ctxKeyTenantID.
func WithTenant(ctx context.Context, tenantID shared.TenantID) context.Context {
	return context.WithValue(ctx, ctxKeyTenantID, tenantID)
}

// TenantFromContext devuelve el tenant del contexto.
func TenantFromContext(ctx context.Context) (shared.TenantID, bool) {
	tid, ok := ctx.Value(ctxKeyTenantID).(shared.TenantID)
	if !ok || tid.IsZero() {
		return "", false
	}
	return tid, true
}

// MustTenant retorna el tenant o error. Es el ÚNICO path recomendado para handlers.
func MustTenant(ctx context.Context) (shared.TenantID, error) {
	tid, ok := TenantFromContext(ctx)
	if !ok {
		return "", errMissingTenant
	}
	return tid, nil
}

// TenantResolver: solo para dev/tests sin JWT. En producción el JWT
// middleware ya inyectó el tenant. No-op si ya está presente.
func TenantResolver(headerName string, isProduction bool) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if _, ok := TenantFromContext(r.Context()); ok {
				next.ServeHTTP(w, r)
				return
			}
			if isProduction {
				writeJSONError(w, http.StatusUnauthorized, "missing_tenant", "tenant no resuelto del token")
				return
			}
			raw := r.Header.Get(headerName)
			if raw == "" {
				writeJSONError(w, http.StatusBadRequest, "missing_tenant_header", headerName+" requerido")
				return
			}
			tid, err := shared.ParseTenantID(raw)
			if err != nil {
				writeJSONError(w, http.StatusBadRequest, "invalid_tenant", "tenant id inválido")
				return
			}
			next.ServeHTTP(w, r.WithContext(WithTenant(r.Context(), tid)))
		})
	}
}

// TenantIDFromContext es un alias de TenantFromContext. Se mantiene por
// compatibilidad con middleware.RateLimit y otros callers que solo necesitan
// leer el tenant. En código nuevo, usa TenantFromContext.
func TenantIDFromContext(ctx context.Context) (shared.TenantID, bool) {
	return TenantFromContext(ctx)
}

// MustTenantID es un alias de MustTenant. Se mantiene por compatibilidad.
// En código nuevo, usa MustTenant.
func MustTenantID(ctx context.Context) (shared.TenantID, error) {
	return MustTenant(ctx)
}
package storefront

import (
	"context"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"

	appstorefront "uvico-api/internal/application/storefront"
)

type ctxKey int

const ctxKeyTenant ctxKey = iota

type SlugResolver struct {
	svc *appstorefront.Service
}

func NewSlugResolver(svc *appstorefront.Service) *SlugResolver {
	return &SlugResolver{svc: svc}
}

func (r *SlugResolver) Middleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			slug := chi.URLParam(req, "slug")
			if slug == "" {
				writeJSONError(w, http.StatusBadRequest, "missing_slug", "slug requerido")
				return
			}

			tenant, err := r.svc.ResolveTenant(req.Context(), slug)
			if err != nil {
				switch {
				case errors.Is(err, appstorefront.ErrTenantNotFound):
					writeJSONError(w, http.StatusNotFound, "store_not_found", "tienda no encontrada")
				case errors.Is(err, appstorefront.ErrStoreUnavailable):
					writeJSONError(w, http.StatusForbidden, "store_unavailable", "tienda no disponible")
				default:
					writeJSONError(w, http.StatusInternalServerError, "internal", "error interno")
				}
				return
			}

			ctx := context.WithValue(req.Context(), ctxKeyTenant, tenant)
			next.ServeHTTP(w, req.WithContext(ctx))
		})
	}
}

func TenantFromContext(ctx context.Context) (*appstorefront.TenantInfo, bool) {
	t, ok := ctx.Value(ctxKeyTenant).(*appstorefront.TenantInfo)
	return t, ok && t != nil
}

func MustTenant(ctx context.Context) (*appstorefront.TenantInfo, error) {
	t, ok := TenantFromContext(ctx)
	if !ok {
		return nil, errors.New("storefront: tenant no en contexto")
	}
	return t, nil
}

func writeJSONError(w http.ResponseWriter, status int, code, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-store")
	w.WriteHeader(status)
	_, _ = w.Write([]byte(`{"error":"` + code + `","message":"` + message + `"}`))
}
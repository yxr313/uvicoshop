package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"

	"uvico-api/internal/domain/shared"
)

// RateLimit implementa ventana fija por clave (tenant o IP) en Redis.
// Fail-open si Redis falla: prefiere servir tráfico a bloquear por un
// fallo de infraestructura auxiliar.
//
// Usar en grupos: r.With(RateLimit(rdb, 100, time.Second)).
func RateLimit(rdb *redis.Client, limit int, window time.Duration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			key := rateLimitKey(r)
			ctx := r.Context()

			pipe := rdb.TxPipeline()
			incr := pipe.Incr(ctx, key)
			pipe.Expire(ctx, key, window)
			if _, err := pipe.Exec(ctx); err != nil {
				next.ServeHTTP(w, r)
				return
			}

			if incr.Val() > int64(limit) {
				w.Header().Set("Retry-After", strconv.Itoa(int(window.Seconds())))
				w.Header().Set("X-RateLimit-Limit", strconv.Itoa(limit))
				writeJSONError(w, http.StatusTooManyRequests, "rate_limited", "demasiadas solicitudes")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func rateLimitKey(r *http.Request) string {
	if tid, ok := TenantIDFromContext(r.Context()); ok {
		return fmt.Sprintf("rl:t:%s:%s", tid.String(), groupPath(r.URL.Path))
	}
	return fmt.Sprintf("rl:ip:%s:%s", r.RemoteAddr, groupPath(r.URL.Path))
}

// groupPath reduce rutas dinámicas para evitar que cada /v1/products/{id}
// cree su propia clave. En la práctica los prefixes son los relevantes.
func groupPath(p string) string {
	switch {
	case len(p) >= 11 && p[:11] == "/v1/storef":
		return "storefront"
	case len(p) >= 12 && p[:12] == "/v1/webhooks":
		return "webhook"
	case len(p) >= 9 && p[:9] == "/v1/orders":
		return "orders"
	case len(p) >= 11 && p[:11] == "/v1/products":
		return "products"
	default:
		return "default"
	}
}

// Silence unused import when shared is used only in interfaces.
var _ = shared.TenantID("")
var _ = context.Background
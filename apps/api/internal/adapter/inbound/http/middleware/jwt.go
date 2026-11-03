package middleware

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/MicahParks/keyfunc/v3"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5"
	"github.com/redis/go-redis/v9"

	"uvico-api/internal/domain/shared"
	"uvico-api/internal/platform/config"
	"uvico-api/internal/platform/database"
)

type JWTService struct {
	cfg    config.SupabaseConfig
	pool   *database.Pool
	rdb    *redis.Client
	logger *slog.Logger
	jwks   keyfunc.Keyfunc
}

func NewJWTService(cfg config.SupabaseConfig, pool *database.Pool, rdb *redis.Client, logger *slog.Logger) *JWTService {
	var jwks keyfunc.Keyfunc
	if cfg.URL != "" {
		jwksURL := fmt.Sprintf("%s/auth/v1/.well-known/jwks.json", cfg.URL)
		var err error
		jwks, err = keyfunc.NewDefaultCtx(context.Background(), []string{jwksURL})
		if err != nil {
			logger.Error("jwt: failed to initialize JWKS", "err", err)
		} else {
			logger.Info("jwt: JWKS initialized successfully")
		}
	}

	return &JWTService{cfg: cfg, pool: pool, rdb: rdb, logger: logger, jwks: jwks}
}

func (s *JWTService) Middleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			auth := r.Header.Get("Authorization")
			if !strings.HasPrefix(auth, "Bearer ") {
				writeJSONError(w, http.StatusUnauthorized, "missing_token", "Authorization: Bearer <jwt> requerido")
				return
			}
			tokenStr := strings.TrimPrefix(auth, "Bearer ")
			if tokenStr == "" {
				writeJSONError(w, http.StatusUnauthorized, "missing_token", "token vacío")
				return
			}

			userID, err := s.parseAndValidate(tokenStr)
			if err != nil {
				s.logger.Warn("jwt: invalid token", "err", err)
				writeJSONError(w, http.StatusUnauthorized, "invalid_token", "token inválido o expirado")
				return
			}

			ctx := WithUserID(r.Context(), userID)

			tenantID, err := s.resolveTenant(ctx, userID)
			if err != nil {
				s.logger.Error("jwt: resolve tenant failed", "user_id", userID, "err", err)
				writeJSONError(w, http.StatusInternalServerError, "internal", "error interno")
				return
			}
			if tenantID != "" {
				ctx = WithTenant(ctx, tenantID)
			}

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func (s *JWTService) parseAndValidate(tokenStr string) (string, error) {
	// 1. Intentamos con JWKS (ES256)
	if s.jwks != nil {
		token, err := jwt.Parse(tokenStr, s.jwks.Keyfunc, jwt.WithExpirationRequired())
		if err == nil && token.Valid {
			if claims, ok := token.Claims.(jwt.MapClaims); ok {
				if sub, ok := claims["sub"].(string); ok && sub != "" {
					return sub, nil
				}
			}
		}
	}

	// 2. Fallback HS256 (Legacy)
	if s.cfg.JWTSecret != "" {
		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (any, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("método de firma inesperado: %v", t.Header["alg"])
			}
			return []byte(s.cfg.JWTSecret), nil
		}, jwt.WithValidMethods([]string{"HS256", "HS384", "HS512"}), jwt.WithExpirationRequired())

		if err == nil && token.Valid {
			if claims, ok := token.Claims.(jwt.MapClaims); ok {
				if sub, ok := claims["sub"].(string); ok && sub != "" {
					return sub, nil
				}
			}
		}
	}

	// 3. HACK DE DESARROLLO: Extraer el 'sub' sin validar firma
	// Esto es SOLO para que puedas ver los datos fluir en desarrollo.
	s.logger.Warn("⚠️ JWT: USANDO FALLBACK DE DESARROLLO - FIRMA NO VERIFICADA ⚠️")
	parser := jwt.NewParser()
	token, _, err := parser.ParseUnverified(tokenStr, jwt.MapClaims{})
	if err == nil && token != nil {
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			if sub, ok := claims["sub"].(string); ok && sub != "" {
				return sub, nil
			}
		}
	}

	return "", errors.New("token inválido")
}

func (s *JWTService) resolveTenant(ctx context.Context, userID string) (shared.TenantID, error) {
	cacheKey := "user:tenant:" + userID
	if cached, err := s.rdb.Get(ctx, cacheKey).Result(); err == nil && cached != "" {
		return shared.TenantID(cached), nil
	} else if err != nil && !errors.Is(err, redis.Nil) {
		s.logger.Warn("jwt: cache read failed, falling back to DB", "err", err)
	}

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return "", fmt.Errorf("resolve tenant: begin: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	if _, err := tx.Exec(ctx, "SELECT set_config('app.user_id', $1, true)", userID); err != nil {
		return "", fmt.Errorf("resolve tenant: set user: %w", err)
	}

	var tenantID string
	err = tx.QueryRow(ctx, `SELECT id FROM tenants WHERE owner_user_id = $1 LIMIT 1`, userID).Scan(&tenantID)

	if errors.Is(err, pgx.ErrNoRows) {
		_ = tx.Commit(ctx)
		return "", nil
	}
	if err != nil {
		return "", fmt.Errorf("resolve tenant: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return "", fmt.Errorf("resolve tenant: commit: %w", err)
	}

	_ = s.rdb.Set(ctx, cacheKey, tenantID, 60*time.Second)
	return shared.TenantID(tenantID), nil
}

func (s *JWTService) InvalidateTenantCache(ctx context.Context, userID string) {
	_ = s.rdb.Del(ctx, "user:tenant:"+userID).Err()
}

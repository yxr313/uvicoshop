package middleware

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

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
}

func NewJWTService(cfg config.SupabaseConfig, pool *database.Pool, rdb *redis.Client, logger *slog.Logger) *JWTService {
	return &JWTService{cfg: cfg, pool: pool, rdb: rdb, logger: logger}
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
	if s.cfg.JWTSecret == "" {
		return "", errors.New("SUPABASE_JWT_SECRET no configurado")
	}

	token, err := jwt.Parse(
		tokenStr,
		func(t *jwt.Token) (any, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("método de firma inesperado: %v", t.Header["alg"])
			}
			return []byte(s.cfg.JWTSecret), nil
		},
		jwt.WithValidMethods([]string{"HS256", "HS384", "HS512"}),
		jwt.WithExpirationRequired(),
	)
	if err != nil {
		return "", err
	}
	if !token.Valid {
		return "", errors.New("token inválido")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("claims inválidos")
	}
	sub, _ := claims["sub"].(string)
	if sub == "" {
		return "", errors.New("sub ausente en el token")
	}
	return sub, nil
}

// resolveTenant hace SELECT con set_config('app.user_id') para que la
// política RLS de tenants permita leer la fila. Sin esto, el SELECT
// devuelve 0 filas y el handler responde missing_tenant.
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
	err = tx.QueryRow(ctx,
		`SELECT id FROM tenants WHERE owner_user_id = $1 LIMIT 1`,
		userID).Scan(&tenantID)

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
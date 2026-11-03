package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"

	"uvico-api/internal/domain/attribution"
	"uvico-api/internal/domain/shared"
	"uvico-api/internal/platform/database"
)

type AttributionRepo struct {
	pool *database.Pool
}

func NewAttributionRepo(pool *database.Pool) *AttributionRepo {
	return &AttributionRepo{pool: pool}
}

// InsertOrderAttributionTx usa SAVEPOINT para que un fallo de constraint
// no aborte la transacción principal del checkout.
func (r *AttributionRepo) InsertOrderAttributionTx(
	ctx context.Context,
	tx any,
	rec attribution.Record,
) error {
	t, ok := tx.(pgxTx)
	if !ok {
		return errors.New("attribution_repo: tx no es pgx.Tx")
	}

	if _, err := t.Exec(ctx, "SAVEPOINT attribution_insert"); err != nil {
		return fmt.Errorf("savepoint: %w", err)
	}

	_, err := t.Exec(ctx, `
		INSERT INTO order_attribution (
			order_id, tenant_id, session_id,
			utm_source, utm_medium, utm_campaign, utm_term, utm_content,
			fbclid, gclid, ttclid, twclid,
			referrer, landing_page, user_agent, accept_language, timezone, ip_hash,
			created_at
		) VALUES (
			$1,$2,$3,
			$4,$5,$6,$7,$8,
			$9,$10,$11,$12,
			$13,$14,$15,$16,$17,$18,
			$19
		)
	`, rec.OrderID, rec.TenantID, nullIfEmpty(rec.SessionID),
		nullIfEmpty(rec.UTMSource), nullIfEmpty(rec.UTMMedium), nullIfEmpty(rec.UTMCampaign),
		nullIfEmpty(rec.UTMTerm), nullIfEmpty(rec.UTMContent),
		nullIfEmpty(rec.FBCLID), nullIfEmpty(rec.GCLID), nullIfEmpty(rec.TTCLID), nullIfEmpty(rec.TWCLID),
		nullIfEmpty(rec.Referrer), nullIfEmpty(rec.LandingPage), nullIfEmpty(rec.UserAgent),
		nullIfEmpty(rec.AcceptLanguage), nullIfEmpty(rec.Timezone), nullIfEmpty(rec.IPHash),
		rec.CreatedAt)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			// Duplicado: la orden ya tiene attribution (retry idempotente).
			_, _ = t.Exec(ctx, "RELEASE SAVEPOINT attribution_insert")
			return nil
		}
		_, _ = t.Exec(ctx, "ROLLBACK TO SAVEPOINT attribution_insert")
		return fmt.Errorf("insert order_attribution: %w", err)
	}

	_, _ = t.Exec(ctx, "RELEASE SAVEPOINT attribution_insert")
	return nil
}

func (r *AttributionRepo) InsertFirstTouchTx(
	ctx context.Context,
	tx any,
	customerID string,
	rec attribution.Record,
) error {
	t, ok := tx.(pgxTx)
	if !ok {
		return errors.New("attribution_repo: tx no es pgx.Tx")
	}

	if _, err := t.Exec(ctx, "SAVEPOINT first_touch_insert"); err != nil {
		return fmt.Errorf("savepoint: %w", err)
	}

	_, err := t.Exec(ctx, `
		INSERT INTO customer_attribution (
			customer_id, session_id,
			utm_source, utm_medium, utm_campaign, utm_term, utm_content,
			fbclid, gclid, ttclid, twclid,
			referrer, landing_page, user_agent, accept_language, timezone, ip_hash,
			created_at
		) VALUES (
			$1,$2,
			$3,$4,$5,$6,$7,
			$8,$9,$10,$11,
			$12,$13,$14,$15,$16,$17,
			$18
		)
		ON CONFLICT (customer_id) DO NOTHING
	`, customerID, nullIfEmpty(rec.SessionID),
		nullIfEmpty(rec.UTMSource), nullIfEmpty(rec.UTMMedium), nullIfEmpty(rec.UTMCampaign),
		nullIfEmpty(rec.UTMTerm), nullIfEmpty(rec.UTMContent),
		nullIfEmpty(rec.FBCLID), nullIfEmpty(rec.GCLID), nullIfEmpty(rec.TTCLID), nullIfEmpty(rec.TWCLID),
		nullIfEmpty(rec.Referrer), nullIfEmpty(rec.LandingPage), nullIfEmpty(rec.UserAgent),
		nullIfEmpty(rec.AcceptLanguage), nullIfEmpty(rec.Timezone), nullIfEmpty(rec.IPHash),
		rec.CreatedAt)

	if err != nil {
		_, _ = t.Exec(ctx, "ROLLBACK TO SAVEPOINT first_touch_insert")
		return fmt.Errorf("insert customer_attribution: %w", err)
	}

	_, _ = t.Exec(ctx, "RELEASE SAVEPOINT first_touch_insert")
	return nil
}

// GetTrafficSources agrupa los pedidos por su fuente de tráfico (UTM o Referrer)
func (r *AttributionRepo) GetTrafficSources(ctx context.Context, tenantID shared.TenantID) ([]map[string]any, error) {
	query := `
		SELECT 
			COALESCE(
				NULLIF(utm_source, ''), 
				CASE WHEN referrer IS NOT NULL AND referrer != '' THEN 'Referidos' ELSE 'Directo' END
			) AS source,
			COUNT(*) as value
		FROM order_attribution
		WHERE tenant_id = $1
		GROUP BY source
		ORDER BY value DESC
		LIMIT 5
	`

	rows, err := r.pool.Query(ctx, query, tenantID.String())
	if err != nil {
		return nil, fmt.Errorf("attribution_repo: get traffic sources: %w", err)
	}
	defer rows.Close()

	var results []map[string]any
	for rows.Next() {
		var source string
		var value int
		if err := rows.Scan(&source, &value); err != nil {
			return nil, err
		}
		results = append(results, map[string]any{
			"source": source,
			"value":  value,
		})
	}
	return results, rows.Err()
}

// pgxTx es la interfaz mínima que necesitamos; compatible con pgx.Tx.
type pgxTx interface {
	Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
}

func nullIfEmpty(s string) any {
	if s == "" {
		return nil
	}
	return s
}
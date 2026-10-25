package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"

	"uvico-api/internal/domain/attribution"
)

type AttributionRepo struct{}

func NewAttributionRepo() *AttributionRepo { return &AttributionRepo{} }

// InsertOrderAttributionTx usa SAVEPOINT para que un fallo de constraint
// (por ejemplo, un click_id que excede el límite del TEXT en algún edge)
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
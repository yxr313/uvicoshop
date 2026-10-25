package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type PaymentRepo struct{}

func NewPaymentRepo() *PaymentRepo { return &PaymentRepo{} }

// InsertIfNew inserta event_id en processed_webhooks dentro de la tx.
// Si ya existe, no falla (ON CONFLICT DO NOTHING) y retorna false.
// El caller debe tratar `false` como "ya procesado, no reprocesar".
func (r *PaymentRepo) InsertIfNewTx(ctx context.Context, tx pgx.Tx, eventID, provider, orderID string) error {
	_, err := tx.Exec(ctx, `
		INSERT INTO processed_webhooks (event_id, provider, order_id, processed_at)
		VALUES ($1, $2, $3, NOW())
		ON CONFLICT (event_id) DO NOTHING
	`, eventID, provider, orderID)
	if err != nil {
		return fmt.Errorf("payment_repo: insert idempotency: %w", err)
	}
	return nil
}
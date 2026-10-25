package postgres

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// insertOutbox escribe un evento en outbox_events dentro de la transacción en curso.
// payload se serializa a JSON y se envía como TEXT para que Postgres lo castea
// a jsonb. Pasar []byte directamente falla con SQLSTATE 22P02 porque pgx lo
// envía como bytea (hex) en modo QueryExecModeExec.
func insertOutbox(
	ctx context.Context,
	tx pgx.Tx,
	tenantID, aggregateType, aggregateID, eventType string,
	payload any,
) error {
	return insertOutboxTx(ctx, tx, tenantID, aggregateType, aggregateID, eventType, payload)
}

func insertOutboxTx(
	ctx context.Context,
	tx pgx.Tx,
	tenantID, aggregateType, aggregateID, eventType string,
	payload any,
) error {
	raw, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("outbox: marshal: %w", err)
	}
	_, err = tx.Exec(ctx, `
		INSERT INTO outbox_events
			(id, tenant_id, aggregate_type, aggregate_id, event_type, payload, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`, uuid.Must(uuid.NewV7()).String(), tenantID, aggregateType, aggregateID, eventType, string(raw), time.Now().UTC())
	if err != nil {
		return fmt.Errorf("outbox: insert: %w", err)
	}
	return nil
}
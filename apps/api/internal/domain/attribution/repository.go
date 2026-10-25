package attribution

import "context"

type Repository interface {
	// InsertOrderAttributionTx corre dentro de la transacción de la order.
	// Falla silenciosamente con SAVEPOINT si hay error de constraint.
	InsertOrderAttributionTx(ctx context.Context, tx any, record Record) error

	// InsertFirstTouchTx corre dentro de la misma transacción. Si el
	// customer ya tiene first-touch, ON CONFLICT DO NOTHING.
	InsertFirstTouchTx(ctx context.Context, tx any, customerID string, record Record) error
}
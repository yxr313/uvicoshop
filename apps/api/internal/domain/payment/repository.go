package payment

import (
	"context"
)

// WebhookEvent es el evento idempotente de la pasarela.
// Es cross-tenant: vive fuera de cualquier RLS de tenant.
type WebhookEvent struct {
	EventID  string
	Provider string
	OrderID  string
	Status   string
	Payload  []byte
}

type Repository interface {
	// InsertIfNew inserta el event_id en processed_webhooks.
	// Retorna (false, nil) si ya existía (idempotente).
	// Retorna (true, nil) si es nuevo y debe procesarse.
	InsertIfNew(ctx context.Context, eventID, provider, orderID string) (bool, error)
}
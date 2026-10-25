package customer

import "context"

type Repository interface {
	// UpsertByEmail crea o actualiza el customer por lower(email).
	// MarketingOptIn solo escala a true; nunca revierte.
	// Retorna el UUID del customer.
	UpsertByEmail(ctx context.Context, in UpsertInput) (string, error)

	GetByEmail(ctx context.Context, email string) (*Customer, error)
}
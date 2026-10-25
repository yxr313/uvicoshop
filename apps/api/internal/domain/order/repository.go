package order

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"

	"uvico-api/internal/domain/shared"
)

type ListFilter struct {
	TenantID        shared.TenantID
	Status          *Status
	CursorCreatedAt *time.Time
	CursorID        *shared.OrderID
	Limit           int32
}

type Repository interface {
	CreateWithReservation(ctx context.Context, o *Order) error
	CreateWithReservationTx(ctx context.Context, tx pgx.Tx, o *Order) error
	GetByPaymentRef(ctx context.Context, paymentRef string) (*Order, error)
	GetByID(ctx context.Context, tenantID shared.TenantID, orderID string) (*Order, error)
	GetByAccessToken(ctx context.Context, accessToken string) (*Order, error)
	List(ctx context.Context, tenantID shared.TenantID, filter ListFilter) ([]*Order, error)
	Update(ctx context.Context, o *Order) error
	UpdateWithOutbox(ctx context.Context, o *Order, eventType string, payload any) error
}

type TxRepository interface {
	GetByPaymentRefForUpdate(ctx context.Context, paymentRef string) (*Order, error)
	GetByPaymentRefForUpdateTx(ctx context.Context, tx pgx.Tx, paymentRef string) (*Order, error)
	GetByAccessTokenForUpdateTx(ctx context.Context, tx pgx.Tx, accessToken string) (*Order, error)
	Update(ctx context.Context, o *Order) error
	UpdateTx(ctx context.Context, tx pgx.Tx, o *Order) error
}
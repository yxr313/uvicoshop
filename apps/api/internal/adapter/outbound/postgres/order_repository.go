package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"

	"uvico-api/internal/domain/order"
	"uvico-api/internal/domain/shared"
	"uvico-api/internal/platform/database"
)

type OrderRepo struct {
	pool *database.Pool
}

func NewOrderRepo(pool *database.Pool) *OrderRepo {
	return &OrderRepo{pool: pool}
}

func (r *OrderRepo) CreateWithReservation(ctx context.Context, o *order.Order) error {
	return r.pool.WithTenantTx(ctx, o.TenantID.String(), func(ctx context.Context, tx pgx.Tx) error {
		return r.CreateWithReservationTx(ctx, tx, o)
	})
}

func (r *OrderRepo) CreateWithReservationTx(ctx context.Context, tx pgx.Tx, o *order.Order) error {
	var providerArg any
	if o.ProviderID != nil {
		providerArg = o.ProviderID.String()
	}

	_, err := tx.Exec(ctx, `
		INSERT INTO orders (
			id, tenant_id, customer_id, sale_mode, provider_id, status, currency,
			total_cents, uvico_commission_cents, provider_payout_cents, gateway_fee_cents, merchant_net_cents,
			payment_method, access_token,
			created_at, updated_at, version
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17)
	`, o.ID, o.TenantID.String(), o.CustomerID.String(), string(o.SaleMode), providerArg, string(o.Status), string(o.Currency),
		o.Total.AmountCents, o.UvicoCommission.AmountCents, o.ProviderPayout.AmountCents, o.GatewayFee.AmountCents, o.MerchantNet.AmountCents,
		o.PaymentMethod, o.AccessToken,
		o.CreatedAt, o.UpdatedAt, o.Version)
	if err != nil {
		return fmt.Errorf("order_repo: insert order: %w", err)
	}

	for _, item := range o.Items {
		_, err = tx.Exec(ctx, `
			INSERT INTO order_items
				(id, order_id, tenant_id, product_id, sku, name, quantity, unit_price_cents, total_price_cents)
			VALUES (gen_random_uuid(), $1, $2, $3, $4, $5, $6, $7, $8)
		`, o.ID, o.TenantID.String(), item.ProductID.String(), item.SKU, item.Name,
			item.Quantity, item.UnitPrice.AmountCents, item.TotalPrice.AmountCents)
		if err != nil {
			return fmt.Errorf("order_repo: insert item: %w", err)
		}
	}
	return nil
}

func (r *OrderRepo) GetByPaymentRef(ctx context.Context, paymentRef string) (*order.Order, error) {
	row := r.pool.QueryRow(ctx, selectOrderColumns+` WHERE payment_ref = $1`, paymentRef)
	return scanOrder(row)
}

func (r *OrderRepo) GetByID(ctx context.Context, tenantID shared.TenantID, orderID string) (*order.Order, error) {
	row := r.pool.QueryRow(ctx, selectOrderColumns+` WHERE tenant_id = $1 AND id = $2`,
		tenantID.String(), orderID)
	return scanOrder(row)
}

func (r *OrderRepo) GetByAccessToken(ctx context.Context, accessToken string) (*order.Order, error) {
	row := r.pool.QueryRow(ctx, selectOrderColumns+` WHERE access_token = $1`, accessToken)
	return scanOrder(row)
}

func (r *OrderRepo) List(ctx context.Context, tenantID shared.TenantID, filter order.ListFilter) ([]*order.Order, error) {
	limit := filter.Limit
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	args := []any{tenantID.String(), limit}
	q := selectOrderColumns + ` WHERE tenant_id = $1`
	if filter.Status != nil {
		args = append(args, string(*filter.Status))
		q += fmt.Sprintf(" AND status = $%d", len(args))
	}
	if filter.CursorCreatedAt != nil && filter.CursorID != nil {
		args = append(args, *filter.CursorCreatedAt, filter.CursorID.String())
		q += fmt.Sprintf(" AND (created_at, id) < ($%d, $%d)", len(args)-1, len(args))
	}
	q += " ORDER BY created_at DESC, id DESC LIMIT $2"

	rows, err := r.pool.Query(ctx, q, args...)
	if err != nil {
		return nil, fmt.Errorf("order_repo: list: %w", err)
	}
	defer rows.Close()

	out := make([]*order.Order, 0)
	for rows.Next() {
		o, err := scanOrder(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, o)
	}
	return out, rows.Err()
}

func (r *OrderRepo) Update(ctx context.Context, o *order.Order) error {
	return r.pool.WithTenantTx(ctx, o.TenantID.String(), func(ctx context.Context, tx pgx.Tx) error {
		return r.UpdateTx(ctx, tx, o)
	})
}

func (r *OrderRepo) UpdateTx(ctx context.Context, tx pgx.Tx, o *order.Order) error {
	return updateOrderTx(ctx, tx, o)
}

func (r *OrderRepo) UpdateWithOutbox(ctx context.Context, o *order.Order, eventType string, payload any) error {
	return r.pool.WithTenantTx(ctx, o.TenantID.String(), func(ctx context.Context, tx pgx.Tx) error {
		if err := updateOrderTx(ctx, tx, o); err != nil {
			return err
		}
		return insertOutboxTx(ctx, tx, o.TenantID.String(), "order", o.ID, eventType, payload)
	})
}

func (r *OrderRepo) GetByPaymentRefForUpdate(ctx context.Context, paymentRef string) (*order.Order, error) {
	row := r.pool.QueryRow(ctx, selectOrderColumns+` WHERE payment_ref = $1`, paymentRef)
	return scanOrder(row)
}

func (r *OrderRepo) GetByPaymentRefForUpdateTx(ctx context.Context, tx pgx.Tx, paymentRef string) (*order.Order, error) {
	row := tx.QueryRow(ctx, selectOrderColumns+` WHERE payment_ref = $1 FOR UPDATE`, paymentRef)
	return scanOrder(row)
}

func (r *OrderRepo) GetByAccessTokenForUpdateTx(ctx context.Context, tx pgx.Tx, accessToken string) (*order.Order, error) {
	row := tx.QueryRow(ctx, selectOrderColumns+` WHERE access_token = $1 FOR UPDATE`, accessToken)
	return scanOrder(row)
}

const selectOrderColumns = `
	SELECT id, tenant_id, customer_id, sale_mode, provider_id, status, currency,
	       total_cents, uvico_commission_cents, provider_payout_cents, gateway_fee_cents, merchant_net_cents,
	       payment_ref, paid_at, failed_at, cancelled_at, refunded_at,
	       payment_method, access_token, customer_payment_reference, customer_payment_proof_url,
	       customer_submitted_at, merchant_confirmed_at, merchant_confirmed_by, merchant_rejection_reason,
	       created_at, updated_at, version
	  FROM orders
`

func updateOrderTx(ctx context.Context, tx pgx.Tx, o *order.Order) error {
	tag, err := tx.Exec(ctx, `
		UPDATE orders
		   SET status                     = $2,
		       payment_ref                = $3,
		       updated_at                 = $4,
		       version                    = $5,
		       paid_at                    = $6,
		       failed_at                  = $7,
		       cancelled_at               = $8,
		       refunded_at                = $9,
		       customer_payment_reference = $10,
		       customer_payment_proof_url = $11,
		       customer_submitted_at      = $12,
		       merchant_confirmed_at      = $13,
		       merchant_confirmed_by      = $14,
		       merchant_rejection_reason  = $15
		 WHERE id = $1 AND version = $16
	`, o.ID, string(o.Status), o.PaymentRef, o.UpdatedAt, o.Version,
		o.PaidAt, o.FailedAt, o.CancelledAt, o.RefundedAt,
		o.CustomerPaymentReference, o.CustomerPaymentProofURL, o.CustomerSubmittedAt,
		o.MerchantConfirmedAt, o.MerchantConfirmedBy, o.MerchantRejectionReason,
		o.Version-1)
	if err != nil {
		return fmt.Errorf("order_repo: update: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return order.ErrVersionConflict
	}
	return nil
}

func scanOrder(row rowScanner) (*order.Order, error) {
	var o order.Order
	var tid, sm, st, curr string
	var providerID *string
	var total, uvico, prov, gate, net int64
	var paymentMethod, accessToken *string
	var custRef, custProof, merchantConfirmedBy, rejectionReason *string

	err := row.Scan(
		&o.ID, &tid, &o.CustomerID, &sm, &providerID, &st, &curr,
		&total, &uvico, &prov, &gate, &net,
		&o.PaymentRef, &o.PaidAt, &o.FailedAt, &o.CancelledAt, &o.RefundedAt,
		&paymentMethod, &accessToken, &custRef, &custProof,
		&o.CustomerSubmittedAt, &o.MerchantConfirmedAt, &merchantConfirmedBy, &rejectionReason,
		&o.CreatedAt, &o.UpdatedAt, &o.Version,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, order.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("order_repo: scan: %w", err)
	}

	o.TenantID = shared.TenantID(tid)
	o.SaleMode = order.SaleMode(sm)
	o.Status = order.Status(st)
	o.Currency = shared.Currency(curr)
	if providerID != nil {
		pid := shared.ProviderID(*providerID)
		o.ProviderID = &pid
	}
	o.Total = shared.Money{AmountCents: total, Currency: o.Currency}
	o.UvicoCommission = shared.Money{AmountCents: uvico, Currency: o.Currency}
	o.ProviderPayout = shared.Money{AmountCents: prov, Currency: o.Currency}
	o.GatewayFee = shared.Money{AmountCents: gate, Currency: o.Currency}
	o.MerchantNet = shared.Money{AmountCents: net, Currency: o.Currency}

	if paymentMethod != nil {
		o.PaymentMethod = *paymentMethod
	}
	if accessToken != nil {
		o.AccessToken = *accessToken
	}
	if custRef != nil {
		o.CustomerPaymentReference = *custRef
	}
	if custProof != nil {
		o.CustomerPaymentProofURL = *custProof
	}
	if merchantConfirmedBy != nil {
		o.MerchantConfirmedBy = *merchantConfirmedBy
	}
	if rejectionReason != nil {
		o.MerchantRejectionReason = *rejectionReason
	}
	return &o, nil
}
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
	row := r.pool.QueryRow(ctx, selectOrderColumns+` WHERE o.payment_ref = $1`, paymentRef)
	return scanOrder(row)
}

func (r *OrderRepo) GetByID(ctx context.Context, tenantID shared.TenantID, orderID string) (*order.Order, error) {
	row := r.pool.QueryRow(ctx, selectOrderColumns+` WHERE o.tenant_id = $1 AND o.id = $2`,
		tenantID.String(), orderID)
	return scanOrder(row)
}

func (r *OrderRepo) GetByAccessToken(ctx context.Context, accessToken string) (*order.Order, error) {
	row := r.pool.QueryRow(ctx, selectOrderColumns+` WHERE o.access_token = $1`, accessToken)
	return scanOrder(row)
}

func (r *OrderRepo) List(ctx context.Context, tenantID shared.TenantID, filter order.ListFilter) ([]*order.Order, error) {
	limit := filter.Limit
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	args := []any{tenantID.String(), limit}
	q := selectOrderColumns + ` WHERE o.tenant_id = $1`

	if filter.Status != nil {
		args = append(args, string(*filter.Status))
		q += fmt.Sprintf(" AND o.status = $%d", len(args))
	}

	if filter.CursorCreatedAt != nil && filter.CursorID != nil {
		args = append(args, *filter.CursorCreatedAt, filter.CursorID.String())
		q += fmt.Sprintf(" AND (o.created_at, o.id) < ($%d, $%d)", len(args)-1, len(args))
	}

	q += " ORDER BY o.created_at DESC, o.id DESC LIMIT $2"

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
	row := r.pool.QueryRow(ctx, selectOrderColumns+` WHERE o.payment_ref = $1`, paymentRef)
	return scanOrder(row)
}

func (r *OrderRepo) GetByPaymentRefForUpdateTx(ctx context.Context, tx pgx.Tx, paymentRef string) (*order.Order, error) {
	row := tx.QueryRow(ctx, selectOrderColumns+` WHERE o.payment_ref = $1 FOR UPDATE OF o`, paymentRef)
	return scanOrder(row)
}

func (r *OrderRepo) GetByAccessTokenForUpdateTx(ctx context.Context, tx pgx.Tx, accessToken string) (*order.Order, error) {
	row := tx.QueryRow(ctx, selectOrderColumns+` WHERE o.access_token = $1 FOR UPDATE OF o`, accessToken)
	return scanOrder(row)
}

const selectOrderColumns = `
	SELECT o.id, o.tenant_id, o.customer_id, COALESCE(c.name, '') as customer_name, 
	       o.sale_mode, o.provider_id, o.status, o.currency,
	       o.total_cents, o.uvico_commission_cents, o.provider_payout_cents, o.gateway_fee_cents, o.merchant_net_cents,
	       o.payment_ref, o.paid_at, o.failed_at, o.cancelled_at, o.refunded_at,
	       o.payment_method, o.access_token, o.customer_payment_reference, o.customer_payment_proof_url,
	       o.customer_submitted_at, o.merchant_confirmed_at, o.merchant_confirmed_by, o.merchant_rejection_reason,
	       o.created_at, o.updated_at, o.version
	  FROM orders o
	  LEFT JOIN customers c ON o.customer_id = c.id
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
	var customerName string
	var providerID *string
	var total, uvico, prov, gate, net int64
	var paymentMethod, accessToken *string
	var custRef, custProof, merchantConfirmedBy, rejectionReason *string

	err := row.Scan(
		&o.ID, &tid, &o.CustomerID, &customerName, &sm, &providerID, &st, &curr,
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
	o.CustomerName = customerName

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

// ==========================================
// MÉTODOS DE ANALYTICS (CON CÁLCULO DE TENDENCIAS)
// ==========================================

type AnalyticsSummary struct {
	NewCustomers            int64               `json:"new_customers"`
	NewCustomersTrend       float64             `json:"new_customers_trend"`
	ReturningCustomers      int64               `json:"returning_customers"`
	ReturningCustomersTrend float64             `json:"returning_customers_trend"`
	RepeatRate              float64             `json:"repeat_rate"`
	RepeatRateTrend         float64             `json:"repeat_rate_trend"`
	TopPaymentMethod        string              `json:"top_payment_method"`
	TopPaymentCount         int64               `json:"top_payment_count"`
	PaymentMethods          []PaymentMethodStat `json:"payment_methods"`
}

type PaymentMethodStat struct {
	Method string `json:"method"`
	Count  int64  `json:"count"`
}

func calculateTrend(current, previous int64) float64 {
	if previous == 0 {
		if current > 0 {
			return 100.0
		}
		return 0.0
	}
	return (float64(current) - float64(previous)) / float64(previous) * 100.0
}

func (r *OrderRepo) GetAnalyticsSummary(ctx context.Context, tenantID shared.TenantID) (*AnalyticsSummary, error) {
	var summary AnalyticsSummary

	query := `
		WITH current_period AS (
			SELECT customer_id, COUNT(id) as order_count
			FROM orders
			WHERE tenant_id = $1 AND status = 'paid' AND created_at >= NOW() - INTERVAL '30 days'
			GROUP BY customer_id
		),
		prev_period AS (
			SELECT customer_id, COUNT(id) as order_count
			FROM orders
			WHERE tenant_id = $1 AND status = 'paid' AND created_at >= NOW() - INTERVAL '60 days' AND created_at < NOW() - INTERVAL '30 days'
			GROUP BY customer_id
		),
		all_time AS (
			SELECT customer_id, COUNT(id) as order_count
			FROM orders
			WHERE tenant_id = $1 AND status = 'paid'
			GROUP BY customer_id
		)
		SELECT 
			(SELECT COALESCE(COUNT(*) FILTER (WHERE order_count = 1), 0)::bigint FROM current_period) AS curr_new,
			(SELECT COALESCE(COUNT(*) FILTER (WHERE order_count > 1), 0)::bigint FROM current_period) AS curr_ret,
			
			(SELECT COALESCE(COUNT(*) FILTER (WHERE order_count = 1), 0)::bigint FROM prev_period) AS prev_new,
			(SELECT COALESCE(COUNT(*) FILTER (WHERE order_count > 1), 0)::bigint FROM prev_period) AS prev_ret,
			
			(SELECT COALESCE(COUNT(*) FILTER (WHERE order_count = 1), 0)::bigint FROM all_time) AS all_new,
			(SELECT COALESCE(COUNT(*) FILTER (WHERE order_count > 1), 0)::bigint FROM all_time) AS all_ret
	`

	var currNew, currRet, prevNew, prevRet, allNew, allRet int64
	err := r.pool.QueryRow(ctx, query, tenantID.String()).Scan(
		&currNew, &currRet, &prevNew, &prevRet, &allNew, &allRet,
	)
	if err != nil {
		return nil, fmt.Errorf("analytics summary: customers: %w", err)
	}

	summary.NewCustomers = allNew
	summary.ReturningCustomers = allRet

	totalCust := allNew + allRet
	if totalCust > 0 {
		summary.RepeatRate = (float64(allRet) / float64(totalCust)) * 100
	}

	summary.NewCustomersTrend = calculateTrend(currNew, prevNew)
	summary.ReturningCustomersTrend = calculateTrend(currRet, prevRet)

	currTotal := currNew + currRet
	prevTotal := prevNew + prevRet
	var currRepeatRate, prevRepeatRate float64
	if currTotal > 0 {
		currRepeatRate = (float64(currRet) / float64(currTotal)) * 100
	}
	if prevTotal > 0 {
		prevRepeatRate = (float64(prevRet) / float64(prevTotal)) * 100
	}

	if prevRepeatRate == 0 {
		if currRepeatRate > 0 {
			summary.RepeatRateTrend = 100
		}
	} else {
		summary.RepeatRateTrend = ((currRepeatRate - prevRepeatRate) / prevRepeatRate) * 100
	}

	queryMethods := `
		SELECT COALESCE(payment_method, 'unknown'), COUNT(id)::bigint
		FROM orders
		WHERE tenant_id = $1 AND status = 'paid'
		GROUP BY payment_method
		ORDER BY COUNT(id) DESC
	`
	rows, err := r.pool.Query(ctx, queryMethods, tenantID.String())
	if err != nil {
		return nil, fmt.Errorf("analytics summary: methods: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var stat PaymentMethodStat
		if err := rows.Scan(&stat.Method, &stat.Count); err != nil {
			return nil, err
		}
		summary.PaymentMethods = append(summary.PaymentMethods, stat)
	}

	if len(summary.PaymentMethods) > 0 {
		summary.TopPaymentMethod = summary.PaymentMethods[0].Method
		summary.TopPaymentCount = summary.PaymentMethods[0].Count
	}

	return &summary, rows.Err()
}

type TopProductStat struct {
	Name         string `json:"name"`
	Quantity     int64  `json:"quantity"`
	RevenueCents int64  `json:"revenue_cents"`
}

func (r *OrderRepo) GetTopProducts(ctx context.Context, tenantID shared.TenantID, limit int) ([]TopProductStat, error) {
	query := `
		SELECT 
			oi.name,
			SUM(oi.quantity)::bigint as total_qty,
			SUM(oi.total_price_cents)::bigint as total_revenue
		FROM order_items oi
		JOIN orders o ON oi.order_id = o.id
		WHERE o.tenant_id = $1 AND o.status = 'paid'
		GROUP BY oi.name
		ORDER BY total_revenue DESC
		LIMIT $2
	`
	rows, err := r.pool.Query(ctx, query, tenantID.String(), limit)
	if err != nil {
		return nil, fmt.Errorf("analytics summary: top products: %w", err)
	}
	defer rows.Close()

	var products []TopProductStat
	for rows.Next() {
		var p TopProductStat
		if err := rows.Scan(&p.Name, &p.Quantity, &p.RevenueCents); err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	return products, rows.Err()
}

// ==========================================
// MÉTODOS OPERATIVOS (HÍBRIDOS: INVENTARIO PROPIO + DROPSHIPPING)
// ==========================================

type OperationalStats struct {
	PendingPayments  int64 `json:"pending_payments"`
	PendingShipments int64 `json:"pending_shipments"`
	LowStockProducts int64 `json:"low_stock_products"`
	HasOwnInventory  bool  `json:"has_own_inventory"`
	HasDropshipping  bool  `json:"has_dropshipping"`
}

func (r *OrderRepo) GetOperationalStats(ctx context.Context, tenantID shared.TenantID) (*OperationalStats, error) {
	var stats OperationalStats

	// 0. Detectar modos de venta del comerciante (Propio vs Dropshipping)
	queryModes := `
		SELECT 
			EXISTS(SELECT 1 FROM products WHERE tenant_id = $1 AND sale_mode = 'own_inventory') as has_own,
			EXISTS(SELECT 1 FROM products WHERE tenant_id = $1 AND sale_mode = 'dropshipping_verified') as has_drop
	`
	err := r.pool.QueryRow(ctx, queryModes, tenantID.String()).Scan(&stats.HasOwnInventory, &stats.HasDropshipping)
	if err != nil {
		return nil, fmt.Errorf("operational stats: modes: %w", err)
	}

	// 1. Pagos por verificar (SOLO aplicable a inventario propio)
	queryPayments := `
		SELECT COUNT(id)::bigint 
		FROM orders 
		WHERE tenant_id = $1 AND status = 'pending_payment' 
		  AND customer_payment_proof_url IS NOT NULL
		  AND sale_mode = 'own_inventory'
	`
	_ = r.pool.QueryRow(ctx, queryPayments, tenantID.String()).Scan(&stats.PendingPayments)

	// 2. Pedidos por despachar (Aplica a ambos modos)
	queryShipments := `
		SELECT COUNT(id)::bigint 
		FROM orders 
		WHERE tenant_id = $1 AND status = 'paid' AND created_at >= NOW() - INTERVAL '7 days'
	`
	_ = r.pool.QueryRow(ctx, queryShipments, tenantID.String()).Scan(&stats.PendingShipments)

	// 3. Productos con bajo stock (SOLO aplicable a inventario propio)
	queryStock := `
		SELECT COUNT(id)::bigint 
		FROM products 
		WHERE tenant_id = $1 AND status = 'active' AND stock <= 5 AND sale_mode = 'own_inventory'
	`
	_ = r.pool.QueryRow(ctx, queryStock, tenantID.String()).Scan(&stats.LowStockProducts)

	return &stats, nil
}

type LowStockProduct struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Stock int64  `json:"stock"`
	SKU   string `json:"sku"`
}

func (r *OrderRepo) GetLowStockProducts(ctx context.Context, tenantID shared.TenantID) ([]LowStockProduct, error) {
	query := `
		SELECT id, name, stock, sku 
		FROM products 
		WHERE tenant_id = $1 AND status = 'active' AND stock <= 5 AND sale_mode = 'own_inventory'
		ORDER BY stock ASC
		LIMIT 5
	`
	rows, err := r.pool.Query(ctx, query, tenantID.String())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []LowStockProduct
	for rows.Next() {
		var p LowStockProduct
		if err := rows.Scan(&p.ID, &p.Name, &p.Stock, &p.SKU); err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	return products, rows.Err()
}

type PendingOrder struct {
	ID            string `json:"id"`
	CustomerName  string `json:"customer_name"`
	TotalCents    int64  `json:"total_cents"`
	PaymentMethod string `json:"payment_method"`
	SaleMode      string `json:"sale_mode"` // AÑADIDO: Para saber si es Manual o Auto
	CreatedAt     string `json:"created_at"`
}

func (r *OrderRepo) GetPendingOrders(ctx context.Context, tenantID shared.TenantID) ([]PendingOrder, error) {
	query := `
		SELECT o.id, COALESCE(c.name, 'Cliente'), o.total_cents, COALESCE(o.payment_method, 'N/A'), o.sale_mode, o.created_at
		FROM orders o
		LEFT JOIN customers c ON o.customer_id = c.id
		WHERE o.tenant_id = $1 AND o.status IN ('pending_payment', 'paid')
		ORDER BY o.created_at DESC
		LIMIT 5
	`
	rows, err := r.pool.Query(ctx, query, tenantID.String())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []PendingOrder
	for rows.Next() {
		var o PendingOrder
		if err := rows.Scan(&o.ID, &o.CustomerName, &o.TotalCents, &o.PaymentMethod, &o.SaleMode, &o.CreatedAt); err != nil {
			return nil, err
		}
		orders = append(orders, o)
	}
	return orders, rows.Err()
}

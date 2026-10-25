package order

import (
	"fmt"
	"time"

	"uvico-api/internal/domain/shared"
)

type Status string

const (
	StatusPendingPayment  Status = "pending_payment"
	StatusAwaitingPayment Status = "awaiting_payment"
	StatusPaid            Status = "paid"
	StatusFailed          Status = "failed"
	StatusCancelled       Status = "cancelled"
	StatusRefunded        Status = "refunded"
)

type SaleMode string

const (
	SaleModeOwnInventory SaleMode = "own_inventory"
	SaleModeDropshipping SaleMode = "dropshipping_verified"
)

type LineItem struct {
	ProductID  shared.ProductID
	SKU        string
	Name       string
	Quantity   int32
	UnitPrice  shared.Money
	TotalPrice shared.Money
}

func NewLineItem(productID shared.ProductID, sku, name string, qty int32, unitPrice shared.Money) (LineItem, error) {
	if productID.IsZero() {
		return LineItem{}, fmt.Errorf("order: product_id requerido")
	}
	if qty <= 0 {
		return LineItem{}, fmt.Errorf("%w: %d", ErrInvalidQuantity, qty)
	}
	if err := unitPrice.Validate(); err != nil {
		return LineItem{}, fmt.Errorf("order: unit_price: %w", err)
	}
	total, err := unitPrice.Multiply(qty)
	if err != nil {
		return LineItem{}, err
	}
	return LineItem{
		ProductID:  productID,
		SKU:        sku,
		Name:       name,
		Quantity:   qty,
		UnitPrice:  unitPrice,
		TotalPrice: total,
	}, nil
}

type Order struct {
	ID         string
	TenantID   shared.TenantID
	CustomerID shared.CustomerID
	SaleMode   SaleMode
	ProviderID *shared.ProviderID
	Status     Status
	Currency   shared.Currency
	Items      []LineItem

	Total           shared.Money
	UvicoCommission shared.Money
	ProviderPayout  shared.Money
	GatewayFee      shared.Money
	MerchantNet     shared.Money

	PaymentRef  string
	Version     int64
	CreatedAt   time.Time
	UpdatedAt   time.Time
	PaidAt      *time.Time
	FailedAt    *time.Time
	CancelledAt *time.Time
	RefundedAt  *time.Time

	// Pago manual: campos introducidos en 0023.
	PaymentMethod             string
	AccessToken               string
	CustomerPaymentReference  string
	CustomerPaymentProofURL   string
	CustomerSubmittedAt       *time.Time
	MerchantConfirmedAt       *time.Time
	MerchantConfirmedBy       string
	MerchantRejectionReason   string
}

type SplitConfig struct {
	UvicoCommissionBps int64
	ProviderCostBps    int64
	GatewayFeeBps      int64
}

func NewOrder(
	tenantID shared.TenantID,
	customerID shared.CustomerID,
	saleMode SaleMode,
	providerID *shared.ProviderID,
	items []LineItem,
	split SplitConfig,
	paymentMethod string,
	accessToken string,
	now time.Time,
) (*Order, error) {
	if tenantID.IsZero() {
		return nil, ErrInvalidTenant
	}
	if customerID.IsZero() {
		return nil, ErrInvalidCustomer
	}
	if len(items) == 0 {
		return nil, ErrInvalidItems
	}
	if paymentMethod == "" {
		return nil, ErrMissingPaymentMethod
	}
	if accessToken == "" {
		return nil, ErrMissingAccessToken
	}
	switch saleMode {
	case SaleModeOwnInventory:
		if providerID != nil {
			return nil, fmt.Errorf("%w: own_inventory no lleva provider", ErrProviderRequired)
		}
	case SaleModeDropshipping:
		if providerID == nil {
			return nil, ErrProviderRequired
		}
	default:
		return nil, fmt.Errorf("order: sale_mode desconocido %q", saleMode)
	}

	currency := items[0].UnitPrice.Currency
	total := shared.Money{Currency: currency}
	for i, it := range items {
		if it.UnitPrice.Currency != currency || it.TotalPrice.Currency != currency {
			return nil, fmt.Errorf("%w: item %d", ErrMixedCurrencies, i)
		}
		var err error
		total, err = total.Add(it.TotalPrice)
		if err != nil {
			return nil, err
		}
	}
	if total.IsZero() || total.IsNegative() {
		return nil, ErrInvalidTotal
	}

	nowUTC := now.UTC()
	o := &Order{
		TenantID:      tenantID,
		CustomerID:    customerID,
		SaleMode:      saleMode,
		ProviderID:    providerID,
		Status:        StatusPendingPayment,
		Currency:      currency,
		Items:         items,
		Total:         total,
		PaymentMethod: paymentMethod,
		AccessToken:   accessToken,
		Version:       1,
		CreatedAt:     nowUTC,
		UpdatedAt:     nowUTC,
	}

	if saleMode == SaleModeOwnInventory {
		o.UvicoCommission = shared.Money{Currency: currency}
		o.ProviderPayout = shared.Money{Currency: currency}
		o.GatewayFee = shared.Money{Currency: currency}
		o.MerchantNet = total
		return o, nil
	}

	parts, merchant, err := shared.SplitBps(total, split.UvicoCommissionBps, split.ProviderCostBps, split.GatewayFeeBps)
	if err != nil {
		return nil, err
	}
	o.UvicoCommission = parts[0]
	o.ProviderPayout = parts[1]
	o.GatewayFee = parts[2]
	o.MerchantNet = merchant

	sum := o.UvicoCommission.AmountCents + o.ProviderPayout.AmountCents + o.GatewayFee.AmountCents + o.MerchantNet.AmountCents
	if sum != total.AmountCents {
		return nil, fmt.Errorf("%w: sum=%d total=%d", ErrSplitInvariant, sum, total.AmountCents)
	}
	return o, nil
}

// SubmitPayment: el comprador envía la referencia bancaria y el comprobante.
// Transición: pending_payment → awaiting_payment.
// Idempotente: si ya está awaiting_payment con la misma referencia, no-op.
func (o *Order) SubmitPayment(reference, proofURL string, now time.Time) error {
	if reference == "" {
		return ErrMissingReference
	}
	if o.Status == StatusAwaitingPayment {
		if o.CustomerPaymentReference == reference {
			return nil
		}
		return fmt.Errorf("%w: awaiting_payment con otra referencia", ErrInvalidTransition)
	}
	if o.Status != StatusPendingPayment {
		return fmt.Errorf("%w: %s -> awaiting_payment", ErrInvalidTransition, o.Status)
	}
	t := now.UTC()
	o.Status = StatusAwaitingPayment
	o.CustomerPaymentReference = reference
	o.CustomerPaymentProofURL = proofURL
	o.CustomerSubmittedAt = &t
	o.UpdatedAt = t
	o.Version++
	return nil
}

// ConfirmPayment: el merchant confirma que recibió el pago.
// Transición: awaiting_payment → paid.
// Idempotente: si ya está paid, no-op.
func (o *Order) ConfirmPayment(merchantUserID string, now time.Time) error {
	if o.Status == StatusPaid {
		return nil
	}
	if o.Status != StatusAwaitingPayment {
		return fmt.Errorf("%w: %s -> paid", ErrInvalidTransition, o.Status)
	}
	if merchantUserID == "" {
		return ErrMissingMerchantID
	}
	t := now.UTC()
	o.Status = StatusPaid
	o.PaidAt = &t
	o.MerchantConfirmedAt = &t
	o.MerchantConfirmedBy = merchantUserID
	o.UpdatedAt = t
	o.Version++
	return nil
}

// RejectPayment: el merchant rechaza el pago (comprobante inválido, monto
// incorrecto, transferencia no recibida, etc.).
// Transición: awaiting_payment → failed.
func (o *Order) RejectPayment(merchantUserID, reason string, now time.Time) error {
	if o.Status == StatusFailed {
		return nil
	}
	if o.Status != StatusAwaitingPayment {
		return fmt.Errorf("%w: %s -> failed", ErrInvalidTransition, o.Status)
	}
	if reason == "" {
		return ErrMissingRejectionReason
	}
	t := now.UTC()
	o.Status = StatusFailed
	o.FailedAt = &t
	o.MerchantConfirmedBy = merchantUserID
	o.MerchantRejectionReason = reason
	o.UpdatedAt = t
	o.Version++
	return nil
}

// MarkPaid: transición directa pending_payment → paid. Solo para webhooks de
// pasarelas (Wayu) que confirman el pago sin pasar por awaiting_payment.
func (o *Order) MarkPaid(paymentRef string, now time.Time) error {
	if o.Status == StatusPaid {
		return nil
	}
	if o.Status != StatusPendingPayment {
		return fmt.Errorf("%w: %s -> paid", ErrInvalidTransition, o.Status)
	}
	if paymentRef == "" {
		return ErrMissingReference
	}
	t := now.UTC()
	o.Status = StatusPaid
	o.PaymentRef = paymentRef
	o.PaidAt = &t
	o.UpdatedAt = t
	o.Version++
	return nil
}

// MarkFailed: acepta pending_payment y awaiting_payment → failed.
func (o *Order) MarkFailed(now time.Time) error {
	if o.Status == StatusFailed {
		return nil
	}
	if o.Status != StatusPendingPayment && o.Status != StatusAwaitingPayment {
		return fmt.Errorf("%w: %s -> failed", ErrInvalidTransition, o.Status)
	}
	t := now.UTC()
	o.Status = StatusFailed
	o.FailedAt = &t
	o.UpdatedAt = t
	o.Version++
	return nil
}

// Cancel: pending_payment, awaiting_payment, o failed → cancelled.
func (o *Order) Cancel(now time.Time) error {
	if o.Status == StatusCancelled {
		return nil
	}
	if o.Status != StatusPendingPayment &&
		o.Status != StatusAwaitingPayment &&
		o.Status != StatusFailed {
		return fmt.Errorf("%w: %s -> cancelled", ErrInvalidTransition, o.Status)
	}
	t := now.UTC()
	o.Status = StatusCancelled
	o.CancelledAt = &t
	o.UpdatedAt = t
	o.Version++
	return nil
}

func (o *Order) Refund(now time.Time) error {
	if o.Status == StatusRefunded {
		return nil
	}
	if o.Status != StatusPaid {
		return fmt.Errorf("%w: %s -> refunded", ErrInvalidTransition, o.Status)
	}
	t := now.UTC()
	o.Status = StatusRefunded
	o.RefundedAt = &t
	o.UpdatedAt = t
	o.Version++
	return nil
}
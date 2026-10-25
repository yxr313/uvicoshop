package order

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"uvico-api/internal/domain/attribution"
	"uvico-api/internal/domain/customer"
	"uvico-api/internal/domain/inventory"
	orderdomain "uvico-api/internal/domain/order"
	"uvico-api/internal/domain/paymentmethod"
	"uvico-api/internal/domain/shared"
	"uvico-api/internal/platform/database"
)

type Service struct {
	pool            *database.Pool
	orderRepo       orderdomain.Repository
	invRepo         inventory.Repository
	customerRepo    customer.Repository
	attribution     attribution.Repository
	pmRepo          paymentmethod.Repository
	idemStore       orderdomain.IdempotencyStore
	logger          *slog.Logger
	split           orderdomain.SplitConfig
	custodyEnabled  bool
	now             func() time.Time
}

type ServiceOptions struct {
	Split          orderdomain.SplitConfig
	CustodyEnabled bool
}

func NewService(
	pool *database.Pool,
	or orderdomain.Repository,
	ir inventory.Repository,
	cr customer.Repository,
	ar attribution.Repository,
	pmr paymentmethod.Repository,
	idem orderdomain.IdempotencyStore,
	logger *slog.Logger,
	opts ServiceOptions,
) *Service {
	return &Service{
		pool:           pool,
		orderRepo:      or,
		invRepo:        ir,
		customerRepo:   cr,
		attribution:    ar,
		pmRepo:         pmr,
		idemStore:      idem,
		logger:         logger,
		split:          opts.Split,
		custodyEnabled: opts.CustodyEnabled,
		now:            time.Now,
	}
}

type CheckoutInput struct {
	TenantID      shared.TenantID
	IdemKey       string
	Customer      customer.UpsertInput
	Attribution   attribution.Attribution
	SaleMode      orderdomain.SaleMode
	ProviderID    *shared.ProviderID
	PaymentMethod string
	Items         []orderdomain.LineItem
}

type CheckoutResult struct {
	Order            *orderdomain.Order
	Idempotent       bool
	PaymentMethod    *paymentmethod.PaymentMethod
}

func (s *Service) Checkout(ctx context.Context, in CheckoutInput) (*CheckoutResult, error) {
	// 0. Bloqueo de dropshipping en MVP
	if in.SaleMode == orderdomain.SaleModeDropshipping && !s.custodyEnabled {
		return nil, orderdomain.ErrDropshippingDisabled
	}

	// 1. Idempotency lookup
	if in.IdemKey != "" {
		reqHash := hashCheckout(in)
		existing, err := s.idemStore.Get(ctx, in.TenantID.String(), in.IdemKey)
		if err != nil {
			s.logger.Warn("checkout: idempotency get failed, continuing", "err", err)
		} else if existing != nil {
			if existing.RequestHash != reqHash {
				return nil, orderdomain.ErrIdempotencyConflict
			}
			o, err := s.orderRepo.GetByID(ctx, in.TenantID, existing.OrderID)
			if err != nil {
				return nil, fmt.Errorf("checkout: load idempotent order: %w", err)
			}
			pm, _ := s.pmRepo.GetByMethod(ctx, in.TenantID, paymentmethod.Method(o.PaymentMethod))
			s.logger.Info("checkout: idempotent replay", "order_id", existing.OrderID, "key", in.IdemKey)
			return &CheckoutResult{Order: o, Idempotent: true, PaymentMethod: pm}, nil
		}
	}

	// 2. Validar método de pago contra la config del merchant
	pm, err := s.pmRepo.GetByMethod(ctx, in.TenantID, paymentmethod.Method(in.PaymentMethod))
	if err != nil {
		return nil, fmt.Errorf("checkout: payment method: %w", err)
	}
	if !pm.Enabled {
		return nil, fmt.Errorf("%w: %s", paymentmethod.ErrDisabled, in.PaymentMethod)
	}

	// 3. Customer upsert con retry
	custID, err := s.upsertCustomerWithRetry(ctx, in.Customer)
	if err != nil {
		return nil, fmt.Errorf("checkout: customer upsert: %w", err)
	}

	// 4. Build order
	accessToken := uuid.Must(uuid.NewV7()).String()
	o, err := orderdomain.NewOrder(
		in.TenantID,
		shared.CustomerID(custID),
		in.SaleMode,
		in.ProviderID,
		in.Items,
		s.split,
		in.PaymentMethod,
		accessToken,
		s.now(),
	)
	if err != nil {
		return nil, err
	}
	o.ID = uuid.Must(uuid.NewV7()).String()

	// 5. Reserve stock (own_inventory)
	if in.SaleMode == orderdomain.SaleModeOwnInventory {
		for _, item := range in.Items {
			if err := s.invRepo.Reserve(ctx, in.TenantID, item.ProductID, int64(item.Quantity), o.ID); err != nil {
				return nil, err
			}
		}
	}

	// 6. Insert order + attribution + outbox
	attributionRecord := attribution.Record{
		OrderID:        o.ID,
		TenantID:       in.TenantID.String(),
		SessionID:      in.Attribution.SessionID,
		UTMSource:      in.Attribution.UTMSource,
		UTMMedium:      in.Attribution.UTMMedium,
		UTMCampaign:    in.Attribution.UTMCampaign,
		UTMTerm:        in.Attribution.UTMTerm,
		UTMContent:     in.Attribution.UTMContent,
		FBCLID:         in.Attribution.FBCLID,
		GCLID:          in.Attribution.GCLID,
		TTCLID:         in.Attribution.TTCLID,
		TWCLID:         in.Attribution.TWCLID,
		Referrer:       in.Attribution.Referrer,
		LandingPage:    in.Attribution.LandingPage,
		UserAgent:      in.Attribution.UserAgent,
		AcceptLanguage: in.Attribution.AcceptLanguage,
		Timezone:       in.Attribution.Timezone,
		IPHash:         in.Attribution.IPHash,
		CreatedAt:      o.CreatedAt,
	}

	err = s.pool.WithTenantTx(ctx, in.TenantID.String(), func(ctx context.Context, tx pgx.Tx) error {
		if err := s.orderRepo.CreateWithReservationTx(ctx, tx, o); err != nil {
			return fmt.Errorf("insert order: %w", err)
		}

		if !in.Attribution.IsEmpty() {
			if err := s.attribution.InsertOrderAttributionTx(ctx, tx, attributionRecord); err != nil {
				s.logger.Warn("checkout: order_attribution insert failed", "order_id", o.ID, "err", err)
			}
			if err := s.attribution.InsertFirstTouchTx(ctx, tx, custID, attributionRecord); err != nil {
				s.logger.Warn("checkout: customer_attribution insert failed", "customer_id", custID, "err", err)
			}
		}

		return insertOutboxTx(ctx, tx, in.TenantID.String(), "order", o.ID, "order.created", map[string]any{
			"order_id":       o.ID,
			"customer_email": in.Customer.Email,
			"total_cents":    o.Total.AmountCents,
			"currency":       string(o.Currency),
			"payment_method": o.PaymentMethod,
		})
	})
	if err != nil {
		return nil, err
	}

	// 7. Idempotency put (best-effort)
	if in.IdemKey != "" {
		if err := s.idemStore.Put(ctx, orderdomain.IdempotencyRecord{
			Key:         in.IdemKey,
			TenantID:    in.TenantID.String(),
			OrderID:     o.ID,
			RequestHash: hashCheckout(in),
		}); err != nil {
			s.logger.Warn("checkout: idempotency put failed, order created",
				"order_id", o.ID, "key", in.IdemKey, "err", err)
			if errors.Is(err, orderdomain.ErrIdempotencyConflict) {
				return nil, err
			}
		}
	}

	s.logger.Info("checkout: order created",
		"order_id", o.ID,
		"tenant_id", in.TenantID.String(),
		"customer_id", custID,
		"payment_method", o.PaymentMethod,
		"utm_source", in.Attribution.UTMSource,
	)
	return &CheckoutResult{Order: o, PaymentMethod: pm}, nil
}

// SubmitPayment: el comprador envía la referencia bancaria y el comprobante.
func (s *Service) SubmitPayment(ctx context.Context, accessToken, reference, proofURL string) (*orderdomain.Order, error) {
	if accessToken == "" {
		return nil, orderdomain.ErrMissingAccessToken
	}
	if reference == "" {
		return nil, orderdomain.ErrMissingReference
	}

	o, err := s.orderRepo.GetByAccessToken(ctx, accessToken)
	if err != nil {
		return nil, err
	}

	if err := o.SubmitPayment(reference, proofURL, s.now()); err != nil {
		return nil, err
	}

	if err := s.orderRepo.UpdateWithOutbox(ctx, o, "payment.submitted", map[string]any{
		"order_id":       o.ID,
		"tenant_id":      o.TenantID.String(),
		"customer_email": "",
		"total_cents":    o.Total.AmountCents,
		"currency":       string(o.Currency),
		"reference":      reference,
	}); err != nil {
		return nil, err
	}

	s.logger.Info("payment submitted",
		"order_id", o.ID,
		"reference", reference,
		"proof_url_set", proofURL != "",
	)
	return o, nil
}

// ConfirmPayment: el merchant confirma el pago.
func (s *Service) ConfirmPayment(ctx context.Context, tenantID shared.TenantID, orderID, merchantUserID string) (*orderdomain.Order, error) {
	o, err := s.orderRepo.GetByID(ctx, tenantID, orderID)
	if err != nil {
		return nil, err
	}
	if err := o.ConfirmPayment(merchantUserID, s.now()); err != nil {
		return nil, err
	}
	if err := s.orderRepo.UpdateWithOutbox(ctx, o, "payment.confirmed", map[string]any{
		"order_id": o.ID,
		"tenant_id": o.TenantID.String(),
	}); err != nil {
		return nil, err
	}
	s.logger.Info("payment confirmed", "order_id", o.ID, "merchant_user_id", merchantUserID)
	return o, nil
}

// RejectPayment: el merchant rechaza el pago.
func (s *Service) RejectPayment(ctx context.Context, tenantID shared.TenantID, orderID, merchantUserID, reason string) (*orderdomain.Order, error) {
	o, err := s.orderRepo.GetByID(ctx, tenantID, orderID)
	if err != nil {
		return nil, err
	}
	if err := o.RejectPayment(merchantUserID, reason, s.now()); err != nil {
		return nil, err
	}
	if err := s.orderRepo.UpdateWithOutbox(ctx, o, "payment.rejected", map[string]any{
		"order_id": o.ID,
		"tenant_id": o.TenantID.String(),
		"reason":    reason,
	}); err != nil {
		return nil, err
	}
	s.logger.Info("payment rejected", "order_id", o.ID, "reason", reason)
	return o, nil
}

func (s *Service) upsertCustomerWithRetry(ctx context.Context, in customer.UpsertInput) (string, error) {
	var lastErr error
	for attempt := 0; attempt < 3; attempt++ {
		custID, err := s.customerRepo.UpsertByEmail(ctx, in)
		if err == nil {
			return custID, nil
		}
		lastErr = err
		if !isTransientPgErr(err) {
			return "", err
		}
		delay := time.Duration(50*(1<<attempt)) * time.Millisecond
		select {
		case <-ctx.Done():
			return "", ctx.Err()
		case <-time.After(delay):
		}
		s.logger.Warn("checkout: customer upsert retry",
			"attempt", attempt+1, "delay_ms", delay.Milliseconds(), "err", err)
	}
	return "", lastErr
}

func (s *Service) ListOrders(ctx context.Context, filter orderdomain.ListFilter) ([]*orderdomain.Order, error) {
	if filter.Limit <= 0 || filter.Limit > 100 {
		filter.Limit = 50
	}
	return s.orderRepo.List(ctx, filter.TenantID, filter)
}

func (s *Service) GetOrder(ctx context.Context, tenantID shared.TenantID, orderID string) (*orderdomain.Order, error) {
	return s.orderRepo.GetByID(ctx, tenantID, orderID)
}

func (s *Service) GetByAccessToken(ctx context.Context, accessToken string) (*orderdomain.Order, error) {
	return s.orderRepo.GetByAccessToken(ctx, accessToken)
}

func hashCheckout(in CheckoutInput) string {
	h := sha256.New()
	_, _ = h.Write([]byte(in.Customer.Email))
	_, _ = h.Write([]byte("|" + in.Customer.Name))
	_, _ = h.Write([]byte("|" + in.Customer.Phone))
	_, _ = h.Write([]byte("|" + string(in.SaleMode)))
	_, _ = h.Write([]byte("|" + in.PaymentMethod))
	if in.ProviderID != nil {
		_, _ = h.Write([]byte("|" + in.ProviderID.String()))
	}
	for _, it := range in.Items {
		_, _ = fmt.Fprintf(h, "|%s:%d:%d", it.ProductID.String(), it.Quantity, it.UnitPrice.AmountCents)
	}
	return hex.EncodeToString(h.Sum(nil))
}

func isTransientPgErr(err error) bool {
	var pgErr interface{ SQLState() string }
	if errors.As(err, &pgErr) {
		switch pgErr.SQLState() {
		case "40001", "40P01", "57014":
			return true
		}
	}
	return false
}

func insertOutboxTx(ctx context.Context, tx pgx.Tx, tenantID, aggregateType, aggregateID, eventType string, payload any) error {
	raw, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("outbox: marshal: %w", err)
	}
	_, err = tx.Exec(ctx, `
		INSERT INTO outbox_events
			(id, tenant_id, aggregate_type, aggregate_id, event_type, payload, created_at)
		VALUES (gen_random_uuid(), $1, $2, $3, $4, $5, NOW())
	`, tenantID, aggregateType, aggregateID, eventType, string(raw))
	if err != nil {
		return fmt.Errorf("outbox: insert: %w", err)
	}
	return nil
}
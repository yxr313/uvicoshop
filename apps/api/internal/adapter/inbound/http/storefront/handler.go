package storefront

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"log/slog"
	"net"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"

	appcart "uvico-api/internal/application/cart"
	apporder "uvico-api/internal/application/order"
	apppm "uvico-api/internal/application/paymentmethod"
	appstorefront "uvico-api/internal/application/storefront"
	"uvico-api/internal/domain/attribution"
	"uvico-api/internal/domain/customer"
	orderdomain "uvico-api/internal/domain/order"
	"uvico-api/internal/domain/paymentmethod"
	"uvico-api/internal/domain/shared"
)

const maxPublicBodyBytes = 1 << 20

type Handler struct {
	svc      *appstorefront.Service
	orderSvc *apporder.Service
	cartSvc  *appcart.Service
	pmSvc    *apppm.Service
	logger   *slog.Logger
	validate *validator.Validate
	hmacKey  []byte
}

func New(
	svc *appstorefront.Service,
	orderSvc *apporder.Service,
	cartSvc *appcart.Service,
	pmSvc *apppm.Service,
	logger *slog.Logger,
	hmacKey []byte,
) *Handler {
	v := validator.New(validator.WithRequiredStructEnabled())
	v.RegisterValidation("uuid_any", func(fl validator.FieldLevel) bool {
		s, ok := fl.Field().Interface().(string)
		if !ok || s == "" {
			return false
		}
		_, err := uuid.Parse(s)
		return err == nil
	})
	return &Handler{
		svc:      svc,
		orderSvc: orderSvc,
		cartSvc:  cartSvc,
		pmSvc:    pmSvc,
		logger:   logger,
		validate: v,
		hmacKey:  hmacKey,
	}
}

func (h *Handler) Routes(resolver *SlugResolver) chi.Router {
	r := chi.NewRouter()
	r.Route("/{slug}", func(r chi.Router) {
		r.Use(resolver.Middleware())
		r.Get("/config", h.getConfig)
		r.Get("/products", h.listProducts)
		r.Get("/products/{productSlug}", h.getProduct)
		r.Get("/payment-methods", h.listPaymentMethods)
		r.Post("/checkout", h.checkout)
		r.Post("/orders/{orderID}/submit-payment", h.submitPayment)
		r.Get("/orders/{orderID}", h.getOrderPublic)
	})
	return r
}

// ─── Existentes ────────────────────────────────────────────────────

func (h *Handler) getConfig(w http.ResponseWriter, r *http.Request) {
	tenant, err := MustTenant(r.Context())
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "missing_tenant", "tenant no resuelto")
		return
	}
	cfg, err := h.svc.GetConfig(r.Context(), tenant.ID)
	if err != nil {
		h.logger.Error("storefront: get config", "err", err, "tenant_id", tenant.ID.String())
		writeJSONError(w, http.StatusInternalServerError, "internal", "error interno")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "public, max-age=60")
	_ = json.NewEncoder(w).Encode(cfg)
}

func (h *Handler) listProducts(w http.ResponseWriter, r *http.Request) {
	tenant, err := MustTenant(r.Context())
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "missing_tenant", "tenant no resuelto")
		return
	}
	products, err := h.svc.ListProducts(r.Context(), tenant.ID, 50)
	if err != nil {
		h.logger.Error("storefront: list products", "err", err, "tenant_id", tenant.ID.String())
		writeJSONError(w, http.StatusInternalServerError, "internal", "error interno")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "public, max-age=30")
	_ = json.NewEncoder(w).Encode(map[string]any{"products": products})
}

func (h *Handler) getProduct(w http.ResponseWriter, r *http.Request) {
	tenant, err := MustTenant(r.Context())
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "missing_tenant", "tenant no resuelto")
		return
	}
	productSlug := chi.URLParam(r, "productSlug")
	p, err := h.svc.GetProduct(r.Context(), tenant.ID, productSlug)
	if err != nil {
		if errors.Is(err, appstorefront.ErrProductNotFound) {
			writeJSONError(w, http.StatusNotFound, "product_not_found", "producto no encontrado")
			return
		}
		h.logger.Error("storefront: get product", "err", err, "tenant_id", tenant.ID.String(), "slug", productSlug)
		writeJSONError(w, http.StatusInternalServerError, "internal", "error interno")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "public, max-age=60")
	_ = json.NewEncoder(w).Encode(p)
}

// ─── Nuevos ────────────────────────────────────────────────────────

func (h *Handler) listPaymentMethods(w http.ResponseWriter, r *http.Request) {
	tenant, err := MustTenant(r.Context())
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "missing_tenant", "tenant no resuelto")
		return
	}
	methods, err := h.pmSvc.ListPublic(r.Context(), tenant.ID)
	if err != nil {
		h.logger.Error("storefront: list payment methods", "err", err, "tenant_id", tenant.ID.String())
		writeJSONError(w, http.StatusInternalServerError, "internal", "error interno")
		return
	}
	out := make([]paymentMethodPublicDTO, 0, len(methods))
	for _, m := range methods {
		out = append(out, paymentMethodPublicDTO{
			Method:   string(m.Method),
			Label:    paymentmethod.MethodLabel(m.Method),
			Position: m.Position,
			Config:   m.PublicView(),
		})
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "public, max-age=60")
	_ = json.NewEncoder(w).Encode(map[string]any{"methods": out})
}

func (h *Handler) checkout(w http.ResponseWriter, r *http.Request) {
	tenant, err := MustTenant(r.Context())
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "missing_tenant", "tenant no resuelto")
		return
	}

	idemKey := r.Header.Get("Idempotency-Key")
	if idemKey == "" {
		idemKey = "public-" + uuid.Must(uuid.NewV7()).String()
	}

	r.Body = http.MaxBytesReader(w, r.Body, maxPublicBodyBytes)
	var req checkoutRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid_json", "cuerpo inválido")
		return
	}
	if err := h.validate.Struct(req); err != nil {
		writeJSONError(w, http.StatusUnprocessableEntity, "validation_error", err.Error())
		return
	}

	// Resolver items: cart_id o items directos.
	var items []orderdomain.LineItem
	if req.CartID != "" {
		lines, err := h.cartSvc.ValidateForCheckout(r.Context(), tenant.ID, req.CartID)
		if err != nil {
			var pcErr *appcart.PriceChangedError
			if errors.As(err, &pcErr) {
				writeJSONError(w, http.StatusConflict, "price_changed",
					"los precios cambiaron, refresca tu carrito")
				return
			}
			writeJSONError(w, http.StatusUnprocessableEntity, "cart_error", err.Error())
			return
		}
		for _, l := range lines {
			item, err := orderdomain.NewLineItem(l.ProductID, l.SKU, l.Name, l.Quantity, l.UnitPrice)
			if err != nil {
				writeJSONError(w, http.StatusInternalServerError, "internal", "error interno")
				return
			}
			items = append(items, item)
		}
		defer func() {
			_ = h.cartSvc.Clear(context.WithoutCancel(r.Context()), tenant.ID, req.CartID)
		}()
	} else if len(req.Items) > 0 {
		for _, dto := range req.Items {
			item, err := orderdomain.NewLineItem(
				shared.ProductID(dto.ProductID),
				dto.SKU, dto.Name, dto.Quantity,
				shared.USD(dto.PriceCents),
			)
			if err != nil {
				writeJSONError(w, http.StatusUnprocessableEntity, "invalid_item", err.Error())
				return
			}
			items = append(items, item)
		}
	} else {
		writeJSONError(w, http.StatusBadRequest, "missing_items", "cart_id o items requerido")
		return
	}

	var providerID *shared.ProviderID
	if req.ProviderID != nil {
		pid := shared.ProviderID(*req.ProviderID)
		providerID = &pid
	}

	attr := attribution.Attribution{
		SessionID:      req.Attribution.SessionID,
		UTMSource:      req.Attribution.UTMSource,
		UTMMedium:      req.Attribution.UTMMedium,
		UTMCampaign:    req.Attribution.UTMCampaign,
		UTMTerm:        req.Attribution.UTMTerm,
		UTMContent:     req.Attribution.UTMContent,
		FBCLID:         req.Attribution.FBCLID,
		GCLID:          req.Attribution.GCLID,
		TTCLID:         req.Attribution.TTCLID,
		TWCLID:         req.Attribution.TWCLID,
		Referrer:       req.Attribution.Referrer,
		LandingPage:    req.Attribution.LandingPage,
		UserAgent:      r.UserAgent(),
		AcceptLanguage: r.Header.Get("Accept-Language"),
		Timezone:       req.Attribution.Timezone,
		IPHash:         h.hashIP(clientIP(r)),
	}
	if err := attr.Validate(); err != nil {
		writeJSONError(w, http.StatusUnprocessableEntity, "invalid_attribution", err.Error())
		return
	}

	result, err := h.orderSvc.Checkout(r.Context(), apporder.CheckoutInput{
		TenantID: tenant.ID,
		IdemKey:  idemKey,
		Customer: customer.UpsertInput{
			Email:          req.CustomerEmail,
			Name:           req.CustomerName,
			Phone:          req.CustomerPhone,
			MarketingOptIn: req.MarketingOptIn,
		},
		Attribution:   attr,
		SaleMode:      orderdomain.SaleMode(req.SaleMode),
		ProviderID:    providerID,
		PaymentMethod: req.PaymentMethod,
		Items:         items,
	})
	if err != nil {
		h.handleOrderError(w, err)
		return
	}

	instructions := paymentmethod.BuildInstructions(result.PaymentMethod, result.Order.ID)

	resp := checkoutResponse{
		OrderID:     result.Order.ID,
		AccessToken: result.Order.AccessToken,
		Status:      string(result.Order.Status),
		TotalCents:  result.Order.Total.AmountCents,
		Currency:    string(result.Order.Currency),
		PaymentInstructions: paymentInstructionsDTO{
			Method:    instructions.Method,
			Label:     instructions.Label,
			Config:    instructions.Config,
			Reference: instructions.Reference,
		},
	}

	w.Header().Set("Idempotency-Key", idemKey)
	if result.Idempotent {
		w.Header().Set("Idempotent-Replay", "true")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(resp)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(resp)
}

func (h *Handler) submitPayment(w http.ResponseWriter, r *http.Request) {
	tenant, err := MustTenant(r.Context())
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "missing_tenant", "tenant no resuelto")
		return
	}
	orderID := chi.URLParam(r, "orderID")
	if _, err := uuid.Parse(orderID); err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid_order_id", "order_id inválido")
		return
	}
	orderToken := r.Header.Get("X-Order-Token")
	if orderToken == "" {
		writeJSONError(w, http.StatusUnauthorized, "missing_order_token", "X-Order-Token requerido")
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, maxPublicBodyBytes)
	var req submitPaymentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid_json", "cuerpo inválido")
		return
	}
	if err := h.validate.Struct(req); err != nil {
		writeJSONError(w, http.StatusUnprocessableEntity, "validation_error", err.Error())
		return
	}

	// Verificar token antes de procesar
	o, err := h.orderSvc.GetByAccessToken(r.Context(), orderToken)
	if err != nil || o.ID != orderID || o.TenantID != tenant.ID {
		writeJSONError(w, http.StatusUnauthorized, "invalid_order_token", "token inválido o no coincide con la orden")
		return
	}

	updated, err := h.orderSvc.SubmitPayment(r.Context(), orderToken, req.Reference, req.ProofURL)
	if err != nil {
		h.handleOrderError(w, err)
		return
	}

	submittedAt := ""
	if updated.CustomerSubmittedAt != nil {
		submittedAt = updated.CustomerSubmittedAt.Format(time.RFC3339)
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(submitPaymentResponse{
		OrderID:             updated.ID,
		Status:              string(updated.Status),
		Reference:           updated.CustomerPaymentReference,
		CustomerSubmittedAt: submittedAt,
	})
}

func (h *Handler) getOrderPublic(w http.ResponseWriter, r *http.Request) {
	tenant, err := MustTenant(r.Context())
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "missing_tenant", "tenant no resuelto")
		return
	}
	orderID := chi.URLParam(r, "orderID")
	if _, err := uuid.Parse(orderID); err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid_order_id", "order_id inválido")
		return
	}
	orderToken := r.Header.Get("X-Order-Token")
	if orderToken == "" {
		writeJSONError(w, http.StatusUnauthorized, "missing_order_token", "X-Order-Token requerido")
		return
	}

	o, err := h.orderSvc.GetByAccessToken(r.Context(), orderToken)
	if err != nil || o.ID != orderID || o.TenantID != tenant.ID {
		writeJSONError(w, http.StatusUnauthorized, "invalid_order_token", "token inválido o no coincide con la orden")
		return
	}

	// Obtener datos del método de pago para las instrucciones
	pm, _ := h.pmSvc.GetForCheckout(r.Context(), tenant.ID, o.PaymentMethod)
	var instructions *paymentInstructionsDTO
	if pm != nil {
		ins := paymentmethod.BuildInstructions(pm, o.ID)
		instructions = &paymentInstructionsDTO{
			Method:    ins.Method,
			Label:     ins.Label,
			Config:    ins.Config,
			Reference: ins.Reference,
		}
	}

	items := make([]orderItemPublicDTO, 0, len(o.Items))
	for _, it := range o.Items {
		items = append(items, orderItemPublicDTO{
			ProductID: it.ProductID.String(),
			SKU:       it.SKU,
			Name:      it.Name,
			Quantity:  it.Quantity,
			UnitPrice: it.UnitPrice.AmountCents,
			Total:     it.TotalPrice.AmountCents,
		})
	}

	resp := orderPublicResponse{
		OrderID:                 o.ID,
		Status:                  string(o.Status),
		TotalCents:              o.Total.AmountCents,
		Currency:                string(o.Currency),
		Items:                   items,
		MerchantRejectionReason: o.MerchantRejectionReason,
	}
	if instructions != nil {
		resp.PaymentInstructions = *instructions
	}
	if o.CustomerSubmittedAt != nil {
		resp.CustomerSubmittedAt = o.CustomerSubmittedAt.Format(time.RFC3339)
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}

// ─── Helpers ───────────────────────────────────────────────────────

func (h *Handler) handleOrderError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, orderdomain.ErrDropshippingDisabled):
		writeJSONError(w, http.StatusUnprocessableEntity, "dropshipping_disabled",
			"el modo dropshipping no está disponible aún; usa sale_mode=own_inventory")
	case errors.Is(err, orderdomain.ErrInvalidTransition):
		writeJSONError(w, http.StatusConflict, "invalid_transition", err.Error())
	case errors.Is(err, orderdomain.ErrMissingReference):
		writeJSONError(w, http.StatusUnprocessableEntity, "missing_reference", err.Error())
	case errors.Is(err, paymentmethod.ErrDisabled),
		errors.Is(err, paymentmethod.ErrNotFound):
		writeJSONError(w, http.StatusUnprocessableEntity, "payment_method_unavailable", err.Error())
	default:
		h.logger.Error("storefront: checkout error", "err", err)
		writeJSONError(w, http.StatusInternalServerError, "internal", "error interno")
	}
}

func clientIP(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		for i := 0; i < len(xff); i++ {
			if xff[i] == ',' {
				return trimSpace(xff[:i])
			}
		}
		return trimSpace(xff)
	}
	if rip := r.Header.Get("X-Real-IP"); rip != "" {
		return trimSpace(rip)
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}

func trimSpace(s string) string {
	start, end := 0, len(s)
	for start < end && (s[start] == ' ' || s[start] == '\t') {
		start++
	}
	for end > start && (s[end-1] == ' ' || s[end-1] == '\t') {
		end--
	}
	return s[start:end]
}

func (h *Handler) hashIP(ip string) string {
	if len(h.hmacKey) == 0 || ip == "" {
		return ""
	}
	mac := hmac.New(sha256.New, h.hmacKey)
	mac.Write([]byte(ip))
	return hex.EncodeToString(mac.Sum(nil))
}
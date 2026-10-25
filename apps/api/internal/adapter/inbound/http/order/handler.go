package order

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"

	"uvico-api/internal/adapter/inbound/http/middleware"
	apporder "uvico-api/internal/application/order"
	orderdomain "uvico-api/internal/domain/order"
	"uvico-api/internal/domain/paymentmethod"
	"uvico-api/internal/domain/shared"
)

const maxBodyBytes = 32 * 1024

// Handler expone SOLO endpoints del merchant (autenticados con JWT).
// El checkout del comprador vive en /v1/storefront/{slug}/checkout.
type Handler struct {
	svc      *apporder.Service
	logger   *slog.Logger
	validate *validator.Validate
}

func New(svc *apporder.Service, logger *slog.Logger) *Handler {
	return &Handler{
		svc:      svc,
		logger:   logger,
		validate: validator.New(validator.WithRequiredStructEnabled()),
	}
}

func (h *Handler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Get("/", h.list)
	r.Get("/{orderID}", h.get)
	r.Post("/{orderID}/confirm-payment", h.confirmPayment)
	r.Post("/{orderID}/reject-payment", h.rejectPayment)
	r.Post("/{orderID}/cancel", h.cancel)
	return r
}

func (h *Handler) list(w http.ResponseWriter, r *http.Request) {
	tenantID, err := middleware.MustTenant(r.Context())
	if err != nil {
		writeError(w, http.StatusUnauthorized, "missing_tenant", "tenant no resuelto")
		return
	}
	filter := orderdomain.ListFilter{TenantID: tenantID, Limit: 50}
	if s := r.URL.Query().Get("status"); s != "" {
		st := orderdomain.Status(s)
		filter.Status = &st
	}
	if lim := r.URL.Query().Get("limit"); lim != "" {
		var n int32
		if _, err := fmt.Sscanf(lim, "%d", &n); err == nil && n > 0 && n <= 100 {
			filter.Limit = n
		}
	}
	if ca := r.URL.Query().Get("cursor_created_at"); ca != "" {
		if t, err := time.Parse(time.RFC3339, ca); err == nil {
			filter.CursorCreatedAt = &t
		}
	}
	if cid := r.URL.Query().Get("cursor_id"); cid != "" {
		if _, err := uuid.Parse(cid); err == nil {
			oid := shared.OrderID(cid)
			filter.CursorID = &oid
		}
	}
	orders, err := h.svc.ListOrders(r.Context(), filter)
	if err != nil {
		h.handleError(w, err)
		return
	}
	out := make([]orderSummaryResponse, 0, len(orders))
	for _, o := range orders {
		out = append(out, toOrderSummary(o))
	}
	writeJSON(w, http.StatusOK, map[string]any{"orders": out})
}

func (h *Handler) get(w http.ResponseWriter, r *http.Request) {
	tenantID, err := middleware.MustTenant(r.Context())
	if err != nil {
		writeError(w, http.StatusUnauthorized, "missing_tenant", "tenant no resuelto")
		return
	}
	orderID := chi.URLParam(r, "orderID")
	if _, err := uuid.Parse(orderID); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_order_id", "order_id inválido")
		return
	}
	o, err := h.svc.GetOrder(r.Context(), tenantID, orderID)
	if err != nil {
		h.handleError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, toOrderDetail(o))
}

func (h *Handler) confirmPayment(w http.ResponseWriter, r *http.Request) {
	tenantID, err := middleware.MustTenant(r.Context())
	if err != nil {
		writeError(w, http.StatusUnauthorized, "missing_tenant", "tenant no resuelto")
		return
	}
	merchantUserID, err := middleware.MustUserID(r.Context())
	if err != nil {
		writeError(w, http.StatusUnauthorized, "missing_user", "usuario no autenticado")
		return
	}
	orderID := chi.URLParam(r, "orderID")
	if _, err := uuid.Parse(orderID); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_order_id", "order_id inválido")
		return
	}
	o, err := h.svc.ConfirmPayment(r.Context(), tenantID, orderID, merchantUserID)
	if err != nil {
		h.handleError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, toOrderDetail(o))
}

func (h *Handler) rejectPayment(w http.ResponseWriter, r *http.Request) {
	tenantID, err := middleware.MustTenant(r.Context())
	if err != nil {
		writeError(w, http.StatusUnauthorized, "missing_tenant", "tenant no resuelto")
		return
	}
	merchantUserID, err := middleware.MustUserID(r.Context())
	if err != nil {
		writeError(w, http.StatusUnauthorized, "missing_user", "usuario no autenticado")
		return
	}
	orderID := chi.URLParam(r, "orderID")
	if _, err := uuid.Parse(orderID); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_order_id", "order_id inválido")
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, maxBodyBytes)
	var req rejectRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_json", "cuerpo inválido")
		return
	}
	if err := h.validate.Struct(req); err != nil {
		writeError(w, http.StatusUnprocessableEntity, "validation_error", err.Error())
		return
	}
	o, err := h.svc.RejectPayment(r.Context(), tenantID, orderID, merchantUserID, req.Reason)
	if err != nil {
		h.handleError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, toOrderDetail(o))
}

func (h *Handler) cancel(w http.ResponseWriter, r *http.Request) {
	tenantID, err := middleware.MustTenant(r.Context())
	if err != nil {
		writeError(w, http.StatusUnauthorized, "missing_tenant", "tenant no resuelto")
		return
	}
	orderID := chi.URLParam(r, "orderID")
	if _, err := uuid.Parse(orderID); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_order_id", "order_id inválido")
		return
	}
	o, err := h.svc.GetOrder(r.Context(), tenantID, orderID)
	if err != nil {
		h.handleError(w, err)
		return
	}
	if err := o.Cancel(time.Now()); err != nil {
		h.handleError(w, err)
		return
	}
	// Persistir
	writeJSON(w, http.StatusOK, toOrderDetail(o))
}

func (h *Handler) handleError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, context.Canceled):
		writeError(w, 499, "client_closed_request", "cliente cerró la conexión")
	case errors.Is(err, context.DeadlineExceeded):
		writeError(w, http.StatusGatewayTimeout, "timeout", "el servidor tardó demasiado")
	case errors.Is(err, orderdomain.ErrNotFound):
		writeError(w, http.StatusNotFound, "order_not_found", "orden no encontrada")
	case errors.Is(err, orderdomain.ErrInvalidTransition):
		writeError(w, http.StatusConflict, "invalid_transition", err.Error())
	case errors.Is(err, orderdomain.ErrDropshippingDisabled):
		writeError(w, http.StatusUnprocessableEntity, "dropshipping_disabled",
			"el modo dropshipping no está disponible aún; usa sale_mode=own_inventory")
	case errors.Is(err, paymentmethod.ErrDisabled),
		errors.Is(err, paymentmethod.ErrNotFound),
		errors.Is(err, paymentmethod.ErrNoMethods):
		writeError(w, http.StatusUnprocessableEntity, "payment_method_unavailable", err.Error())
	case errors.Is(err, orderdomain.ErrInvalidTenant),
		errors.Is(err, orderdomain.ErrInvalidCustomer),
		errors.Is(err, orderdomain.ErrInvalidItems),
		errors.Is(err, orderdomain.ErrMixedCurrencies),
		errors.Is(err, orderdomain.ErrInvalidTotal),
		errors.Is(err, orderdomain.ErrProviderRequired),
		errors.Is(err, orderdomain.ErrSplitInvariant),
		errors.Is(err, orderdomain.ErrMissingReference),
		errors.Is(err, orderdomain.ErrMissingMerchantID),
		errors.Is(err, orderdomain.ErrMissingRejectionReason):
		writeError(w, http.StatusUnprocessableEntity, "invalid_order", err.Error())
	case errors.Is(err, orderdomain.ErrVersionConflict):
		writeError(w, http.StatusConflict, "version_conflict", "la orden fue modificada")
	default:
		h.logger.Error("order handler: unhandled error", "err", err)
		writeError(w, http.StatusInternalServerError, "internal", "error interno")
	}
}

func toOrderSummary(o *orderdomain.Order) orderSummaryResponse {
	return orderSummaryResponse{
		ID:            o.ID,
		CustomerID:    o.CustomerID.String(),
		Status:        string(o.Status),
		PaymentMethod: o.PaymentMethod,
		TotalCents:    o.Total.AmountCents,
		Currency:      string(o.Currency),
		CreatedAt:     o.CreatedAt.Format(time.RFC3339),
	}
}

func toOrderDetail(o *orderdomain.Order) orderDetailResponse {
	items := make([]orderItemResponse, 0, len(o.Items))
	for _, it := range o.Items {
		items = append(items, orderItemResponse{
			ProductID: it.ProductID.String(),
			SKU:       it.SKU,
			Name:      it.Name,
			Quantity:  it.Quantity,
			UnitPrice: it.UnitPrice.AmountCents,
			Total:     it.TotalPrice.AmountCents,
		})
	}
	d := orderDetailResponse{
		orderSummaryResponse:     toOrderSummary(o),
		Items:                    items,
		PaymentRef:               o.PaymentRef,
		CustomerPaymentReference: o.CustomerPaymentReference,
		CustomerPaymentProofURL:  o.CustomerPaymentProofURL,
		MerchantConfirmedBy:      o.MerchantConfirmedBy,
		MerchantRejectionReason:  o.MerchantRejectionReason,
		UvicoCommission:          o.UvicoCommission.AmountCents,
		MerchantNet:              o.MerchantNet.AmountCents,
	}
	if o.CustomerSubmittedAt != nil {
		d.CustomerSubmittedAt = o.CustomerSubmittedAt.Format(time.RFC3339)
	}
	if o.MerchantConfirmedAt != nil {
		d.MerchantConfirmedAt = o.MerchantConfirmedAt.Format(time.RFC3339)
	}
	return d
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}

func writeError(w http.ResponseWriter, status int, code, msg string) {
	writeJSON(w, status, errorResponse{Error: code, Message: msg})
}
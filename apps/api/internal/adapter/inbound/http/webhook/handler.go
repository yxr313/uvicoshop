package webhook

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"

	apppayment "uvico-api/internal/application/payment"
	orderdomain "uvico-api/internal/domain/order"
	"uvico-api/internal/platform/config"
)

type Handler struct {
	svc    *apppayment.WebhookService
	cfg    *config.Config
	logger *slog.Logger
}

func New(svc *apppayment.WebhookService, cfg *config.Config, logger *slog.Logger) *Handler {
	return &Handler{svc: svc, cfg: cfg, logger: logger}
}

func (h *Handler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Post("/{provider}/payment", h.handlePayment)
	return r
}

type webhookPayload struct {
	EventID  string `json:"event_id"`
	OrderID  string `json:"order_id"`
	Status   string `json:"status"`
	Amount   int64  `json:"amount_cents"`
	Currency string `json:"currency"`
}

func (h *Handler) handlePayment(w http.ResponseWriter, r *http.Request) {
	provider := chi.URLParam(r, "provider")
	gateway, ok := h.resolveGateway(provider)
	if !ok {
		writeError(w, http.StatusNotFound, "unknown_provider", "pasarela desconocida")
		return
	}
	if gateway.WebhookSecret == "" {
		h.logger.Error("webhook: secret no configurado", "provider", provider)
		writeError(w, http.StatusServiceUnavailable, "misconfigured", "webhook secret no configurado")
		return
	}

	body, err := io.ReadAll(http.MaxBytesReader(w, r.Body, 1<<20))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid_body", "no se pudo leer el cuerpo")
		return
	}

	sig := r.Header.Get(gateway.WebhookHeader)
	if sig == "" {
		writeError(w, http.StatusUnauthorized, "missing_signature", "header de firma ausente")
		return
	}

	if !verifyHMAC(body, []byte(gateway.WebhookSecret), sig) {
		h.logger.Warn("webhook: firma inválida", "provider", provider)
		writeError(w, http.StatusUnauthorized, "invalid_signature", "firma inválida")
		return
	}

	var payload webhookPayload
	if err := json.Unmarshal(body, &payload); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_payload", "payload inválido")
		return
	}
	if payload.EventID == "" || payload.OrderID == "" || payload.Status == "" {
		writeError(w, http.StatusBadRequest, "incomplete_payload", "event_id, order_id y status son requeridos")
		return
	}

	targetStatus, err := mapStatus(payload.Status)
	if err != nil {
		writeError(w, http.StatusUnprocessableEntity, "unknown_status", err.Error())
		return
	}

	result, err := h.svc.Process(r.Context(), provider, payload.EventID, payload.OrderID, targetStatus)
	if err != nil {
		h.handleError(w, err, provider)
		return
	}

	h.logger.Info("webhook procesado",
		"provider", provider,
		"event_id", payload.EventID,
		"order_id", payload.OrderID,
		"status", string(targetStatus),
		"already_processed", result.AlreadyProcessed,
	)
	writeJSON(w, http.StatusOK, map[string]any{
		"ok":                true,
		"already_processed": result.AlreadyProcessed,
	})
}

func (h *Handler) resolveGateway(name string) (config.PaymentGatewayConfig, bool) {
	switch name {
	case "wayu":
		return h.cfg.Webhooks.Wayu, h.cfg.Webhooks.Wayu.BaseURL != ""
	case "veripay":
		return h.cfg.Webhooks.Veripay, h.cfg.Webhooks.Veripay.BaseURL != ""
	case "vexpay":
		return h.cfg.Webhooks.VEXPay, h.cfg.Webhooks.VEXPay.BaseURL != ""
	case "payku":
		return h.cfg.Webhooks.Payku, h.cfg.Webhooks.Payku.BaseURL != ""
	case "pabilo":
		return h.cfg.Webhooks.Pabilo, h.cfg.Webhooks.Pabilo.BaseURL != ""
	default:
		return config.PaymentGatewayConfig{}, false
	}
}

func (h *Handler) handleError(w http.ResponseWriter, err error, provider string) {
	switch {
	case errors.Is(err, orderdomain.ErrNotFound):
		writeError(w, http.StatusNotFound, "order_not_found", "orden no encontrada")
	case errors.Is(err, orderdomain.ErrInvalidTransition):
		writeError(w, http.StatusConflict, "invalid_transition", err.Error())
	case errors.Is(err, orderdomain.ErrMissingReference),
		errors.Is(err, orderdomain.ErrInvalidStatus):
		writeError(w, http.StatusUnprocessableEntity, "invalid_webhook", err.Error())
	default:
		h.logger.Error("webhook: error interno", "provider", provider, "err", err)
		writeError(w, http.StatusInternalServerError, "internal", "error interno")
	}
}

func mapStatus(s string) (orderdomain.Status, error) {
	switch s {
	case "approved", "paid", "succeeded", "completed":
		return orderdomain.StatusPaid, nil
	case "rejected", "failed", "declined":
		return orderdomain.StatusFailed, nil
	case "cancelled", "canceled":
		return orderdomain.StatusCancelled, nil
	default:
		return "", errors.New("webhook: status no reconocido: " + s)
	}
}

func verifyHMAC(body, secret []byte, signature string) bool {
	mac := hmac.New(sha256.New, secret)
	mac.Write(body)
	expected := hex.EncodeToString(mac.Sum(nil))
	return hmac.Equal([]byte(expected), []byte(signature))
}

type errorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-store")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}

func writeError(w http.ResponseWriter, status int, code, msg string) {
	writeJSON(w, status, errorResponse{Error: code, Message: msg})
}

// Silence unused import when ctx is only used in interfaces.
var _ = context.Background
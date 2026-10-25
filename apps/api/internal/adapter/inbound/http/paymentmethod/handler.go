package paymentmethod

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"

	"uvico-api/internal/adapter/inbound/http/middleware"
	apppm "uvico-api/internal/application/paymentmethod"
	domain "uvico-api/internal/domain/paymentmethod"
)

const maxPMBodyBytes = 32 * 1024

type Handler struct {
	svc      *apppm.Service
	logger   *slog.Logger
	validate *validator.Validate
}

func New(svc *apppm.Service, logger *slog.Logger) *Handler {
	return &Handler{
		svc:      svc,
		logger:   logger,
		validate: validator.New(validator.WithRequiredStructEnabled()),
	}
}

// MerchantRoutes: /v1/store/payment-methods (autenticado, tenant-scoped)
func (h *Handler) MerchantRoutes() chi.Router {
	r := chi.NewRouter()
	r.Get("/", h.listMine)
	r.Put("/", h.replaceMine)
	return r
}

type methodDTO struct {
	Method   string         `json:"method" validate:"required,oneof=pago_movil zelle usdt transfer cash"`
	Enabled  bool           `json:"enabled"`
	Position int            `json:"position" validate:"gte=0,lte=20"`
	Config   map[string]any `json:"config"`
}

type replaceRequest struct {
	Methods []methodDTO `json:"methods" validate:"required,max=10,dive"`
}

func (h *Handler) replaceMine(w http.ResponseWriter, r *http.Request) {
	tenantID, err := middleware.MustTenant(r.Context())
	if err != nil {
		writeError(w, http.StatusUnauthorized, "missing_tenant", "tenant no resuelto")
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, maxPMBodyBytes)
	var req replaceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_json", "cuerpo inválido")
		return
	}
	if err := h.validate.Struct(req); err != nil {
		writeError(w, http.StatusUnprocessableEntity, "validation_error", err.Error())
		return
	}

	inputs := make([]apppm.MethodInput, 0, len(req.Methods))
	for _, m := range req.Methods {
		inputs = append(inputs, apppm.MethodInput{
			Method:   m.Method,
			Enabled:  m.Enabled,
			Position: m.Position,
			Config:   m.Config,
		})
	}

	if err := h.svc.ReplaceMethods(r.Context(), tenantID, inputs); err != nil {
		h.handleError(w, err)
		return
	}

	methods, _ := h.svc.ListMine(r.Context(), tenantID)
	writeJSON(w, http.StatusOK, map[string]any{"methods": toMethodViews(methods, true)})
}

func (h *Handler) listMine(w http.ResponseWriter, r *http.Request) {
	tenantID, err := middleware.MustTenant(r.Context())
	if err != nil {
		writeError(w, http.StatusUnauthorized, "missing_tenant", "tenant no resuelto")
		return
	}
	methods, err := h.svc.ListMine(r.Context(), tenantID)
	if err != nil {
		h.handleError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"methods": toMethodViews(methods, true)})
}

type methodView struct {
	Method   string         `json:"method"`
	Label    string         `json:"label"`
	Enabled  bool           `json:"enabled"`
	Position int            `json:"position"`
	Config   map[string]any `json:"config"`
}

func toMethodViews(ms []*domain.PaymentMethod, includeDisabled bool) []methodView {
	out := make([]methodView, 0, len(ms))
	for _, m := range ms {
		if !includeDisabled && !m.Enabled {
			continue
		}
		cfg := m.Config
		if !includeDisabled {
			cfg = m.PublicView()
		}
		out = append(out, methodView{
			Method:   string(m.Method),
			Label:    domain.MethodLabel(m.Method),
			Enabled:  m.Enabled,
			Position: m.Position,
			Config:   cfg,
		})
	}
	return out
}

func (h *Handler) handleError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domain.ErrInvalidMethod),
		errors.Is(err, domain.ErrInvalidConfig):
		writeError(w, http.StatusUnprocessableEntity, "invalid_payment_method", err.Error())
	case errors.Is(err, domain.ErrNotFound):
		writeError(w, http.StatusNotFound, "not_found", "método no encontrado")
	default:
		h.logger.Error("paymentmethod handler: unhandled error", "err", err)
		writeError(w, http.StatusInternalServerError, "internal", "error interno")
	}
}

type errorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
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
package billing

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"

	appbilling "uvico-api/internal/application/billing"
	billdomain "uvico-api/internal/domain/billing"
	"uvico-api/internal/domain/shared"
	"uvico-api/internal/platform/config"
)

type Handler struct {
	svc    *appbilling.Service
	cfg    *config.Config
	logger *slog.Logger
}

func New(svc *appbilling.Service, cfg *config.Config, logger *slog.Logger) *Handler {
	return &Handler{svc: svc, cfg: cfg, logger: logger}
}

func (h *Handler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Get("/plans", h.listPlans)
	r.Get("/subscription", h.getSubscription)
	return r
}

func (h *Handler) listPlans(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{
		"plans": h.svc.ListPlans(),
	})
}

func (h *Handler) getSubscription(w http.ResponseWriter, r *http.Request) {
	tenantIDStr := r.Header.Get("X-Tenant-ID")
	if tenantIDStr == "" {
		writeError(w, http.StatusBadRequest, "missing_tenant", "X-Tenant-ID header requerido")
		return
	}
	tenantID, err := shared.ParseTenantID(tenantIDStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid_tenant", err.Error())
		return
	}
	sub, err := h.svc.GetSubscription(r.Context(), tenantID)
	if errors.Is(err, billdomain.ErrNotFound) {
		writeError(w, http.StatusNotFound, "not_found", "suscripción no encontrada")
		return
	}
	if err != nil {
		h.logger.Error("get subscription", "err", err)
		writeError(w, http.StatusInternalServerError, "internal", "error interno")
		return
	}
	writeJSON(w, http.StatusOK, sub)
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}

func writeError(w http.ResponseWriter, status int, code, message string) {
	writeJSON(w, status, map[string]string{
		"error":   code,
		"message": message,
	})
}
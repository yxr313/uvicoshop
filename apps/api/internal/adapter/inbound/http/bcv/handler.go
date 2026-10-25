package bcv

import (
	"encoding/json"
	"log/slog"
	"net/http"

	appstorefront "uvico-api/internal/application/storefront"
)

type Handler struct {
	svc    *appstorefront.Service
	logger *slog.Logger
}

func New(svc *appstorefront.Service, logger *slog.Logger) *Handler {
	return &Handler{svc: svc, logger: logger}
}

func (h *Handler) GetRate(w http.ResponseWriter, r *http.Request) {
	rate, err := h.svc.GetBCVRate(r.Context())
	if err != nil {
		h.logger.Warn("bcv: rate no disponible", "err", err)
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Cache-Control", "no-store")
		w.WriteHeader(http.StatusServiceUnavailable)
		_, _ = w.Write([]byte(`{"error":"bcv_unavailable","message":"tasa BCV temporalmente no disponible"}`))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "public, max-age=300")
	_ = json.NewEncoder(w).Encode(rate)
}
package analytics

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"uvico-api/internal/adapter/inbound/http/middleware"
	"uvico-api/internal/adapter/outbound/postgres"
	"uvico-api/internal/domain/shared"
)

type Handler struct {
	attributionRepo *postgres.AttributionRepo
	orderRepo       *postgres.OrderRepo
	analyticsRepo   *postgres.AnalyticsRepo
}

func New(attr *postgres.AttributionRepo, ord *postgres.OrderRepo, ana *postgres.AnalyticsRepo) *Handler {
	return &Handler{
		attributionRepo: attr,
		orderRepo:       ord,
		analyticsRepo:   ana,
	}
}

func (h *Handler) PrivateRoutes() chi.Router {
	r := chi.NewRouter()
	r.Get("/traffic", h.getTraffic)
	r.Get("/summary", h.getSummary)
	r.Get("/visits", h.getVisits)
	r.Get("/operations", h.getOperations) // <-- NUEVA RUTA
	return r
}

type trackVisitRequest struct {
	TenantID  string `json:"tenant_id"`
	SessionID string `json:"session_id"`
	Path      string `json:"path"`
	Device    string `json:"device"`
	Browser   string `json:"browser"`
	OS        string `json:"os"`
	Referrer  string `json:"referrer"`
	UTMSource string `json:"utm_source"`
}

// Endpoint PÚBLICO para el rastreador
func (h *Handler) TrackVisit(w http.ResponseWriter, r *http.Request) {
	var req trackVisitRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_json", "cuerpo inválido")
		return
	}

	if req.TenantID == "" || req.Path == "" {
		writeError(w, http.StatusBadRequest, "missing_fields", "tenant_id y path son requeridos")
		return
	}

	err := h.analyticsRepo.TrackVisit(r.Context(), postgres.VisitInput{
		TenantID:  shared.TenantID(req.TenantID),
		SessionID: req.SessionID,
		Path:      req.Path,
		Device:    req.Device,
		Browser:   req.Browser,
		OS:        req.OS,
		Referrer:  req.Referrer,
		UTMSource: req.UTMSource,
	})

	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal", "error guardando visita")
		return
	}

	writeJSON(w, http.StatusOK, map[string]bool{"success": true})
}

// Endpoints PRIVADOS
func (h *Handler) getTraffic(w http.ResponseWriter, r *http.Request) {
	tenantID, err := middleware.MustTenant(r.Context())
	if err != nil {
		writeError(w, http.StatusUnauthorized, "missing_tenant", "tenant no resuelto")
		return
	}

	sources, err := h.attributionRepo.GetTrafficSources(r.Context(), tenantID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal", "error interno")
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{"sources": sources})
}

func (h *Handler) getSummary(w http.ResponseWriter, r *http.Request) {
	tenantID, err := middleware.MustTenant(r.Context())
	if err != nil {
		writeError(w, http.StatusUnauthorized, "missing_tenant", "tenant no resuelto")
		return
	}

	summary, err := h.orderRepo.GetAnalyticsSummary(r.Context(), tenantID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal", "error obteniendo resumen")
		return
	}

	topProducts, err := h.orderRepo.GetTopProducts(r.Context(), tenantID, 5)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal", "error obteniendo productos")
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"summary":      summary,
		"top_products": topProducts,
	})
}

func (h *Handler) getVisits(w http.ResponseWriter, r *http.Request) {
	tenantID, err := middleware.MustTenant(r.Context())
	if err != nil {
		writeError(w, http.StatusUnauthorized, "missing_tenant", "tenant no resuelto")
		return
	}

	timeRange := r.URL.Query().Get("range")
	if timeRange == "" {
		timeRange = "30d"
	}

	visits, err := h.analyticsRepo.GetVisitsTimeSeries(r.Context(), tenantID, timeRange)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal", "error obteniendo visitas")
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{"visits": visits})
}

// NUEVO ENDPOINT OPERATIVO
func (h *Handler) getOperations(w http.ResponseWriter, r *http.Request) {
	tenantID, err := middleware.MustTenant(r.Context())
	if err != nil {
		writeError(w, http.StatusUnauthorized, "missing_tenant", "tenant no resuelto")
		return
	}

	stats, err := h.orderRepo.GetOperationalStats(r.Context(), tenantID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal", "error obteniendo estadísticas operativas")
		return
	}

	lowStock, err := h.orderRepo.GetLowStockProducts(r.Context(), tenantID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal", "error obteniendo productos con bajo stock")
		return
	}

	pendingOrders, err := h.orderRepo.GetPendingOrders(r.Context(), tenantID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal", "error obteniendo pedidos pendientes")
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"stats":          stats,
		"low_stock":      lowStock,
		"pending_orders": pendingOrders,
	})
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}

func writeError(w http.ResponseWriter, status int, code, msg string) {
	writeJSON(w, status, map[string]string{"error": code, "message": msg})
}

package catalog

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"

	appcatalog "uvico-api/internal/application/catalog"
	catalogdomain "uvico-api/internal/domain/catalog"
	"uvico-api/internal/domain/shared"
	"uvico-api/internal/adapter/inbound/http/middleware"
)

type Handler struct {
	svc      *appcatalog.Service
	logger   *slog.Logger
	validate *validator.Validate
}

func New(svc *appcatalog.Service, logger *slog.Logger) *Handler {
	return &Handler{
		svc:      svc,
		logger:   logger,
		validate: validator.New(validator.WithRequiredStructEnabled()),
	}
}

func (h *Handler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Post("/", h.create)
	r.Get("/", h.list)
	r.Get("/{productID}", h.get)
	r.Patch("/{productID}", h.update)
	r.Post("/{productID}/publish", h.publish)
	r.Post("/{productID}/archive", h.archive)
	return r
}

func (h *Handler) create(w http.ResponseWriter, r *http.Request) {
	tenantID, err := middleware.MustTenant(r.Context())
	if err != nil {
		writeError(w, http.StatusUnauthorized, "missing_tenant", "tenant no resuelto")
		return
	}

	var req createProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_json", "cuerpo inválido")
		return
	}
	if err := h.validate.Struct(req); err != nil {
		writeError(w, http.StatusUnprocessableEntity, "validation_error", err.Error())
		return
	}

	var providerID *shared.ProviderID
	if req.ProviderID != nil {
		pid := shared.ProviderID(*req.ProviderID)
		providerID = &pid
	}

	p, err := h.svc.CreateProduct(r.Context(), appcatalog.CreateProductInput{
		TenantID:     tenantID,
		SaleMode:     catalogdomain.SaleMode(req.SaleMode),
		ProviderID:   providerID,
		SKU:          req.SKU,
		Slug:         req.Slug,
		Name:         req.Name,
		Description:  req.Description,
		PriceCents:   req.PriceCents,
		InitialStock: req.InitialStock,
	})
	if err != nil {
		h.handleServiceError(w, err, "create product")
		return
	}
	writeJSON(w, http.StatusCreated, toProductResponse(p))
}

func (h *Handler) list(w http.ResponseWriter, r *http.Request) {
	tenantID, err := middleware.MustTenant(r.Context())
	if err != nil {
		writeError(w, http.StatusUnauthorized, "missing_tenant", "tenant no resuelto")
		return
	}

	filter := catalogdomain.ListFilter{TenantID: tenantID, Limit: 50}
	if s := r.URL.Query().Get("status"); s != "" {
		st := catalogdomain.Status(s)
		filter.Status = &st
	}
	if lim := r.URL.Query().Get("limit"); lim != "" {
		var n int32
		if _, err := jsonNumberParse(lim, &n); err == nil {
			filter.Limit = n
		}
	}

	products, err := h.svc.ListProducts(r.Context(), filter)
	if err != nil {
		h.handleServiceError(w, err, "list products")
		return
	}

	resp := make([]productResponse, 0, len(products))
	for _, p := range products {
		resp = append(resp, toProductResponse(p))
	}
	writeJSON(w, http.StatusOK, map[string]any{"products": resp})
}

func (h *Handler) get(w http.ResponseWriter, r *http.Request) {
	tenantID, err := middleware.MustTenant(r.Context())
	if err != nil {
		writeError(w, http.StatusUnauthorized, "missing_tenant", "tenant no resuelto")
		return
	}
	productID, err := shared.ParseProductID(chi.URLParam(r, "productID"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid_product_id", err.Error())
		return
	}
	p, err := h.svc.GetProduct(r.Context(), tenantID, productID)
	if err != nil {
		h.handleServiceError(w, err, "get product")
		return
	}
	writeJSON(w, http.StatusOK, toProductResponse(p))
}

func (h *Handler) update(w http.ResponseWriter, r *http.Request) {
	tenantID, err := middleware.MustTenant(r.Context())
	if err != nil {
		writeError(w, http.StatusUnauthorized, "missing_tenant", "tenant no resuelto")
		return
	}
	productID, err := shared.ParseProductID(chi.URLParam(r, "productID"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid_product_id", err.Error())
		return
	}

	var req updateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_json", "cuerpo inválido")
		return
	}
	if err := h.validate.Struct(req); err != nil {
		writeError(w, http.StatusUnprocessableEntity, "validation_error", err.Error())
		return
	}

	p, err := h.svc.UpdateProduct(r.Context(), appcatalog.UpdateProductInput{
		TenantID:    tenantID,
		ProductID:   productID,
		Name:        req.Name,
		Description: req.Description,
		Slug:        req.Slug,
		PriceCents:  req.PriceCents,
	})
	if err != nil {
		h.handleServiceError(w, err, "update product")
		return
	}
	writeJSON(w, http.StatusOK, toProductResponse(p))
}

func (h *Handler) publish(w http.ResponseWriter, r *http.Request) {
	tenantID, err := middleware.MustTenant(r.Context())
	if err != nil {
		writeError(w, http.StatusUnauthorized, "missing_tenant", "tenant no resuelto")
		return
	}
	productID, err := shared.ParseProductID(chi.URLParam(r, "productID"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid_product_id", err.Error())
		return
	}
	p, err := h.svc.PublishProduct(r.Context(), tenantID, productID)
	if err != nil {
		h.handleServiceError(w, err, "publish product")
		return
	}
	writeJSON(w, http.StatusOK, toProductResponse(p))
}

func (h *Handler) archive(w http.ResponseWriter, r *http.Request) {
	tenantID, err := middleware.MustTenant(r.Context())
	if err != nil {
		writeError(w, http.StatusUnauthorized, "missing_tenant", "tenant no resuelto")
		return
	}
	productID, err := shared.ParseProductID(chi.URLParam(r, "productID"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid_product_id", err.Error())
		return
	}
	p, err := h.svc.ArchiveProduct(r.Context(), tenantID, productID)
	if err != nil {
		h.handleServiceError(w, err, "archive product")
		return
	}
	writeJSON(w, http.StatusOK, toProductResponse(p))
}

// handleServiceError mapea errores de dominio a códigos HTTP. Los errores
// internos se loggean y se devuelven como 500 sin detalles.
func (h *Handler) handleServiceError(w http.ResponseWriter, err error, op string) {
	switch {
	case errors.Is(err, catalogdomain.ErrNotFound):
		writeError(w, http.StatusNotFound, "not_found", "producto no encontrado")
	case errors.Is(err, catalogdomain.ErrSlugTaken):
		writeError(w, http.StatusConflict, "slug_taken", "el slug ya está en uso")
	case errors.Is(err, catalogdomain.ErrSKUTaken):
		writeError(w, http.StatusConflict, "sku_taken", "el sku ya está en uso")
	case errors.Is(err, catalogdomain.ErrVersionConflict):
		writeError(w, http.StatusConflict, "version_conflict", "el producto fue modificado; recarga y reintenta")
	case errors.Is(err, catalogdomain.ErrArchived):
		writeError(w, http.StatusConflict, "archived", "el producto está archivado")
	case errors.Is(err, catalogdomain.ErrInvalidName),
		errors.Is(err, catalogdomain.ErrInvalidSlug),
		errors.Is(err, catalogdomain.ErrInvalidSKU),
		errors.Is(err, catalogdomain.ErrInvalidPrice),
		errors.Is(err, catalogdomain.ErrInvalidSaleMode),
		errors.Is(err, catalogdomain.ErrInvalidStatus),
		errors.Is(err, catalogdomain.ErrProviderReq):
		writeError(w, http.StatusUnprocessableEntity, "invalid_product", err.Error())
	default:
		h.logger.Error("catalog handler: unhandled error", "op", op, "err", err)
		writeError(w, http.StatusInternalServerError, "internal", "error interno")
	}
}

func toProductResponse(p *catalogdomain.Product) productResponse {
	var providerID *string
	if p.ProviderID != nil {
		s := p.ProviderID.String()
		providerID = &s
	}
	return productResponse{
		ID:          p.ID.String(),
		TenantID:    p.TenantID.String(),
		ProviderID:  providerID,
		SaleMode:    string(p.SaleMode),
		SKU:         p.SKU,
		Slug:        p.Slug,
		Name:        p.Name,
		Description: p.Description,
		PriceCents:  p.Price.AmountCents,
		Currency:    string(p.Price.Currency),
		Status:      string(p.Status),
		CreatedAt:   p.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   p.UpdatedAt.Format(time.RFC3339),
		Version:     p.Version,
	}
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}

func writeError(w http.ResponseWriter, status int, code, msg string) {
	writeJSON(w, status, errorResponse{Error: code, Message: msg})
}

func jsonNumberParse(s string, out *int32) (int, error) {
	var n int
	for _, r := range s {
		if r < '0' || r > '9' {
			return 0, errors.New("no numérico")
		}
		n = n*10 + int(r-'0')
		if n > 100 {
			n = 100
		}
	}
	*out = int32(n)
	return 1, nil
}
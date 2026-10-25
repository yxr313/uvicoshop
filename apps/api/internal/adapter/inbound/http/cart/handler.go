package cart

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"context"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"

	"uvico-api/internal/adapter/inbound/http/middleware"
	appcart "uvico-api/internal/application/cart"
	domain "uvico-api/internal/domain/cart"
	"uvico-api/internal/domain/shared"
)

const maxCartBodyBytes = 32 * 1024

type Handler struct {
	svc      *appcart.Service
	logger   *slog.Logger
	validate *validator.Validate
}

func New(svc *appcart.Service, logger *slog.Logger) *Handler {
	v := validator.New(validator.WithRequiredStructEnabled())
	v.RegisterValidation("uuid_any", func(fl validator.FieldLevel) bool {
		s, ok := fl.Field().Interface().(string)
		if !ok || s == "" {
			return false
		}
		_, err := uuid.Parse(s)
		return err == nil
	})
	return &Handler{svc: svc, logger: logger, validate: v}
}

func (h *Handler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Use(h.requireSession)
	r.Get("/", h.get)
	r.Delete("/", h.clear)
	r.Post("/items", h.addItem)
	r.Patch("/items/{productID}", h.setQuantity)
	r.Delete("/items/{productID}", h.removeItem)
	return r
}

type ctxKey int

const ctxKeySession ctxKey = iota

func (h *Handler) requireSession(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sid := r.Header.Get("X-Session-ID")
		if sid == "" {
			sid = appcart.GenerateSessionID()
		}
		if len(sid) < 8 || len(sid) > 128 {
			writeError(w, http.StatusBadRequest, "invalid_session_id",
				"X-Session-ID debe tener 8..128 caracteres")
			return
		}
		ctx := contextWithSession(r.Context(), sid)
		w.Header().Set("X-Session-ID", sid)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (h *Handler) get(w http.ResponseWriter, r *http.Request) {
	tenantID, sid, ok := h.ctxValues(w, r)
	if !ok {
		return
	}
	c, err := h.svc.Get(r.Context(), tenantID, sid)
	if err != nil {
		h.handleError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, c)
}

func (h *Handler) clear(w http.ResponseWriter, r *http.Request) {
	tenantID, sid, ok := h.ctxValues(w, r)
	if !ok {
		return
	}
	if err := h.svc.Clear(r.Context(), tenantID, sid); err != nil {
		h.handleError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

type addRequest struct {
	ProductID string `json:"product_id" validate:"required,uuid_any"`
	Quantity  int32  `json:"quantity" validate:"required,gt=0,lte=1000"`
}

func (h *Handler) addItem(w http.ResponseWriter, r *http.Request) {
	tenantID, sid, ok := h.ctxValues(w, r)
	if !ok {
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, maxCartBodyBytes)
	var req addRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_json", "cuerpo inválido")
		return
	}
	if err := h.validate.Struct(req); err != nil {
		writeError(w, http.StatusUnprocessableEntity, "validation_error", err.Error())
		return
	}

	c, err := h.svc.AddItem(r.Context(), appcart.AddInput{
		TenantID:  tenantID,
		SessionID: sid,
		ProductID: shared.ProductID(req.ProductID),
		Quantity:  req.Quantity,
	})
	if err != nil {
		h.handleError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, c)
}

type setQtyRequest struct {
	Quantity int32 `json:"quantity" validate:"required,gt=0,lte=1000"`
}

func (h *Handler) setQuantity(w http.ResponseWriter, r *http.Request) {
	tenantID, sid, ok := h.ctxValues(w, r)
	if !ok {
		return
	}
	pid := chi.URLParam(r, "productID")
	if _, err := uuid.Parse(pid); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_product_id", "product_id inválido")
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, maxCartBodyBytes)
	var req setQtyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_json", "cuerpo inválido")
		return
	}
	if err := h.validate.Struct(req); err != nil {
		writeError(w, http.StatusUnprocessableEntity, "validation_error", err.Error())
		return
	}
	c, err := h.svc.SetQuantity(r.Context(), tenantID, sid, shared.ProductID(pid), req.Quantity)
	if err != nil {
		h.handleError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, c)
}

func (h *Handler) removeItem(w http.ResponseWriter, r *http.Request) {
	tenantID, sid, ok := h.ctxValues(w, r)
	if !ok {
		return
	}
	pid := chi.URLParam(r, "productID")
	if _, err := uuid.Parse(pid); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_product_id", "product_id inválido")
		return
	}
	c, err := h.svc.RemoveItem(r.Context(), tenantID, sid, shared.ProductID(pid))
	if err != nil {
		h.handleError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, c)
}

func (h *Handler) ctxValues(w http.ResponseWriter, r *http.Request) (shared.TenantID, string, bool) {
	tenantID, err := middleware.MustTenant(r.Context())
	if err != nil {
		writeError(w, http.StatusUnauthorized, "missing_tenant", "tenant no resuelto")
		return "", "", false
	}
	sid, _ := sessionFromContext(r.Context())
	if sid == "" {
		writeError(w, http.StatusBadRequest, "missing_session", "session requerida")
		return "", "", false
	}
	return tenantID, sid, true
}

func (h *Handler) handleError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domain.ErrEmptySessionID):
		writeError(w, http.StatusBadRequest, "missing_session", "session requerida")
	case errors.Is(err, domain.ErrItemNotFound):
		writeError(w, http.StatusNotFound, "item_not_found", "item no encontrado")
	case errors.Is(err, domain.ErrProductNotFound):
		writeError(w, http.StatusNotFound, "product_not_found", "producto no encontrado o no publicado")
	case errors.Is(err, domain.ErrInvalidQuantity):
		writeError(w, http.StatusUnprocessableEntity, "invalid_quantity", err.Error())
	case errors.Is(err, domain.ErrPriceChanged):
		writeError(w, http.StatusConflict, "price_changed", err.Error())
	default:
		h.logger.Error("cart handler: unhandled error", "err", err)
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

func contextWithSession(ctx context.Context, sid string) context.Context {
	return context.WithValue(ctx, ctxKeySession, sid)
}

func sessionFromContext(ctx context.Context) (string, bool) {
	v, ok := ctx.Value(ctxKeySession).(string)
	return v, ok && v != ""
}
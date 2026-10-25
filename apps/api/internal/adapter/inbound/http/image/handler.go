package image

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"

	"uvico-api/internal/adapter/inbound/http/middleware"
	appimage "uvico-api/internal/application/image"
	"uvico-api/internal/domain/cart"
	"uvico-api/internal/domain/image"
	"uvico-api/internal/domain/shared"
)

const maxImageBodyBytes = 32 * 1024

type Handler struct {
	svc      *appimage.Service
	logger   *slog.Logger
	validate *validator.Validate
}

func New(svc *appimage.Service, logger *slog.Logger) *Handler {
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

// UploadRoutes monta POST /uploads/presigned.
func (h *Handler) UploadRoutes() chi.Router {
	r := chi.NewRouter()
	r.Post("/presigned", h.presign)
	return r
}

// ImageRoutes monta /products/{productID}/images/*.
// El productID viene del URL param inyectado por el router padre.
// Este handler resuelve el param en el middleware interno para que
// cada subhandler lo tenga en el contexto.
func (h *Handler) ImageRoutes() chi.Router {
	r := chi.NewRouter()
	r.Use(h.resolveProduct())
	r.Get("/", h.list)
	r.Post("/", h.attach)
	r.Patch("/reorder", h.reorder)
	r.Delete("/{imageID}", h.deleteImage)
	return r
}

type ctxKey int

const ctxKeyProduct ctxKey = iota

func productFromContext(ctx context.Context) (shared.ProductID, bool) {
	v, ok := ctx.Value(ctxKeyProduct).(shared.ProductID)
	return v, ok && !v.IsZero()
}

func (h *Handler) resolveProduct() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			raw := chi.URLParam(r, "productID")
			if raw == "" {
				writeError(w, http.StatusBadRequest, "missing_product_id", "product_id requerido en la URL")
				return
			}
			if _, err := uuid.Parse(raw); err != nil {
				writeError(w, http.StatusBadRequest, "invalid_product_id", "product_id inválido: "+raw)
				return
			}
			ctx := context.WithValue(r.Context(), ctxKeyProduct, shared.ProductID(raw))
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// ─── POST /uploads/presigned ──────────────────────────────────────

func (h *Handler) presign(w http.ResponseWriter, r *http.Request) {
	tenantID, err := middleware.MustTenant(r.Context())
	if err != nil {
		writeError(w, http.StatusUnauthorized, "missing_tenant", "tenant no resuelto")
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, maxImageBodyBytes)
	var req struct {
		ContentType string `json:"content_type" validate:"required,oneof=image/jpeg image/png image/webp image/avif"`
		FileSize    int64  `json:"file_size" validate:"required,gt=0"`
		ProductID   string `json:"product_id" validate:"required,uuid_any"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_json", "cuerpo inválido")
		return
	}
	if err := h.validate.Struct(req); err != nil {
		writeError(w, http.StatusUnprocessableEntity, "validation_error", err.Error())
		return
	}

	res, err := h.svc.CreatePresignedUpload(r.Context(), appimage.PresignInput{
		TenantID:    tenantID,
		ProductID:   shared.ProductID(req.ProductID),
		ContentType: req.ContentType,
		FileSize:    req.FileSize,
	})
	if err != nil {
		h.handleError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, presignResponse{
		UploadURL:   res.UploadURL,
		StoragePath: res.StoragePath,
		ExpiresAt:   res.ExpiresAt.Format("2006-01-02T15:04:05Z07:00"),
	})
}

// ─── POST /products/{productID}/images ────────────────────────────

func (h *Handler) attach(w http.ResponseWriter, r *http.Request) {
	tenantID, err := middleware.MustTenant(r.Context())
	if err != nil {
		writeError(w, http.StatusUnauthorized, "missing_tenant", "tenant no resuelto")
		return
	}
	productID, ok := productFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusBadRequest, "missing_product_id", "product_id no en contexto")
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, maxImageBodyBytes)
	var req attachRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_json", "cuerpo inválido")
		return
	}
	if err := h.validate.Struct(req); err != nil {
		writeError(w, http.StatusUnprocessableEntity, "validation_error", err.Error())
		return
	}

	img, err := h.svc.AttachImage(r.Context(), appimage.AttachInput{
		TenantID:    tenantID,
		ProductID:   productID,
		StoragePath: req.StoragePath,
		AltText:     req.AltText,
	})
	if err != nil {
		h.handleError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, toImageResponse(img))
}

// ─── GET /products/{productID}/images ─────────────────────────────

func (h *Handler) list(w http.ResponseWriter, r *http.Request) {
	tenantID, err := middleware.MustTenant(r.Context())
	if err != nil {
		writeError(w, http.StatusUnauthorized, "missing_tenant", "tenant no resuelto")
		return
	}
	productID, ok := productFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusBadRequest, "missing_product_id", "product_id no en contexto")
		return
	}
	images, err := h.svc.List(r.Context(), tenantID, productID)
	if err != nil {
		h.handleError(w, err)
		return
	}
	out := make([]imageResponse, 0, len(images))
	for _, img := range images {
		out = append(out, toImageResponse(img))
	}
	writeJSON(w, http.StatusOK, map[string]any{"images": out})
}

// ─── DELETE /products/{productID}/images/{imageID} ────────────────

func (h *Handler) deleteImage(w http.ResponseWriter, r *http.Request) {
	tenantID, err := middleware.MustTenant(r.Context())
	if err != nil {
		writeError(w, http.StatusUnauthorized, "missing_tenant", "tenant no resuelto")
		return
	}
	productID, ok := productFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusBadRequest, "missing_product_id", "product_id no en contexto")
		return
	}
	imageID := chi.URLParam(r, "imageID")
	if _, err := uuid.Parse(imageID); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_image_id", "image_id inválido")
		return
	}
	if err := h.svc.Delete(r.Context(), tenantID, productID, imageID); err != nil {
		h.handleError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// ─── PATCH /products/{productID}/images/reorder ───────────────────

func (h *Handler) reorder(w http.ResponseWriter, r *http.Request) {
	tenantID, err := middleware.MustTenant(r.Context())
	if err != nil {
		writeError(w, http.StatusUnauthorized, "missing_tenant", "tenant no resuelto")
		return
	}
	productID, ok := productFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusBadRequest, "missing_product_id", "product_id no en contexto")
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, maxImageBodyBytes)
	var req reorderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_json", "cuerpo inválido")
		return
	}
	if err := h.validate.Struct(req); err != nil {
		writeError(w, http.StatusUnprocessableEntity, "validation_error", err.Error())
		return
	}
	if err := h.svc.Reorder(r.Context(), tenantID, productID, req.Order); err != nil {
		h.handleError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) handleError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, image.ErrProductNotFound),
		errors.Is(err, cart.ErrProductNotFound):
		writeError(w, http.StatusNotFound, "product_not_found", "producto no encontrado")
	case errors.Is(err, image.ErrNotFound):
		writeError(w, http.StatusNotFound, "image_not_found", "imagen no encontrada")
	case errors.Is(err, image.ErrInvalidMIMEType),
		errors.Is(err, image.ErrFileTooLarge),
		errors.Is(err, image.ErrInvalidPath),
		errors.Is(err, image.ErrMaxImages),
		errors.Is(err, image.ErrPositionConflict):
		writeError(w, http.StatusUnprocessableEntity, "invalid_image", err.Error())
	default:
		h.logger.Error("image handler: unhandled error", "err", err)
		writeError(w, http.StatusInternalServerError, "internal", "error interno")
	}
}

func toImageResponse(img *image.Image) imageResponse {
	return imageResponse{
		ID:        img.ID,
		PublicURL: img.StoragePath,
		AltText:   img.AltText,
		Position:  img.Position,
		CreatedAt: img.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
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
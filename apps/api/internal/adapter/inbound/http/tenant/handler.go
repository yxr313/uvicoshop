package tenant

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"

	"uvico-api/internal/adapter/inbound/http/middleware"
	apptenant "uvico-api/internal/application/tenant"
)

type Handler struct {
	svc      *apptenant.Service
	jwt      *middleware.JWTService
	logger   *slog.Logger
	validate *validator.Validate
}

func New(svc *apptenant.Service, jwt *middleware.JWTService, logger *slog.Logger) *Handler {
	v := validator.New(validator.WithRequiredStructEnabled())

	// Registro de validaciones custom. Solo en este constructor para no
	// afectar otros handlers que no las necesitan.
	v.RegisterValidation("alphanumdash", func(fl validator.FieldLevel) bool {
		s, ok := fl.Field().Interface().(string)
		if !ok {
			return false
		}
		for _, r := range s {
			switch {
			case r >= 'a' && r <= 'z':
			case r >= '0' && r <= '9':
			case r == '-':
			default:
				return false
			}
		}
		return true
	})

	return &Handler{
		svc:      svc,
		jwt:      jwt,
		logger:   logger,
		validate: v,
	}
}

func (h *Handler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Post("/", h.onboard)
	return r
}

type onboardRequest struct {
	Name     string `json:"name" validate:"required,min=1,max=120"`
	Slug     string `json:"slug" validate:"required,min=1,max=60,lowercase,alphanumdash"`
	PlanID   string `json:"plan_id" validate:"required,oneof=founder essential professional enterprise"`
	Interval string `json:"interval" validate:"required,oneof=monthly annual"`
}

type onboardResponse struct {
	TenantID       string `json:"tenant_id"`
	SubscriptionID string `json:"subscription_id"`
	Slug           string `json:"slug"`
	TrialEndsAt    string `json:"trial_ends_at"`
}

func (h *Handler) onboard(w http.ResponseWriter, r *http.Request) {
	userID, err := middleware.MustUserID(r.Context())
	if err != nil {
		writeError(w, http.StatusUnauthorized, "missing_user", "usuario no autenticado")
		return
	}

	if _, hasTenant := middleware.TenantIDFromContext(r.Context()); hasTenant {
		writeError(w, http.StatusConflict, "tenant_exists", "el usuario ya tiene un tenant")
		return
	}

	var req onboardRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_json", "cuerpo inválido")
		return
	}
	if err := h.validate.Struct(req); err != nil {
		writeError(w, http.StatusUnprocessableEntity, "validation_error", err.Error())
		return
	}

	out, err := h.svc.Onboard(r.Context(), apptenant.OnboardInput{
		OwnerUserID: userID,
		Name:        req.Name,
		Slug:        req.Slug,
		PlanID:      req.PlanID,
		Interval:    req.Interval,
		TrialDays:   7,
	})
	if err != nil {
		if errors.Is(err, apptenant.ErrSlugTaken) {
			writeError(w, http.StatusConflict, "slug_taken", "el slug ya está en uso")
			return
		}
		h.logger.Error("onboard failed", "err", err, "user_id", userID)
		writeError(w, http.StatusInternalServerError, "internal", "error interno")
		return
	}

	h.jwt.InvalidateTenantCache(r.Context(), userID)

	writeJSON(w, http.StatusCreated, onboardResponse{
		TenantID:       out.TenantID.String(),
		SubscriptionID: out.SubscriptionID,
		Slug:           out.Slug,
		TrialEndsAt:    out.TrialEndsAt.Format("2006-01-02T15:04:05Z07:00"),
	})
}

type errorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}

func writeError(w http.ResponseWriter, status int, code, msg string) {
	writeJSON(w, status, errorResponse{Error: code, Message: msg})
}
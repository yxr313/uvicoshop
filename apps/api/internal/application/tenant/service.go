package tenant

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"uvico-api/internal/domain/shared"
	"uvico-api/internal/platform/database"
)

type Service struct {
	pool   *database.Pool
	logger *slog.Logger
	now    func() time.Time
}

func NewService(pool *database.Pool, logger *slog.Logger) *Service {
	return &Service{pool: pool, logger: logger, now: time.Now}
}

type OnboardInput struct {
	OwnerUserID  string
	Name         string
	Slug         string
	PlanID       string // founder | essential | professional | enterprise
	Interval     string // monthly | annual
	TrialDays    int
}

type OnboardOutput struct {
	TenantID       shared.TenantID
	SubscriptionID string
	Slug           string
	TrialEndsAt    time.Time
}

// Onboard crea tenant + subscription (trialing) + store_settings + página raíz
// en una única transacción. Si cualquier paso falla, rollback total.
func (s *Service) Onboard(ctx context.Context, in OnboardInput) (*OnboardOutput, error) {
	if in.OwnerUserID == "" {
		return nil, errors.New("tenant: owner_user_id requerido")
	}
	if in.Name == "" || in.Slug == "" {
		return nil, errors.New("tenant: name y slug requeridos")
	}
	if in.PlanID == "" || in.Interval == "" {
		return nil, errors.New("tenant: plan_id y interval requeridos")
	}
	if in.TrialDays <= 0 {
		in.TrialDays = 7
	}

	tenantID := shared.NewTenantID()
	subscriptionID := uuid.Must(uuid.NewV7()).String()
	now := s.now().UTC()
	trialEndsAt := now.AddDate(0, 0, in.TrialDays)

	out := &OnboardOutput{
		TenantID:       tenantID,
		SubscriptionID: subscriptionID,
		Slug:           in.Slug,
		TrialEndsAt:    trialEndsAt,
	}

	err := s.pool.WithTenantTx(ctx, tenantID.String(), func(ctx context.Context, tx pgx.Tx) error {
		// 1. Tenant (trigger crea ledger accounts automáticamente)
		if _, err := tx.Exec(ctx, `
			INSERT INTO tenants (id, name, slug, owner_user_id, status, created_at, updated_at, version)
			VALUES ($1, $2, $3, $4, 'trial', $5, $6, 1)
		`, tenantID.String(), in.Name, in.Slug, in.OwnerUserID, now, now); err != nil {
			if isUniqueViolation(err) {
				return ErrSlugTaken
			}
			return fmt.Errorf("onboard: insert tenant: %w", err)
		}

		// 2. Subscription en trialing
		if _, err := tx.Exec(ctx, `
			INSERT INTO subscriptions (
				id, tenant_id, plan_id, billing_interval, status,
				trial_ends_at, billing_cycles_completed,
				created_at, updated_at, version
			) VALUES ($1, $2, $3, $4, 'trialing', $5, 0, $6, $7, 1)
		`, subscriptionID, tenantID.String(), in.PlanID, in.Interval, trialEndsAt, now, now); err != nil {
			return fmt.Errorf("onboard: insert subscription: %w", err)
		}

		// 3. Store settings default
		if _, err := tx.Exec(ctx, `
			INSERT INTO store_settings (tenant_id, currency_display, updated_at)
			VALUES ($1, 'USD', $2)
		`, tenantID.String(), now); err != nil {
			return fmt.Errorf("onboard: insert store settings: %w", err)
		}

		// 4. Página raíz vacía (Puck inicia con JSON vacío)
		defaultPuck, _ := json.Marshal(map[string]any{
			"root":    map[string]any{"props": map[string]any{}},
			"content": []any{},
		})
		if _, err := tx.Exec(ctx, `
			INSERT INTO store_pages (tenant_id, page_slug, puck_data, published, updated_at)
			VALUES ($1, '/', $2, false, $3)
		`, tenantID.String(), string(defaultPuck), now); err != nil {
			return fmt.Errorf("onboard: insert store page: %w", err)
		}

		// 5. Outbox: tenant.created
		payload, _ := json.Marshal(map[string]any{
			"tenant_id":       tenantID.String(),
			"subscription_id": subscriptionID,
			"plan_id":         in.PlanID,
			"interval":        in.Interval,
			"trial_ends_at":   trialEndsAt,
		})
		if _, err := tx.Exec(ctx, `
			INSERT INTO outbox_events (id, tenant_id, aggregate_type, aggregate_id, event_type, payload, created_at)
			VALUES (gen_random_uuid(), $1, 'tenant', $2, 'tenant.created', $3, $4)
		`, tenantID.String(), tenantID.String(), string(payload), now); err != nil {
			return fmt.Errorf("onboard: insert outbox: %w", err)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	s.logger.Info("tenant onboarded",
		"tenant_id", tenantID.String(),
		"user_id", in.OwnerUserID,
		"slug", in.Slug,
		"plan", in.PlanID,
	)
	return out, nil
}

var ErrSlugTaken = errors.New("tenant: slug ya existe")

func isUniqueViolation(err error) bool {
	type pgErr interface{ SQLState() string }
	var e pgErr
	if errors.As(err, &e) {
		return e.SQLState() == "23505"
	}
	return false
}
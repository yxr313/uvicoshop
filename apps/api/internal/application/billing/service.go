package billing

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"

	billdomain "uvico-api/internal/domain/billing"
	"uvico-api/internal/domain/shared"
	"uvico-api/internal/domain/tenant"
)

type Service struct {
	tenants tenant.Repository
	subs    billdomain.SubscriptionRepository
	plans   map[billdomain.PlanID]billdomain.Plan
	logger  *slog.Logger
}

func NewService(
	tenants tenant.Repository,
	subs billdomain.SubscriptionRepository,
	plans map[billdomain.PlanID]billdomain.Plan,
	logger *slog.Logger,
) *Service {
	return &Service{tenants: tenants, subs: subs, plans: plans, logger: logger}
}

type CreateTenantInput struct {
	Name            string
	Slug            string
	OwnerUserID     string
	PlanID          billdomain.PlanID
	BillingInterval billdomain.BillingInterval
}

type CreateTenantOutput struct {
	Tenant       *tenant.Tenant
	Subscription *billdomain.Subscription
}

// CreateTenantWithTrial crea Tenant + Subscription en estado trialing.
// No hay cobro. La activación ocurre cuando llega el primer pago.
func (s *Service) CreateTenantWithTrial(ctx context.Context, in CreateTenantInput, now time.Time) (*CreateTenantOutput, error) {
	plan, ok := s.plans[in.PlanID]
	if !ok {
		return nil, fmt.Errorf("%w: %q", billdomain.ErrInvalidPlan, in.PlanID)
	}

	exists, err := s.tenants.SlugExists(ctx, in.Slug)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, tenant.ErrSlugTaken
	}

	tenantID := shared.NewTenantID()
	t, err := tenant.New(tenantID, in.Name, in.Slug, in.OwnerUserID, now)
	if err != nil {
		return nil, err
	}

	sub, err := billdomain.StartTrial(uuid.Must(uuid.NewV7()).String(), tenantID, plan, in.BillingInterval, now)
	if err != nil {
		return nil, err
	}

	if err := s.tenants.Create(ctx, t); err != nil {
		return nil, err
	}
	if err := s.subs.Create(ctx, sub); err != nil {
		// Rollback manual: borrar tenant. Sin transacción atómica cross-repo
		// por ahora. Se puede mejorar con WithTenantTx cuando ambos repos
		// compartan la misma tx.
		return nil, fmt.Errorf("create subscription (tenant created, orphan): %w", err)
	}
	return &CreateTenantOutput{Tenant: t, Subscription: sub}, nil
}

func (s *Service) GetSubscription(ctx context.Context, tenantID shared.TenantID) (*billdomain.Subscription, error) {
	return s.subs.GetByTenantID(ctx, tenantID)
}

type PlanView struct {
	ID           billdomain.PlanID `json:"id"`
	MonthlyCents int64             `json:"monthly_cents"`
	AnnualCents  int64             `json:"annual_cents"`
	TrialDays    int               `json:"trial_days"`
}

func (s *Service) ListPlans() []PlanView {
	out := make([]PlanView, 0, len(s.plans))
	for _, p := range s.plans {
		out = append(out, PlanView{
			ID:           p.ID,
			MonthlyCents: p.MonthlyCents,
			AnnualCents:  p.AnnualCents,
			TrialDays:    p.TrialDays,
		})
	}
	return out
}

// Silence unused import in case logger isn't used in some build paths.
var _ = errors.New
package billing

import (
	"errors"
	"fmt"
	"time"

	"uvico-api/internal/domain/shared"
)

type SubscriptionStatus string

const (
	SubTrialing        SubscriptionStatus = "trialing"
	SubActive          SubscriptionStatus = "active"
	SubPastDue         SubscriptionStatus = "past_due"
	SubPaused          SubscriptionStatus = "paused"
	SubPendingDeletion SubscriptionStatus = "pending_deletion"
	SubCancelled       SubscriptionStatus = "cancelled"
	SubDeleted         SubscriptionStatus = "deleted"
)

var (
	ErrInvalidTransition = errors.New("billing: transición inválida")
	ErrNotFound          = errors.New("billing: suscripción no encontrada")
	ErrAlreadyExists     = errors.New("billing: suscripción ya existe para el tenant")
	ErrVersionConflict   = errors.New("billing: version conflict")
	ErrInvalidTenant     = errors.New("billing: tenant requerido")
)

type Subscription struct {
	ID                     string
	TenantID               shared.TenantID
	PlanID                 PlanID
	BillingInterval        BillingInterval
	Status                 SubscriptionStatus
	TrialEndsAt            *time.Time
	CurrentPeriodStartsAt  *time.Time
	CurrentPeriodEndsAt    *time.Time
	BillingCyclesCompleted int32
	PastDueSince           *time.Time
	PausedAt               *time.Time
	DeletionScheduledAt    *time.Time
	CancelledAt            *time.Time
	CreatedAt              time.Time
	UpdatedAt              time.Time
	Version                int64
}

// StartTrial crea una suscripción en estado trialing. Sin cargo.
// El período de billing comienza cuando se activa (post-trial).
func StartTrial(
	id string,
	tenantID shared.TenantID,
	plan Plan,
	interval BillingInterval,
	now time.Time,
) (*Subscription, error) {
	if tenantID.IsZero() {
		return nil, ErrInvalidTenant
	}
	if _, err := plan.PriceFor(interval); err != nil {
		return nil, err
	}
	if plan.TrialDays <= 0 {
		return nil, fmt.Errorf("billing: plan %q no define trial", plan.ID)
	}
	t := now.UTC()
	trialEnd := t.AddDate(0, 0, plan.TrialDays)
	return &Subscription{
		ID:              id,
		TenantID:        tenantID,
		PlanID:          plan.ID,
		BillingInterval: interval,
		Status:          SubTrialing,
		TrialEndsAt:     &trialEnd,
		CreatedAt:       t,
		UpdatedAt:       t,
		Version:         1,
	}, nil
}

// NextAmount devuelve el monto a cobrar por el PRÓXIMO ciclo.
// Aplica promo de founder si BillingCyclesCompleted < FounderPromoMonths.
func (s *Subscription) NextAmount(plans map[PlanID]Plan) (shared.Money, error) {
	plan, ok := plans[s.PlanID]
	if !ok {
		return shared.Money{}, fmt.Errorf("%w: %q", ErrInvalidPlan, s.PlanID)
	}
	if plan.FounderPromoMonths > 0 &&
		s.BillingCyclesCompleted < int32(plan.FounderPromoMonths) {
		return shared.USD(plan.FounderPromoMonthlyCents), nil
	}
	return plan.PriceFor(s.BillingInterval)
}

// Activate se llama cuando el primer pago se confirmó.
// Establece el período de billing (1 mes o 1 año desde now).
func (s *Subscription) Activate(now time.Time) error {
	if s.Status != SubTrialing {
		return fmt.Errorf("%w: %s -> active", ErrInvalidTransition, s.Status)
	}
	t := now.UTC()
	periodEnd := s.computePeriodEnd(t)
	s.Status = SubActive
	s.CurrentPeriodStartsAt = &t
	s.CurrentPeriodEndsAt = &periodEnd
	s.BillingCyclesCompleted = 1
	s.UpdatedAt = t
	s.Version++
	return nil
}

// AdvanceCycle se llama cuando el pago del siguiente ciclo se confirma.
// El período previo debe estar expirado o próximo a expirar.
func (s *Subscription) AdvanceCycle(now time.Time) error {
	if s.Status != SubActive && s.Status != SubPastDue {
		return fmt.Errorf("%w: %s -> advance", ErrInvalidTransition, s.Status)
	}
	t := now.UTC()
	periodEnd := s.computePeriodEnd(t)
	s.Status = SubActive
	s.CurrentPeriodStartsAt = &t
	s.CurrentPeriodEndsAt = &periodEnd
	s.BillingCyclesCompleted++
	s.PastDueSince = nil
	s.UpdatedAt = t
	s.Version++
	return nil
}

// MarkPastDue se llama desde el cron cuando trial vence sin pago o
// el período activo expira sin renovación.
func (s *Subscription) MarkPastDue(now time.Time) error {
	if s.Status == SubPastDue {
		return nil
	}
	if s.Status != SubTrialing && s.Status != SubActive {
		return fmt.Errorf("%w: %s -> past_due", ErrInvalidTransition, s.Status)
	}
	t := now.UTC()
	s.Status = SubPastDue
	s.PastDueSince = &t
	s.UpdatedAt = t
	s.Version++
	return nil
}

// Pause aplica al 6to día de past_due (T&C §9.3).
func (s *Subscription) Pause(now time.Time) error {
	if s.Status == SubPaused {
		return nil
	}
	if s.Status != SubPastDue {
		return fmt.Errorf("%w: %s -> paused", ErrInvalidTransition, s.Status)
	}
	t := now.UTC()
	s.Status = SubPaused
	s.PausedAt = &t
	s.UpdatedAt = t
	s.Version++
	return nil
}

// MarkPendingDeletion aplica al día 25 desde paused (aviso 5 días).
func (s *Subscription) MarkPendingDeletion(now time.Time) error {
	if s.Status == SubPendingDeletion {
		return nil
	}
	if s.Status != SubPaused {
		return fmt.Errorf("%w: %s -> pending_deletion", ErrInvalidTransition, s.Status)
	}
	t := now.UTC()
	scheduled := t.AddDate(0, 0, 5)
	s.Status = SubPendingDeletion
	s.DeletionScheduledAt = &scheduled
	s.UpdatedAt = t
	s.Version++
	return nil
}

// Cancel es voluntario. Efectivo al final del período pagado.
func (s *Subscription) Cancel(now time.Time) error {
	if s.Status == SubCancelled {
		return nil
	}
	if s.Status != SubActive && s.Status != SubTrialing && s.Status != SubPastDue {
		return fmt.Errorf("%w: %s -> cancelled", ErrInvalidTransition, s.Status)
	}
	t := now.UTC()
	s.Status = SubCancelled
	s.CancelledAt = &t
	s.UpdatedAt = t
	s.Version++
	return nil
}

func (s *Subscription) IsInTrial() bool {
	return s.Status == SubTrialing
}

func (s *Subscription) TrialDaysRemaining(now time.Time) int {
	if s.TrialEndsAt == nil || s.Status != SubTrialing {
		return 0
	}
	d := int(s.TrialEndsAt.Sub(now).Hours() / 24)
	if d < 0 {
		return 0
	}
	return d
}

func (s *Subscription) IsFounderPromoActive(plans map[PlanID]Plan) bool {
	plan, ok := plans[s.PlanID]
	if !ok || plan.FounderPromoMonths == 0 {
		return false
	}
	return s.BillingCyclesCompleted < int32(plan.FounderPromoMonths)
}

func (s *Subscription) computePeriodEnd(from time.Time) time.Time {
	switch s.BillingInterval {
	case IntervalMonthly:
		return from.AddDate(0, 1, 0)
	case IntervalAnnual:
		return from.AddDate(1, 0, 0)
	}
	return from
}
package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"

	"uvico-api/internal/domain/billing"
	"uvico-api/internal/domain/shared"
	"uvico-api/internal/platform/database"
)

type SubscriptionRepo struct {
	pool *database.Pool
}

func NewSubscriptionRepo(pool *database.Pool) *SubscriptionRepo {
	return &SubscriptionRepo{pool: pool}
}

func (r *SubscriptionRepo) Create(ctx context.Context, s *billing.Subscription) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO subscriptions (
			id, tenant_id, plan_id, billing_interval, status,
			trial_ends_at, current_period_starts_at, current_period_ends_at,
			billing_cycles_completed, past_due_since, paused_at,
			deletion_scheduled_at, cancelled_at,
			created_at, updated_at, version
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16)
	`, s.ID, string(s.TenantID), string(s.PlanID), string(s.BillingInterval), string(s.Status),
		s.TrialEndsAt, s.CurrentPeriodStartsAt, s.CurrentPeriodEndsAt,
		s.BillingCyclesCompleted, s.PastDueSince, s.PausedAt,
		s.DeletionScheduledAt, s.CancelledAt,
		s.CreatedAt, s.UpdatedAt, s.Version)
	if err != nil {
		return fmt.Errorf("subscription_repo: create: %w", err)
	}
	return nil
}

func (r *SubscriptionRepo) Update(ctx context.Context, s *billing.Subscription) error {
	tag, err := r.pool.Exec(ctx, `
		UPDATE subscriptions
		   SET plan_id = $2,
		       billing_interval = $3,
		       status = $4,
		       trial_ends_at = $5,
		       current_period_starts_at = $6,
		       current_period_ends_at = $7,
		       billing_cycles_completed = $8,
		       past_due_since = $9,
		       paused_at = $10,
		       deletion_scheduled_at = $11,
		       cancelled_at = $12,
		       updated_at = $13,
		       version = $14
		 WHERE id = $1 AND version = $15
	`, s.ID, string(s.PlanID), string(s.BillingInterval), string(s.Status),
		s.TrialEndsAt, s.CurrentPeriodStartsAt, s.CurrentPeriodEndsAt,
		s.BillingCyclesCompleted, s.PastDueSince, s.PausedAt,
		s.DeletionScheduledAt, s.CancelledAt,
		s.UpdatedAt, s.Version, s.Version-1)
	if err != nil {
		return fmt.Errorf("subscription_repo: update: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return billing.ErrVersionConflict
	}
	return nil
}

func (r *SubscriptionRepo) GetByTenantID(ctx context.Context, tenantID shared.TenantID) (*billing.Subscription, error) {
	var s billing.Subscription
	var planID, interval, status string

	err := r.pool.QueryRow(ctx, `
		SELECT id, tenant_id, plan_id, billing_interval, status,
		       trial_ends_at, current_period_starts_at, current_period_ends_at,
		       billing_cycles_completed, past_due_since, paused_at,
		       deletion_scheduled_at, cancelled_at,
		       created_at, updated_at, version
		  FROM subscriptions
		 WHERE tenant_id = $1
	`, string(tenantID)).Scan(
		&s.ID, (*string)(&s.TenantID), &planID, &interval, &status,
		&s.TrialEndsAt, &s.CurrentPeriodStartsAt, &s.CurrentPeriodEndsAt,
		&s.BillingCyclesCompleted, &s.PastDueSince, &s.PausedAt,
		&s.DeletionScheduledAt, &s.CancelledAt,
		&s.CreatedAt, &s.UpdatedAt, &s.Version,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, billing.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("subscription_repo: get: %w", err)
	}
	s.PlanID = billing.PlanID(planID)
	s.BillingInterval = billing.BillingInterval(interval)
	s.Status = billing.SubscriptionStatus(status)
	return &s, nil
}

func (r *SubscriptionRepo) ListExpiredTrials(ctx context.Context, now time.Time, limit int) ([]*billing.Subscription, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id FROM subscriptions
		 WHERE status = 'trialing'
		   AND trial_ends_at <= $1
		 ORDER BY trial_ends_at
		 LIMIT $2
	`, now, limit)
	if err != nil {
		return nil, fmt.Errorf("subscription_repo: list expired trials: %w", err)
	}
	defer rows.Close()

	var out []*billing.Subscription
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		out = append(out, &billing.Subscription{ID: id})
	}
	return out, rows.Err()
}
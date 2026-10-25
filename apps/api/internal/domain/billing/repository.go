package billing

import (
	"context"
	"time"

	"uvico-api/internal/domain/shared"
)

type SubscriptionRepository interface {
	Create(ctx context.Context, s *Subscription) error
	Update(ctx context.Context, s *Subscription) error
	GetByTenantID(ctx context.Context, tenantID shared.TenantID) (*Subscription, error)
	ListExpiredTrials(ctx context.Context, now time.Time, limit int) ([]*Subscription, error)
}
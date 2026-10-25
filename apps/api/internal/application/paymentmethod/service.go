package paymentmethod

import (
	"context"
	"fmt"
	"time"

	domain "uvico-api/internal/domain/paymentmethod"
	"uvico-api/internal/domain/shared"
)

type Service struct {
	repo domain.Repository
	now  func() time.Time
}

func NewService(repo domain.Repository) *Service {
	return &Service{repo: repo, now: time.Now}
}

type MethodInput struct {
	Method   string         `json:"method"`
	Enabled  bool           `json:"enabled"`
	Position int            `json:"position"`
	Config   map[string]any `json:"config"`
}

// ReplaceMethods atomically replaces the tenant's method list.
// Methods not present in the input are disabled (soft delete).
func (s *Service) ReplaceMethods(ctx context.Context, tenantID shared.TenantID, inputs []MethodInput) error {
	methods := make([]*domain.PaymentMethod, 0, len(inputs))
	now := s.now()
	for i, in := range inputs {
		pos := in.Position
		if pos == 0 {
			pos = i
		}
		m, err := domain.New(tenantID, domain.Method(in.Method), pos, domain.Config(in.Config), now)
		if err != nil {
			return err
		}
		m.Enabled = in.Enabled
		methods = append(methods, m)
	}
	return s.repo.Replace(ctx, tenantID, methods)
}

func (s *Service) ListMine(ctx context.Context, tenantID shared.TenantID) ([]*domain.PaymentMethod, error) {
	return s.repo.ListAll(ctx, tenantID)
}

func (s *Service) ListPublic(ctx context.Context, tenantID shared.TenantID) ([]*domain.PaymentMethod, error) {
	return s.repo.ListEnabled(ctx, tenantID)
}

// GetForCheckout returns the specific method the buyer picked, failing
// with ErrDisabled if the merchant hasn't enabled it.
func (s *Service) GetForCheckout(ctx context.Context, tenantID shared.TenantID, method string) (*domain.PaymentMethod, error) {
	m, err := s.repo.GetByMethod(ctx, tenantID, domain.Method(method))
	if err != nil {
		return nil, err
	}
	if !m.Enabled {
		return nil, fmt.Errorf("%w: %s", domain.ErrDisabled, method)
	}
	return m, nil
}
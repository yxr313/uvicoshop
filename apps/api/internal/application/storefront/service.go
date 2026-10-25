package storefront

import (
	"context"
	"errors"
	"fmt"

	"uvico-api/internal/domain/shared"
	domain "uvico-api/internal/domain/storefront"
)

var (
	ErrTenantNotFound   = errors.New("storefront: tenant no encontrado")
	ErrProductNotFound  = errors.New("storefront: producto no encontrado")
	ErrStoreUnavailable = errors.New("storefront: tienda no disponible")
)

type TenantInfo struct {
	ID     shared.TenantID
	Slug   string
	Status string
}

// Repository es implementado por el adaptador de Postgres.
type Repository interface {
	GetTenantBySlug(ctx context.Context, slug string) (*TenantInfo, error)
	GetStoreConfig(ctx context.Context, tenantID shared.TenantID) (*domain.StoreConfig, error)
	ListPublishedProducts(ctx context.Context, tenantID shared.TenantID, limit int) ([]*domain.ProductSummary, error)
	GetPublishedProduct(ctx context.Context, tenantID shared.TenantID, productSlug string) (*domain.ProductDetail, error)
}

// BCVRepository es implementado por el adaptador de Redis.
type BCVRepository interface {
	GetRate(ctx context.Context) (*domain.BCVRate, error)
}

type Service struct {
	repo Repository
	bcv  BCVRepository
}

func NewService(repo Repository, bcv BCVRepository) *Service {
	return &Service{repo: repo, bcv: bcv}
}

func (s *Service) ResolveTenant(ctx context.Context, slug string) (*TenantInfo, error) {
	if slug == "" {
		return nil, ErrTenantNotFound
	}
	tenant, err := s.repo.GetTenantBySlug(ctx, slug)
	if err != nil {
		return nil, err
	}
	switch tenant.Status {
	case "trial", "active", "past_due":
		return tenant, nil
	default:
		return nil, ErrStoreUnavailable
	}
}

func (s *Service) GetConfig(ctx context.Context, tenantID shared.TenantID) (*domain.StoreConfig, error) {
	return s.repo.GetStoreConfig(ctx, tenantID)
}

func (s *Service) ListProducts(ctx context.Context, tenantID shared.TenantID, limit int) ([]*domain.ProductSummary, error) {
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	return s.repo.ListPublishedProducts(ctx, tenantID, limit)
}

func (s *Service) GetProduct(ctx context.Context, tenantID shared.TenantID, productSlug string) (*domain.ProductDetail, error) {
	if productSlug == "" {
		return nil, ErrProductNotFound
	}
	return s.repo.GetPublishedProduct(ctx, tenantID, productSlug)
}

func (s *Service) GetBCVRate(ctx context.Context) (*domain.BCVRate, error) {
	rate, err := s.bcv.GetRate(ctx)
	if err != nil {
		return nil, fmt.Errorf("storefront: bcv rate: %w", err)
	}
	return rate, nil
}
package catalog

import (
	"context"
	"errors"
	"fmt"
	"time"

	"uvico-api/internal/domain/catalog"
	"uvico-api/internal/domain/shared"
)

type Service struct {
	repo catalog.Repository
	now  func() time.Time
}

func NewService(repo catalog.Repository) *Service {
	return &Service{repo: repo, now: time.Now}
}

type CreateProductInput struct {
	TenantID     shared.TenantID
	SaleMode     catalog.SaleMode
	ProviderID   *shared.ProviderID
	SKU          string
	Slug         string
	Name         string
	Description  string
	PriceCents   int64
	InitialStock int64 // solo own_inventory; ignorado en dropshipping
}

func (s *Service) CreateProduct(ctx context.Context, in CreateProductInput) (*catalog.Product, error) {
	p, err := catalog.NewProduct(
		shared.NewProductID(),
		in.TenantID,
		in.ProviderID,
		in.SaleMode,
		in.SKU,
		in.Slug,
		in.Name,
		in.Description,
		shared.USD(in.PriceCents),
		s.now(),
	)
	if err != nil {
		return nil, err
	}

	if in.SaleMode == catalog.SaleOwnInventory && in.InitialStock < 0 {
		return nil, fmt.Errorf("catalog: initial_stock no puede ser negativo")
	}

	if err := s.repo.Create(ctx, p, in.InitialStock); err != nil {
		return nil, err
	}
	return p, nil
}

type UpdateProductInput struct {
	TenantID    shared.TenantID
	ProductID   shared.ProductID
	Name        string
	Description string
	Slug        string
	PriceCents  int64
}

func (s *Service) UpdateProduct(ctx context.Context, in UpdateProductInput) (*catalog.Product, error) {
	p, err := s.repo.GetByID(ctx, in.TenantID, in.ProductID)
	if err != nil {
		return nil, err
	}
	if err := p.Update(in.Name, in.Description, in.Slug, shared.USD(in.PriceCents), s.now()); err != nil {
		return nil, err
	}
	if err := s.repo.Update(ctx, p); err != nil {
		return nil, err
	}
	return p, nil
}

func (s *Service) PublishProduct(ctx context.Context, tenantID shared.TenantID, productID shared.ProductID) (*catalog.Product, error) {
	p, err := s.repo.GetByID(ctx, tenantID, productID)
	if err != nil {
		return nil, err
	}
	if err := p.Publish(s.now()); err != nil {
		return nil, err
	}
	if err := s.repo.Update(ctx, p); err != nil {
		return nil, err
	}
	return p, nil
}

func (s *Service) ArchiveProduct(ctx context.Context, tenantID shared.TenantID, productID shared.ProductID) (*catalog.Product, error) {
	p, err := s.repo.GetByID(ctx, tenantID, productID)
	if err != nil {
		return nil, err
	}
	if err := p.Archive(s.now()); err != nil {
		return nil, err
	}
	if err := s.repo.Update(ctx, p); err != nil {
		return nil, err
	}
	return p, nil
}

func (s *Service) GetProduct(ctx context.Context, tenantID shared.TenantID, productID shared.ProductID) (*catalog.Product, error) {
	return s.repo.GetByID(ctx, tenantID, productID)
}

func (s *Service) ListProducts(ctx context.Context, f catalog.ListFilter) ([]*catalog.Product, error) {
	if f.Limit <= 0 || f.Limit > 100 {
		f.Limit = 50
	}
	return s.repo.List(ctx, f)
}

// ExposedErrors lista los errores de dominio que los handlers pueden mapear a códigos HTTP.
var _ = errors.Is
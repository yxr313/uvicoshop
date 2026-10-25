package catalog

import (
	"fmt"
	"strings"
	"time"
	"unicode/utf8"

	"uvico-api/internal/domain/shared"
)

type SaleMode string

const (
	SaleOwnInventory SaleMode = "own_inventory"
	SaleDropshipping SaleMode = "dropshipping_verified"
)

type Status string

const (
	StatusDraft     Status = "draft"
	StatusPublished Status = "published"
	StatusArchived  Status = "archived"
)

const (
	MaxNameLen  = 120
	MaxDescLen  = 5000
	MaxSlugLen  = 80
	MaxSKULen   = 64
	MaxPriceCent = int64(1_000_000_000_00) // $1B cap defensivo
)

type Product struct {
	ID          shared.ProductID
	TenantID    shared.TenantID
	ProviderID  *shared.ProviderID
	SaleMode    SaleMode
	SKU         string
	Slug        string
	Name        string
	Description string
	Price       shared.Money
	Status      Status
	PublishedAt *time.Time
	ArchivedAt  *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Version     int64
}

// NewProduct crea y valida un producto en estado draft.
func NewProduct(
	id shared.ProductID,
	tenantID shared.TenantID,
	providerID *shared.ProviderID,
	saleMode SaleMode,
	sku, slug, name, description string,
	price shared.Money,
	now time.Time,
) (*Product, error) {
	if id.IsZero() {
		return nil, fmt.Errorf("catalog: product id requerido")
	}
	if tenantID.IsZero() {
		return nil, ErrInvalidTenant
	}
	if err := validateSaleModeProvider(saleMode, providerID); err != nil {
		return nil, err
	}
	name = strings.TrimSpace(name)
	if name == "" || utf8.RuneCountInString(name) > MaxNameLen {
		return nil, ErrInvalidName
	}
	if utf8.RuneCountInString(description) > MaxDescLen {
		return nil, fmt.Errorf("catalog: descripción excede %d caracteres", MaxDescLen)
	}
	if !isValidSlug(slug) {
		return nil, ErrInvalidSlug
	}
	if !isValidSKU(sku) {
		return nil, ErrInvalidSKU
	}
	if err := price.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidPrice, err)
	}
	if price.IsZero() {
		return nil, ErrInvalidPrice
	}
	if price.AmountCents > MaxPriceCent {
		return nil, fmt.Errorf("%w: excede el límite permitido", ErrInvalidPrice)
	}

	t := now.UTC()
	return &Product{
		ID:          id,
		TenantID:    tenantID,
		ProviderID:  providerID,
		SaleMode:    saleMode,
		SKU:         sku,
		Slug:        slug,
		Name:        name,
		Description: description,
		Price:       price,
		Status:      StatusDraft,
		CreatedAt:   t,
		UpdatedAt:   t,
		Version:     1,
	}, nil
}

// Update aplica cambios en name/description/slug/price.
// Solo permitido en draft o published. Nunca en archived.
func (p *Product) Update(name, description, slug string, price shared.Money, now time.Time) error {
	if p.Status == StatusArchived {
		return ErrArchived
	}
	name = strings.TrimSpace(name)
	if name == "" || utf8.RuneCountInString(name) > MaxNameLen {
		return ErrInvalidName
	}
	if utf8.RuneCountInString(description) > MaxDescLen {
		return fmt.Errorf("catalog: descripción excede %d caracteres", MaxDescLen)
	}
	if slug != p.Slug && !isValidSlug(slug) {
		return ErrInvalidSlug
	}
	if err := price.Validate(); err != nil || price.IsZero() {
		return ErrInvalidPrice
	}
	if price.AmountCents > MaxPriceCent {
		return fmt.Errorf("%w: excede el límite permitido", ErrInvalidPrice)
	}
	p.Name = name
	p.Description = description
	p.Slug = slug
	p.Price = price
	p.UpdatedAt = now.UTC()
	p.Version++
	return nil
}

// Publish solo es válido desde draft. Idempotente en published.
func (p *Product) Publish(now time.Time) error {
	if p.Status == StatusPublished {
		return nil
	}
	if p.Status == StatusArchived {
		return ErrArchived
	}
	if p.Status != StatusDraft {
		return fmt.Errorf("%w: %s -> published", ErrInvalidStatus, p.Status)
	}
	t := now.UTC()
	p.Status = StatusPublished
	p.PublishedAt = &t
	p.UpdatedAt = t
	p.Version++
	return nil
}

// Archive es irreversible. Idempotente.
func (p *Product) Archive(now time.Time) error {
	if p.Status == StatusArchived {
		return nil
	}
	t := now.UTC()
	p.Status = StatusArchived
	p.ArchivedAt = &t
	p.UpdatedAt = t
	p.Version++
	return nil
}

// CanSell retorna true si el producto puede venderse hoy.
func (p *Product) CanSell() bool {
	return p.Status == StatusPublished
}

func validateSaleModeProvider(sm SaleMode, pid *shared.ProviderID) error {
	switch sm {
	case SaleOwnInventory:
		if pid != nil {
			return fmt.Errorf("%w: own_inventory no lleva provider", ErrInvalidSaleMode)
		}
	case SaleDropshipping:
		if pid == nil {
			return ErrProviderReq
		}
	default:
		return fmt.Errorf("%w: %q", ErrInvalidSaleMode, sm)
	}
	return nil
}

func isValidSlug(s string) bool {
	if s == "" || len(s) > MaxSlugLen {
		return false
	}
	for i, r := range s {
		switch {
		case r >= 'a' && r <= 'z':
		case r >= '0' && r <= '9':
		case r == '-':
		default:
			return false
		}
		if i == 0 && r == '-' {
			return false
		}
	}
	return s[len(s)-1] != '-'
}

func isValidSKU(s string) bool {
	if s == "" || len(s) > MaxSKULen {
		return false
	}
	for _, r := range s {
		switch {
		case r >= 'A' && r <= 'Z':
		case r >= 'a' && r <= 'z':
		case r >= '0' && r <= '9':
		case r == '-', r == '_':
		default:
			return false
		}
	}
	return true
}
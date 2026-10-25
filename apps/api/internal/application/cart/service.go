package cart

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	domain "uvico-api/internal/domain/cart"
	"uvico-api/internal/domain/shared"
)

// ProductReader es la interfaz mínima del catálogo para validar
// precio actual al hacer checkout. La implementa CatalogRepo.
type ProductReader interface {
	GetPublishedPrice(ctx context.Context, tenantID shared.TenantID, productID shared.ProductID) (sku, name string, priceCents int64, currency string, err error)
}

type Service struct {
	repo domain.Repository
	cat  ProductReader
	now  func() time.Time
}

func NewService(repo domain.Repository, cat ProductReader) *Service {
	return &Service{repo: repo, cat: cat, now: time.Now}
}

type AddInput struct {
	TenantID  shared.TenantID
	SessionID string
	ProductID shared.ProductID
	Quantity  int32
}

func (s *Service) AddItem(ctx context.Context, in AddInput) (*domain.Cart, error) {
	if in.SessionID == "" {
		return nil, domain.ErrEmptySessionID
	}
	if err := domain.ValidateQuantity(in.Quantity); err != nil {
		return nil, err
	}

	sku, name, price, currency, err := s.cat.GetPublishedPrice(ctx, in.TenantID, in.ProductID)
	if err != nil {
		return nil, err
	}

	return s.repo.AddItem(ctx, in.TenantID, in.SessionID, domain.Item{
		ProductID:       in.ProductID,
		SKU:             sku,
		Name:            name,
		Quantity:        in.Quantity,
		PriceCentsAtAdd: price,
		Currency:        currency,
		AddedAt:         s.now().UTC(),
	})
}

func (s *Service) SetQuantity(ctx context.Context, tenantID shared.TenantID, sessionID string, productID shared.ProductID, qty int32) (*domain.Cart, error) {
	if sessionID == "" {
		return nil, domain.ErrEmptySessionID
	}
	if err := domain.ValidateQuantity(qty); err != nil {
		return nil, err
	}
	return s.repo.SetQuantity(ctx, tenantID, sessionID, productID, qty)
}

func (s *Service) RemoveItem(ctx context.Context, tenantID shared.TenantID, sessionID string, productID shared.ProductID) (*domain.Cart, error) {
	if sessionID == "" {
		return nil, domain.ErrEmptySessionID
	}
	return s.repo.RemoveItem(ctx, tenantID, sessionID, productID)
}

func (s *Service) Get(ctx context.Context, tenantID shared.TenantID, sessionID string) (*domain.Cart, error) {
	if sessionID == "" {
		return nil, domain.ErrEmptySessionID
	}
	return s.repo.Get(ctx, tenantID, sessionID)
}

func (s *Service) Clear(ctx context.Context, tenantID shared.TenantID, sessionID string) error {
	if sessionID == "" {
		return domain.ErrEmptySessionID
	}
	return s.repo.Clear(ctx, tenantID, sessionID)
}

// ValidateForCheckout carga el cart, valida precios actuales y retorna
// los items como LineItems del dominio order. Falla con ErrPriceChanged
// si algún precio cambió (el frontend debe reconfirmar).
func (s *Service) ValidateForCheckout(
	ctx context.Context,
	tenantID shared.TenantID,
	sessionID string,
) ([]CartLine, error) {
	c, err := s.repo.Get(ctx, tenantID, sessionID)
	if err != nil {
		return nil, err
	}
	if len(c.Items) == 0 {
		return nil, domain.ErrEmptyCart
	}

	currentPrices := make(map[shared.ProductID]int64, len(c.Items))
	for _, it := range c.Items {
		_, _, price, _, err := s.cat.GetPublishedPrice(ctx, tenantID, it.ProductID)
		if err != nil {
			return nil, fmt.Errorf("cart: producto %s: %w", it.ProductID, err)
		}
		currentPrices[it.ProductID] = price
	}

	changes := c.ValidatePrices(currentPrices)
	if len(changes) > 0 {
		return nil, &PriceChangedError{Changes: changes}
	}

	out := make([]CartLine, 0, len(c.Items))
	for _, it := range c.Items {
		out = append(out, CartLine{
			ProductID:  it.ProductID,
			SKU:        it.SKU,
			Name:       it.Name,
			Quantity:   it.Quantity,
			UnitPrice:  shared.Money{AmountCents: currentPrices[it.ProductID], Currency: shared.Currency(it.Currency)},
			Currency:   it.Currency,
		})
	}
	return out, nil
}

type CartLine struct {
	ProductID shared.ProductID
	SKU       string
	Name      string
	Quantity  int32
	UnitPrice shared.Money
	Currency  string
}

type PriceChangedError struct {
	Changes []domain.PriceChange
}

func (e *PriceChangedError) Error() string {
	return "cart: " + fmt.Sprint(len(e.Changes)) + " precios cambiaron"
}

func (e *PriceChangedError) Is(target error) bool {
	return target == domain.ErrPriceChanged
}

// GenerateSessionID genera un session ID si el cliente no lo envía.
func GenerateSessionID() string {
	return uuid.Must(uuid.NewV7()).String()
}

var _ = errors.New
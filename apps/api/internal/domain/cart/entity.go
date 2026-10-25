package cart

import (
	"errors"
	"fmt"
	"time"

	"uvico-api/internal/domain/shared"
)

var (
	ErrEmptySessionID  = errors.New("cart: session_id requerido")
	ErrItemNotFound    = errors.New("cart: item no encontrado")
	ErrInvalidQuantity = errors.New("cart: cantidad inválida")
	ErrEmptyCart       = errors.New("cart: carrito vacío")
	ErrPriceChanged    = errors.New("cart: precio cambió desde que se agregó al carrito")
	ErrProductNotFound = errors.New("cart: producto no encontrado")
	ErrProductArchived = errors.New("cart: producto archivado")
)

const MaxItemsPerCart = 100

type Item struct {
	ProductID         shared.ProductID `json:"product_id"`
	SKU               string           `json:"sku"`
	Name              string           `json:"name"`
	Quantity          int32            `json:"quantity"`
	PriceCentsAtAdd   int64            `json:"price_cents_at_add"`
	PriceCentsCurrent int64            `json:"price_cents_current,omitempty"`
	Currency          string           `json:"currency"`
	AddedAt           time.Time        `json:"added_at"`
}

type Cart struct {
	TenantID  shared.TenantID `json:"tenant_id"`
	SessionID string          `json:"session_id"`
	Items     []Item          `json:"items"`
	Subtotal  int64           `json:"subtotal_cents"`
	Currency  string          `json:"currency"`
	UpdatedAt time.Time       `json:"updated_at"`
}

type PriceChange struct {
	ProductID       shared.ProductID `json:"product_id"`
	Name            string           `json:"name"`
	PriceCentsAtAdd int64            `json:"price_cents_at_add"`
	PriceCentsNow   int64            `json:"price_cents_now"`
}

func (c *Cart) RecalculateSubtotal() {
	var total int64
	for _, it := range c.Items {
		total += it.PriceCentsAtAdd * int64(it.Quantity)
	}
	c.Subtotal = total
}

// ValidatePrices compara el precio actual contra el precio al agregar.
// Retorna la lista de cambios; vacío = todo bien.
func (c *Cart) ValidatePrices(current map[shared.ProductID]int64) []PriceChange {
	var changes []PriceChange
	for _, it := range c.Items {
		now, ok := current[it.ProductID]
		if !ok {
			continue
		}
		if now != it.PriceCentsAtAdd {
			changes = append(changes, PriceChange{
				ProductID:       it.ProductID,
				Name:            it.Name,
				PriceCentsAtAdd: it.PriceCentsAtAdd,
				PriceCentsNow:   now,
			})
		}
	}
	return changes
}

func ValidateQuantity(q int32) error {
	if q <= 0 || q > 1000 {
		return fmt.Errorf("%w: %d", ErrInvalidQuantity, q)
	}
	return nil
}
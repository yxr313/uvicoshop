package inventory

import (
	"fmt"
	"time"

	"uvico-api/internal/domain/shared"
)

type Reason string

const (
	ReasonReserved   Reason = "reserved"
	ReasonCommitted  Reason = "committed"
	ReasonReleased   Reason = "released"
	ReasonAdjustment Reason = "adjustment"
)

type Item struct {
	TenantID  shared.TenantID
	ProductID shared.ProductID
	OnHand    int64
	Reserved  int64
	Version   int64
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (i *Item) Available() int64 { return i.OnHand - i.Reserved }

// Reserve valida invariantes y devuelve el delta de movimiento.
func (i *Item) Reserve(qty int64, now time.Time) (Reason, int64, error) {
	if qty <= 0 {
		return "", 0, ErrInvalidQuantity
	}
	if i.Available() < qty {
		return "", 0, ErrInsufficientStock
	}
	i.Reserved += qty
	i.Version++
	i.UpdatedAt = now.UTC()
	return ReasonReserved, -qty, nil
}

func (i *Item) Adjust(delta int64, now time.Time) (Reason, int64, error) {
	if delta == 0 {
		return "", 0, ErrInvalidQuantity
	}
	next := i.OnHand + delta
	if next < 0 {
		return "", 0, fmt.Errorf("%w: on_hand=%d delta=%d", ErrOnHandNegative, i.OnHand, delta)
	}
	if next < i.Reserved {
		return "", 0, fmt.Errorf("%w: on_hand=%d < reserved=%d", ErrInsufficientStock, next, i.Reserved)
	}
	i.OnHand = next
	i.Version++
	i.UpdatedAt = now.UTC()
	return ReasonAdjustment, delta, nil
}
package shared

import (
	"fmt"
	"math"
)

type Currency string

const (
	CurrencyUSD Currency = "USD"
	CurrencyVES Currency = "VES"
)

type Money struct {
	AmountCents int64
	Currency    Currency
}

func USD(cents int64) Money { return Money{AmountCents: cents, Currency: CurrencyUSD} }
func VES(cents int64) Money { return Money{AmountCents: cents, Currency: CurrencyVES} }

func (m Money) Validate() error {
	switch m.Currency {
	case CurrencyUSD, CurrencyVES:
	default:
		return fmt.Errorf("%w: %q", ErrInvalidCurrency, m.Currency)
	}
	if m.AmountCents < 0 {
		return fmt.Errorf("%w: monto negativo %d", ErrInvariant, m.AmountCents)
	}
	return nil
}

func (m Money) IsZero() bool     { return m.AmountCents == 0 }
func (m Money) IsNegative() bool { return m.AmountCents < 0 }

func (m Money) Add(other Money) (Money, error) {
	if m.Currency != other.Currency {
		return Money{}, ErrCurrencyMismatch
	}
	sum := m.AmountCents + other.AmountCents
	if (m.AmountCents > 0 && other.AmountCents > 0 && sum < m.AmountCents) ||
		(m.AmountCents < 0 && other.AmountCents < 0 && sum > m.AmountCents) {
		return Money{}, ErrOverflow
	}
	return Money{AmountCents: sum, Currency: m.Currency}, nil
}

func (m Money) Subtract(other Money) (Money, error) {
	if m.Currency != other.Currency {
		return Money{}, ErrCurrencyMismatch
	}
	result := m.AmountCents - other.AmountCents
	if result < 0 {
		return Money{}, fmt.Errorf("%w: resultado negativo %d", ErrInvariant, result)
	}
	return Money{AmountCents: result, Currency: m.Currency}, nil
}

func (m Money) Multiply(qty int32) (Money, error) {
	if qty < 0 {
		return Money{}, ErrNegativeQuantity
	}
	if qty == 0 || m.AmountCents == 0 {
		return Money{AmountCents: 0, Currency: m.Currency}, nil
	}
	if m.AmountCents > math.MaxInt64/int64(qty) {
		return Money{}, ErrOverflow
	}
	return Money{AmountCents: m.AmountCents * int64(qty), Currency: m.Currency}, nil
}

// SplitBps distribuye total en N partes según basis points.
// Invariante garantizada: Sum(parts) + merchant == total (exacto, sin residuo flotante).
// El residuo del redondeo (hasta N-1 centavos) va al merchant (customer-friendly).
func SplitBps(total Money, bps ...int64) ([]Money, Money, error) {
	if len(bps) == 0 {
		return nil, total, nil
	}
	var sumBps int64
	for _, b := range bps {
		if b < 0 || b > 10000 {
			return nil, Money{}, fmt.Errorf("%w: bps %d", ErrInvalidSplit, b)
		}
		sumBps += b
	}
	if sumBps > 10000 {
		return nil, Money{}, fmt.Errorf("%w: suma %d > 10000", ErrInvalidSplit, sumBps)
	}

	parts := make([]Money, len(bps))
	var allocated int64
	for i, b := range bps {
		if b > 0 && total.AmountCents > math.MaxInt64/b {
			return nil, Money{}, ErrOverflow
		}
		amount := (total.AmountCents * b) / 10000
		parts[i] = Money{AmountCents: amount, Currency: total.Currency}
		allocated += amount
	}
	return parts, Money{AmountCents: total.AmountCents - allocated, Currency: total.Currency}, nil
}
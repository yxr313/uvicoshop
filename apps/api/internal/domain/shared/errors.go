package shared

import "errors"

var (
	ErrInvalidCurrency  = errors.New("shared: moneda inválida")
	ErrCurrencyMismatch = errors.New("shared: monedas distintas")
	ErrOverflow         = errors.New("shared: overflow aritmético")
	ErrInvalidSplit     = errors.New("shared: split inválido")
	ErrNegativeQuantity = errors.New("shared: cantidad negativa")
	ErrInvariant        = errors.New("shared: invariante violada")
)
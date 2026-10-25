package inventory

import "errors"

var (
	ErrInsufficientStock = errors.New("inventory: stock insuficiente")
	ErrInvalidQuantity   = errors.New("inventory: cantidad inválida")
	ErrNotFound          = errors.New("inventory: item no encontrado")
	ErrNotOwnInventory   = errors.New("inventory: solo aplica a productos de inventario propio")
	ErrOnHandNegative    = errors.New("inventory: on_hand no puede ser negativo")
	ErrReservedNegative  = errors.New("inventory: reserved no puede ser negativo")
	ErrVersionConflict   = errors.New("inventory: version conflict")
	ErrInvalidReason     = errors.New("inventory: razón inválida")
)
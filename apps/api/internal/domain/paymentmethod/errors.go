package paymentmethod

import "errors"

var (
	ErrNotFound      = errors.New("paymentmethod: método no encontrado")
	ErrInvalidMethod = errors.New("paymentmethod: método inválido")
	ErrInvalidConfig = errors.New("paymentmethod: configuración inválida")
	ErrDisabled      = errors.New("paymentmethod: método no habilitado para esta tienda")
	ErrNoMethods     = errors.New("paymentmethod: la tienda no acepta pagos configurados")
)
package order

import "errors"

var (
	ErrInvalidTenant         = errors.New("order: tenant requerido")
	ErrInvalidCustomer       = errors.New("order: customer requerido")
	ErrInvalidItems          = errors.New("order: al menos un item requerido")
	ErrMixedCurrencies       = errors.New("order: monedas mixtas en items")
	ErrMixedSaleModes        = errors.New("order: sale_modes mixtas en items")
	ErrInvalidTotal          = errors.New("order: total inválido")
	ErrInvalidTransition     = errors.New("order: transición de estado inválida")
	ErrMissingReference      = errors.New("order: referencia de pago requerida")
	ErrMissingPaymentMethod  = errors.New("order: método de pago requerido")
	ErrMissingAccessToken    = errors.New("order: access token requerido")
	ErrMissingMerchantID     = errors.New("order: merchant user id requerido")
	ErrMissingRejectionReason = errors.New("order: razón de rechazo requerida")
	ErrVersionConflict       = errors.New("order: version conflict")
	ErrProviderRequired      = errors.New("order: dropshipping requiere provider")
	ErrNotFound              = errors.New("order: orden no encontrada")
	ErrSplitInvariant        = errors.New("order: split no suma al total")
	ErrInvalidQuantity       = errors.New("order: cantidad inválida")
	ErrInvalidStatus         = errors.New("order: status inválido")
	ErrDropshippingDisabled  = errors.New("order: dropshipping deshabilitado en MVP")
	ErrAlreadySubmitted      = errors.New("order: la orden ya tiene un pago registrado")
)
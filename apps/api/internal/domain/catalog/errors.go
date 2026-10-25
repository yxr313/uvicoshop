package catalog

import "errors"

var (
	ErrInvalidTenant   = errors.New("catalog: tenant requerido")
	ErrInvalidName     = errors.New("catalog: nombre inválido (1..120)")
	ErrInvalidSlug     = errors.New("catalog: slug inválido")
	ErrInvalidSKU      = errors.New("catalog: sku inválido (1..64, alfanumérico, - o _)")
	ErrInvalidPrice    = errors.New("catalog: precio inválido")
	ErrInvalidSaleMode = errors.New("catalog: sale_mode inválido")
	ErrInvalidStatus   = errors.New("catalog: status inválido")
	ErrSlugTaken       = errors.New("catalog: slug ya existe para este tenant")
	ErrSKUTaken        = errors.New("catalog: sku ya existe para este tenant")
	ErrNotFound        = errors.New("catalog: producto no encontrado")
	ErrVersionConflict = errors.New("catalog: version conflict")
	ErrArchived        = errors.New("catalog: producto archivado")
	ErrProviderReq     = errors.New("catalog: dropshipping requiere provider_id")
	ErrProviderNotImpl = errors.New("catalog: provider_id no soportado aún")
)
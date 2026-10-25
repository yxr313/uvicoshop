package image

import "errors"

var (
	ErrInvalidMIMEType  = errors.New("image: tipo MIME no permitido")
	ErrFileTooLarge     = errors.New("image: archivo excede el tamaño máximo")
	ErrInvalidPath      = errors.New("image: path inválido")
	ErrNotFound         = errors.New("image: imagen no encontrada")
	ErrProductNotFound  = errors.New("image: producto no encontrado")
	ErrMaxImages        = errors.New("image: se alcanzó el máximo de imágenes por producto")
	ErrPositionConflict = errors.New("image: posición inválida")
)
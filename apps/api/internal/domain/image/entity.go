package image

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"uvico-api/internal/domain/shared"
)

const (
	MaxImagesPerProduct = 10
	MaxAltTextLen       = 200
	MinPosition         = 0
	MaxPosition         = MaxImagesPerProduct - 1
)

var (
	storagePathRegex = regexp.MustCompile(`^products/[0-9a-f-]{36}/[0-9a-f-]{36}/[0-9a-f-]{36}\.(jpg|jpeg|png|webp|avif)$`)
	extByMIME        = map[string]string{
		"image/jpeg": "jpg",
		"image/png":  "png",
		"image/webp": "webp",
		"image/avif": "avif",
	}
)

type Image struct {
	ID          string
	TenantID    shared.TenantID
	ProductID   shared.ProductID
	StoragePath string
	AltText     string
	Position    int
	Width       *int
	Height      *int
	SizeBytes   *int
	CreatedAt   time.Time
}

// ExtensionForMIME retorna la extensión para un MIME permitido.
func ExtensionForMIME(mime string) (string, bool) {
	ext, ok := extByMIME[strings.ToLower(strings.TrimSpace(mime))]
	return ext, ok
}

// BuildStoragePath construye el path canónico dentro del bucket.
// Formato: products/{tenant_id}/{product_id}/{uuid}.{ext}
// El bucket se llama "products", así que el path empieza con el prefijo
// del tenant para que las URLs públicas sean legibles.
func BuildStoragePath(tenantID shared.TenantID, productID shared.ProductID, ext string) string {
	return fmt.Sprintf("products/%s/%s/%s.%s",
		tenantID.String(), productID.String(), newUUID(), ext)
}

// ValidateStoragePath verifica que un path cumple el formato esperado.
func ValidateStoragePath(tenantID shared.TenantID, productID shared.ProductID, p string) error {
	if !storagePathRegex.MatchString(p) {
		return fmt.Errorf("%w: formato inesperado %q", ErrInvalidPath, p)
	}
	parts := strings.Split(p, "/")
	if len(parts) != 4 {
		return ErrInvalidPath
	}
	if parts[0] != "products" {
		return ErrInvalidPath
	}
	if parts[1] != tenantID.String() {
		return fmt.Errorf("%w: tenant mismatch (esperado %s, recibido %s)",
			ErrInvalidPath, tenantID.String(), parts[1])
	}
	if parts[2] != productID.String() {
		return fmt.Errorf("%w: product mismatch (esperado %s, recibido %s)",
			ErrInvalidPath, productID.String(), parts[2])
	}
	return nil
}

func ValidateAltText(s string) (string, error) {
	s = strings.TrimSpace(s)
	if len(s) > MaxAltTextLen {
		return "", fmt.Errorf("image: alt_text excede %d caracteres", MaxAltTextLen)
	}
	return s, nil
}

func ValidatePosition(p int) error {
	if p < MinPosition || p > MaxPosition {
		return fmt.Errorf("%w: %d fuera de [%d,%d]", ErrPositionConflict, p, MinPosition, MaxPosition)
	}
	return nil
}

func newUUID() string {
	return mustUUID()
}

var _ = errors.New
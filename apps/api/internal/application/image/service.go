package image

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	domain "uvico-api/internal/domain/image"
	"uvico-api/internal/domain/shared"
)

// Storage es el puerto hacia Supabase Storage. La implementación real vive
// en internal/adapter/outbound/storage.
type Storage interface {
	CreateSignedUploadURL(ctx context.Context, path string, ttl time.Duration) (string, time.Time, error)
	Delete(ctx context.Context, path string) error
}

// ProductReader valida que el producto exista y pertenece al tenant.
type ProductReader interface {
	GetPublishedPrice(ctx context.Context, tenantID shared.TenantID, productID shared.ProductID) (string, string, int64, string, error)
	GetAnyStatusPrice(ctx context.Context, tenantID shared.TenantID, productID shared.ProductID) (string, string, int64, string, string, error)
}

type Options struct {
	MaxFileSizeBytes int64
	AllowedMIMETypes []string
	PresignedURLTTL  time.Duration
}

type Service struct {
	repo    domain.Repository
	cat     ProductReader
	storage Storage
	opts    Options
	now     func() time.Time
}

func NewService(repo domain.Repository, cat ProductReader, storage Storage, opts Options) *Service {
	if opts.PresignedURLTTL <= 0 {
		opts.PresignedURLTTL = 5 * time.Minute
	}
	return &Service{repo: repo, cat: cat, storage: storage, opts: opts, now: time.Now}
}

type PresignInput struct {
	TenantID    shared.TenantID
	ProductID   shared.ProductID
	ContentType string
	FileSize    int64
}

type PresignResult struct {
	UploadURL   string    `json:"upload_url"`
	StoragePath string    `json:"storage_path"`
	ExpiresAt   time.Time `json:"expires_at"`
	PublicURL   string    `json:"public_url"`
}

func (s *Service) CreatePresignedUpload(ctx context.Context, in PresignInput) (*PresignResult, error) {
	if in.FileSize <= 0 {
		return nil, fmt.Errorf("%w: file_size requerido y positivo", domain.ErrInvalidPath)
	}
	if in.FileSize > s.opts.MaxFileSizeBytes {
		return nil, fmt.Errorf("%w: max %d bytes", domain.ErrFileTooLarge, s.opts.MaxFileSizeBytes)
	}
	ct := strings.ToLower(strings.TrimSpace(in.ContentType))
	if !s.isAllowedMIME(ct) {
		return nil, fmt.Errorf("%w: %q no permitido", domain.ErrInvalidMIMEType, ct)
	}
	ext, ok := domain.ExtensionForMIME(ct)
	if !ok {
		return nil, fmt.Errorf("%w: %q sin extensión conocida", domain.ErrInvalidMIMEType, ct)
	}

	// Verificar que el producto exista en el tenant. Usa GetAnyStatusPrice
	// para aceptar productos en draft también (el merchant sube fotos antes
	// de publicar).
	if _, _, _, _, _, err := s.cat.GetAnyStatusPrice(ctx, in.TenantID, in.ProductID); err != nil {
		return nil, fmt.Errorf("image: %w", err)
	}

	count, err := s.repo.CountByProduct(ctx, in.TenantID, in.ProductID)
	if err != nil {
		return nil, err
	}
	if count >= domain.MaxImagesPerProduct {
		return nil, fmt.Errorf("%w: max %d por producto", domain.ErrMaxImages, domain.MaxImagesPerProduct)
	}

	storagePath := domain.BuildStoragePath(in.TenantID, in.ProductID, ext)

	uploadURL, expiresAt, err := s.storage.CreateSignedUploadURL(ctx, storagePath, s.opts.PresignedURLTTL)
	if err != nil {
		return nil, fmt.Errorf("image: presign: %w", err)
	}

	return &PresignResult{
		UploadURL:   uploadURL,
		StoragePath: storagePath,
		ExpiresAt:   expiresAt,
	}, nil
}

type AttachInput struct {
	TenantID    shared.TenantID
	ProductID   shared.ProductID
	StoragePath string
	AltText     string
}

func (s *Service) AttachImage(ctx context.Context, in AttachInput) (*domain.Image, error) {
	if err := domain.ValidateStoragePath(in.TenantID, in.ProductID, in.StoragePath); err != nil {
		return nil, err
	}
	alt, err := domain.ValidateAltText(in.AltText)
	if err != nil {
		return nil, err
	}
	count, err := s.repo.CountByProduct(ctx, in.TenantID, in.ProductID)
	if err != nil {
		return nil, err
	}
	if count >= domain.MaxImagesPerProduct {
		return nil, fmt.Errorf("%w: max %d por producto", domain.ErrMaxImages, domain.MaxImagesPerProduct)
	}

	img := &domain.Image{
		ID:          mustID(),
		TenantID:    in.TenantID,
		ProductID:   in.ProductID,
		StoragePath: in.StoragePath,
		AltText:     alt,
		Position:    count,
		CreatedAt:   s.now().UTC(),
	}
	if err := s.repo.Insert(ctx, img); err != nil {
		return nil, err
	}
	return img, nil
}

func (s *Service) List(ctx context.Context, tenantID shared.TenantID, productID shared.ProductID) ([]*domain.Image, error) {
	return s.repo.ListByProduct(ctx, tenantID, productID)
}

func (s *Service) Delete(ctx context.Context, tenantID shared.TenantID, productID shared.ProductID, imageID string) error {
	img, err := s.repo.GetByID(ctx, tenantID, productID, imageID)
	if err != nil {
		return err
	}
	// Borrar del storage primero; si falla, no borramos la fila para que
	// se pueda reintentar. Idempotente: si el objeto ya no existe, nil.
	if err := s.storage.Delete(ctx, img.StoragePath); err != nil {
		return fmt.Errorf("image: storage delete: %w", err)
	}
	return s.repo.Delete(ctx, tenantID, productID, imageID)
}

func (s *Service) Reorder(ctx context.Context, tenantID shared.TenantID, productID shared.ProductID, order []string) error {
	if len(order) == 0 {
		return errors.New("image: order vacío")
	}
	if len(order) > domain.MaxImagesPerProduct {
		return fmt.Errorf("%w: max %d", domain.ErrPositionConflict, domain.MaxImagesPerProduct)
	}
	return s.repo.Reorder(ctx, tenantID, productID, order)
}

func (s *Service) isAllowedMIME(ct string) bool {
	for _, allowed := range s.opts.AllowedMIMETypes {
		if strings.EqualFold(allowed, ct) {
			return true
		}
	}
	return false
}

func mustID() string {
	return uuidV7()
}
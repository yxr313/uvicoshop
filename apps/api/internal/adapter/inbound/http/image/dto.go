package image

type presignRequest struct {
	ContentType string `json:"content_type" validate:"required,oneof=image/jpeg image/png image/webp image/avif"`
	FileSize    int64  `json:"file_size" validate:"required,gt=0"`
}

type presignResponse struct {
	UploadURL   string `json:"upload_url"`
	StoragePath string `json:"storage_path"`
	ExpiresAt   string `json:"expires_at"`
}

type attachRequest struct {
	StoragePath string `json:"storage_path" validate:"required,min=1,max=512"`
	AltText     string `json:"alt_text" validate:"max=200"`
}

type imageResponse struct {
	ID        string `json:"id"`
	PublicURL string `json:"public_url"`
	AltText   string `json:"alt_text"`
	Position  int    `json:"position"`
	CreatedAt string `json:"created_at"`
}

type reorderRequest struct {
	Order []string `json:"order" validate:"required,min=1,max=10,dive,uuid_any"`
}

type errorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}
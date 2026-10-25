package catalog

type createProductRequest struct {
	SaleMode     string  `json:"sale_mode" validate:"required,oneof=own_inventory dropshipping_verified"`
	ProviderID   *string `json:"provider_id,omitempty" validate:"omitempty,uuid4"`
	SKU          string  `json:"sku" validate:"required,min=1,max=64"`
	Slug         string  `json:"slug" validate:"required,min=1,max=80"`
	Name         string  `json:"name" validate:"required,min=1,max=120"`
	Description  string  `json:"description" validate:"max=5000"`
	PriceCents   int64   `json:"price_cents" validate:"required,gt=0"`
	InitialStock int64   `json:"initial_stock" validate:"gte=0"`
}

type updateProductRequest struct {
	Name        string `json:"name" validate:"required,min=1,max=120"`
	Description string `json:"description" validate:"max=5000"`
	Slug        string `json:"slug" validate:"required,min=1,max=80"`
	PriceCents  int64  `json:"price_cents" validate:"required,gt=0"`
}

type productResponse struct {
	ID          string  `json:"id"`
	TenantID    string  `json:"tenant_id"`
	ProviderID  *string `json:"provider_id,omitempty"`
	SaleMode    string  `json:"sale_mode"`
	SKU         string  `json:"sku"`
	Slug        string  `json:"slug"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	PriceCents  int64   `json:"price_cents"`
	Currency    string  `json:"currency"`
	Status      string  `json:"status"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
	Version     int64   `json:"version"`
}

type errorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}
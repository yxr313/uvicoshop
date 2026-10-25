package storefront

type checkoutRequest struct {
	// Datos del comprador (guest checkout, sin cuenta Uvico).
	CustomerEmail string `json:"customer_email" validate:"required,email,max=254"`
	CustomerName  string `json:"customer_name" validate:"required,min=1,max=120"`
	CustomerPhone string `json:"customer_phone" validate:"omitempty,max=30"`
	MarketingOptIn bool  `json:"marketing_opt_in"`

	// Atribución (opcional).
	Attribution attributionDTO `json:"attribution"`

	// Pago y modo de venta.
	PaymentMethod string `json:"payment_method" validate:"required"`
	SaleMode      string `json:"sale_mode" validate:"required,oneof=own_inventory dropshipping_verified"`
	ProviderID    *string `json:"provider_id,omitempty" validate:"omitempty,uuid_any"`

	// Items: cart_id o items directos.
	CartID string            `json:"cart_id,omitempty"`
	Items  []checkoutLineDTO `json:"items,omitempty" validate:"omitempty,min=1,max=100,dive"`
}

type attributionDTO struct {
	SessionID   string `json:"session_id,omitempty"`
	UTMSource   string `json:"utm_source,omitempty"`
	UTMMedium   string `json:"utm_medium,omitempty"`
	UTMCampaign string `json:"utm_campaign,omitempty"`
	UTMTerm     string `json:"utm_term,omitempty"`
	UTMContent  string `json:"utm_content,omitempty"`
	FBCLID      string `json:"fbclid,omitempty"`
	GCLID       string `json:"gclid,omitempty"`
	TTCLID      string `json:"ttclid,omitempty"`
	TWCLID      string `json:"twclid,omitempty"`
	Referrer    string `json:"referrer,omitempty"`
	LandingPage string `json:"landing_page,omitempty"`
	Timezone    string `json:"timezone,omitempty"`
}

type checkoutLineDTO struct {
	ProductID  string `json:"product_id" validate:"required,uuid_any"`
	SKU        string `json:"sku" validate:"required,max=64"`
	Name       string `json:"name" validate:"required,max=120"`
	Quantity   int32  `json:"quantity" validate:"required,gt=0,lte=1000"`
	PriceCents int64  `json:"price_cents" validate:"required,gt=0"`
}

type paymentInstructionsDTO struct {
	Method    string         `json:"method"`
	Label     string         `json:"label"`
	Config    map[string]any `json:"config"`
	Reference string         `json:"reference"`
}

type checkoutResponse struct {
	OrderID             string                `json:"order_id"`
	AccessToken         string                `json:"access_token"`
	Status              string                `json:"status"`
	TotalCents          int64                 `json:"total_cents"`
	Currency            string                `json:"currency"`
	PaymentInstructions paymentInstructionsDTO `json:"payment_instructions"`
}

type paymentMethodPublicDTO struct {
	Method   string         `json:"method"`
	Label    string         `json:"label"`
	Position int            `json:"position"`
	Config   map[string]any `json:"config"`
}

type submitPaymentRequest struct {
	Reference string `json:"reference" validate:"required,min=3,max=200"`
	ProofURL  string `json:"proof_url" validate:"omitempty,max=512"`
}

type submitPaymentResponse struct {
	OrderID           string `json:"order_id"`
	Status            string `json:"status"`
	Reference         string `json:"reference"`
	CustomerSubmittedAt string `json:"customer_submitted_at"`
}

type orderPublicResponse struct {
	OrderID             string                `json:"order_id"`
	Status              string                `json:"status"`
	TotalCents          int64                 `json:"total_cents"`
	Currency            string                `json:"currency"`
	Items               []orderItemPublicDTO  `json:"items"`
	PaymentInstructions paymentInstructionsDTO `json:"payment_instructions"`
	CustomerSubmittedAt string                `json:"customer_submitted_at,omitempty"`
	MerchantRejectionReason string            `json:"merchant_rejection_reason,omitempty"`
}

type orderItemPublicDTO struct {
	ProductID string `json:"product_id"`
	SKU       string `json:"sku"`
	Name      string `json:"name"`
	Quantity  int32  `json:"quantity"`
	UnitPrice int64  `json:"unit_price_cents"`
	Total     int64  `json:"total_price_cents"`
}
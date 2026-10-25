package storefront

import "time"

type ProductImageRef struct {
	Path string `json:"path"`
	Alt  string `json:"alt"`
}

type ProductSummary struct {
	ID          string            `json:"id"`
	SKU         string            `json:"sku"`
	Slug        string            `json:"slug"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	PriceCents  int64             `json:"price_cents"`
	Currency    string            `json:"currency"`
	SaleMode    string            `json:"sale_mode"`
	Available   *int64            `json:"available"` // nil en dropshipping
	Images      []ProductImageRef `json:"images"`
}

type ProductDetail struct {
	ProductSummary
	UpdatedAt time.Time `json:"updated_at"`
}

type StoreConfig struct {
	TenantID        string         `json:"tenant_id"`
	Slug            string         `json:"slug"`
	Name            string         `json:"name"`
	LogoURL         *string        `json:"logo_url,omitempty"`
	PrimaryColor    *string        `json:"primary_color,omitempty"`
	ContactEmail    *string        `json:"contact_email,omitempty"`
	ContactPhone    *string        `json:"contact_phone,omitempty"`
	WhatsApp        *string        `json:"whatsapp,omitempty"`
	CurrencyDisplay string         `json:"currency_display"`
	SocialInstagram *string        `json:"social_instagram,omitempty"`
	SocialTikTok    *string        `json:"social_tiktok,omitempty"`
	MetaTitle       *string        `json:"meta_title,omitempty"`
	MetaDescription *string        `json:"meta_description,omitempty"`
	PuckData        map[string]any `json:"puck_data"`
}

type BCVRate struct {
	RateMicros int64     `json:"rate_micros"`
	Rate       string    `json:"rate"`
	UpdatedAt  time.Time `json:"updated_at"`
}
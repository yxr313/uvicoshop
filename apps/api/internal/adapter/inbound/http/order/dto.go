package order

type orderSummaryResponse struct {
	ID            string `json:"id"`
	CustomerID    string `json:"customer_id"`
	Status        string `json:"status"`
	PaymentMethod string `json:"payment_method,omitempty"`
	TotalCents    int64  `json:"total_cents"`
	Currency      string `json:"currency"`
	CreatedAt     string `json:"created_at"`
}

type orderDetailResponse struct {
	orderSummaryResponse
	Items                    []orderItemResponse `json:"items"`
	PaymentRef               string              `json:"payment_ref,omitempty"`
	CustomerPaymentReference string              `json:"customer_payment_reference,omitempty"`
	CustomerPaymentProofURL  string              `json:"customer_payment_proof_url,omitempty"`
	CustomerSubmittedAt      string              `json:"customer_submitted_at,omitempty"`
	MerchantConfirmedAt      string              `json:"merchant_confirmed_at,omitempty"`
	MerchantConfirmedBy      string              `json:"merchant_confirmed_by,omitempty"`
	MerchantRejectionReason  string              `json:"merchant_rejection_reason,omitempty"`
	UvicoCommission          int64               `json:"uvico_commission_cents"`
	MerchantNet              int64               `json:"merchant_net_cents"`
}

type orderItemResponse struct {
	ProductID string `json:"product_id"`
	SKU       string `json:"sku"`
	Name      string `json:"name"`
	Quantity  int32  `json:"quantity"`
	UnitPrice int64  `json:"unit_price_cents"`
	Total     int64  `json:"total_price_cents"`
}

type rejectRequest struct {
	Reason string `json:"reason" validate:"required,min=3,max=500"`
}

type errorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}
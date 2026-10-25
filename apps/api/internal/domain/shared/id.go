package shared

import (
	"fmt"

	"github.com/google/uuid"
)

type (
	TenantID   string
	ProductID  string
	ProviderID string
	OrderID    string
)

type CustomerID string

func NewCustomerID() CustomerID { return CustomerID(uuid.Must(uuid.NewV7()).String()) }

func (id CustomerID) String() string { return string(id) }
func (id CustomerID) IsZero() bool   { return id == "" }

func ParseCustomerID(s string) (CustomerID, error) {
	if _, err := uuid.Parse(s); err != nil {
		return "", fmt.Errorf("shared: customer id inválido: %w", err)
	}
	return CustomerID(s), nil
}

func NewTenantID() TenantID     { return TenantID(uuid.Must(uuid.NewV7()).String()) }
func NewProductID() ProductID   { return ProductID(uuid.Must(uuid.NewV7()).String()) }
func NewProviderID() ProviderID { return ProviderID(uuid.Must(uuid.NewV7()).String()) }
func NewOrderID() OrderID       { return OrderID(uuid.Must(uuid.NewV7()).String()) }

func (id TenantID) String() string   { return string(id) }
func (id ProductID) String() string  { return string(id) }
func (id ProviderID) String() string { return string(id) }
func (id OrderID) String() string    { return string(id) }

func (id TenantID) IsZero() bool   { return id == "" }
func (id ProductID) IsZero() bool  { return id == "" }
func (id ProviderID) IsZero() bool { return id == "" }

func ParseTenantID(s string) (TenantID, error) {
	if _, err := uuid.Parse(s); err != nil {
		return "", fmt.Errorf("shared: tenant id inválido: %w", err)
	}
	return TenantID(s), nil
}

func ParseProductID(s string) (ProductID, error) {
	if _, err := uuid.Parse(s); err != nil {
		return "", fmt.Errorf("shared: product id inválido: %w", err)
	}
	return ProductID(s), nil
}


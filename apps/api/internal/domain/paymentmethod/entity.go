package paymentmethod

import (
	"fmt"
	"strings"
	"time"

	"uvico-api/internal/domain/shared"
)

type Method string

const (
	MethodPagoMovil Method = "pago_movil"
	MethodZelle     Method = "zelle"
	MethodUSDT      Method = "usdt"
	MethodTransfer  Method = "transfer"
	MethodCash      Method = "cash"
)

func (m Method) IsValid() bool {
	switch m {
	case MethodPagoMovil, MethodZelle, MethodUSDT, MethodTransfer, MethodCash:
		return true
	}
	return false
}

// Config es el payload JSONB que varía por método. La estructura concreta
// se valida en New/validateConfig según el tipo.
type Config map[string]any

// PaymentMethod es un método que el comerciante habilita para su tienda.
// Un mismo tenant puede tener varios, cada uno con su propia configuración.
type PaymentMethod struct {
	TenantID  shared.TenantID
	Method    Method
	Enabled   bool
	Position  int
	Config    Config
	CreatedAt time.Time
	UpdatedAt time.Time
}

// New valida el método y la configuración específica.
func New(tenantID shared.TenantID, method Method, position int, cfg Config, now time.Time) (*PaymentMethod, error) {
	if tenantID.IsZero() {
		return nil, fmt.Errorf("%w: tenant requerido", ErrInvalidMethod)
	}
	if !method.IsValid() {
		return nil, fmt.Errorf("%w: %q", ErrInvalidMethod, method)
	}
	if position < 0 || position > 20 {
		return nil, fmt.Errorf("%w: posición fuera de rango", ErrInvalidConfig)
	}
	if err := validateConfig(method, cfg); err != nil {
		return nil, err
	}
	t := now.UTC()
	return &PaymentMethod{
		TenantID:  tenantID,
		Method:    method,
		Enabled:   true,
		Position:  position,
		Config:    cfg,
		CreatedAt: t,
		UpdatedAt: t,
	}, nil
}

// PublicView devuelve la config con datos enmascarados según el método.
// El comprador nunca ve la cédula completa, solo los últimos 4 dígitos.
func (p *PaymentMethod) PublicView() Config {
	out := make(Config, len(p.Config))
	for k, v := range p.Config {
		out[k] = v
	}
	switch p.Method {
	case MethodPagoMovil:
		if phone, ok := out["phone"].(string); ok && len(phone) >= 4 {
			out["phone"] = "****" + phone[len(phone)-4:]
		}
		if cedula, ok := out["cedula"].(string); ok && len(cedula) >= 4 {
			out["cedula"] = "****" + cedula[len(cedula)-4:]
		}
	case MethodZelle:
		if email, ok := out["email"].(string); ok && strings.Contains(email, "@") {
			parts := strings.SplitN(email, "@", 2)
			if len(parts[0]) > 2 {
				out["email"] = parts[0][:2] + "***@" + parts[1]
			}
		}
	case MethodUSDT:
		if addr, ok := out["address"].(string); ok && len(addr) > 10 {
			out["address"] = addr[:6] + "..." + addr[len(addr)-4:]
		}
	case MethodTransfer:
		if account, ok := out["account"].(string); ok && len(account) >= 4 {
			out["account"] = "****" + account[len(account)-4:]
		}
	}
	return out
}

func validateConfig(method Method, cfg Config) error {
	switch method {
	case MethodPagoMovil:
		return requireFields(cfg, "bank", "phone", "cedula", "holder")
	case MethodZelle:
		return requireFields(cfg, "email", "holder")
	case MethodUSDT:
		return requireFields(cfg, "network", "address")
	case MethodTransfer:
		return requireFields(cfg, "bank", "account", "holder")
	case MethodCash:
		return nil
	}
	return nil
}

func requireFields(cfg Config, fields ...string) error {
	var missing []string
	for _, f := range fields {
		v, ok := cfg[f]
		if !ok {
			missing = append(missing, f)
			continue
		}
		s, _ := v.(string)
		if strings.TrimSpace(s) == "" {
			missing = append(missing, f)
		}
	}
	if len(missing) > 0 {
		return fmt.Errorf("%w: faltan %v", ErrInvalidConfig, missing)
	}
	return nil
}
package tenant

import (
	"errors"
	"fmt"
	"strings"
	"time"
	"unicode/utf8"

	"uvico-api/internal/domain/shared"
)

type Status string

const (
	StatusTrial           Status = "trial"
	StatusActive          Status = "active"
	StatusPastDue         Status = "past_due"
	StatusPaused          Status = "paused"
	StatusPendingDeletion Status = "pending_deletion"
	StatusDeleted         Status = "deleted"
)

var (
	ErrInvalidName    = errors.New("tenant: nombre inválido")
	ErrInvalidSlug    = errors.New("tenant: slug inválido")
	ErrInvalidOwner   = errors.New("tenant: owner_user_id requerido")
	ErrSlugTaken      = errors.New("tenant: slug ya existe")
	ErrNotFound       = errors.New("tenant: no encontrado")
	ErrVersionConflict = errors.New("tenant: version conflict")
)

const (
	MaxNameLen = 120
	MaxSlugLen = 60
)

type Tenant struct {
	ID          shared.TenantID
	Name        string
	Slug        string
	OwnerUserID string
	Status      Status
	Version     int64
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func New(id shared.TenantID, name, slug, ownerUserID string, now time.Time) (*Tenant, error) {
	if id.IsZero() {
		return nil, errors.New("tenant: id requerido")
	}
	name = strings.TrimSpace(name)
	if name == "" || utf8.RuneCountInString(name) > MaxNameLen {
		return nil, fmt.Errorf("%w: 1..%d caracteres", ErrInvalidName, MaxNameLen)
	}
	if !isValidSlug(slug) {
		return nil, fmt.Errorf("%w: %q", ErrInvalidSlug, slug)
	}
	if ownerUserID == "" {
		return nil, ErrInvalidOwner
	}
	t := now.UTC()
	return &Tenant{
		ID:          id,
		Name:        name,
		Slug:        slug,
		OwnerUserID: ownerUserID,
		Status:      StatusTrial,
		Version:     1,
		CreatedAt:   t,
		UpdatedAt:   t,
	}, nil
}

func (t *Tenant) SetStatus(s Status, now time.Time) error {
	switch s {
	case StatusTrial, StatusActive, StatusPastDue, StatusPaused, StatusPendingDeletion, StatusDeleted:
	default:
		return fmt.Errorf("tenant: status desconocido %q", s)
	}
	if t.Status == s {
		return nil
	}
	t.Status = s
	t.UpdatedAt = now.UTC()
	t.Version++
	return nil
}

// CanAcceptOrders retorna true solo si el estado permite vender.
// past_due sigue vendiendo (T&C §9.2: 5 días de gracia con banner).
func (t *Tenant) CanAcceptOrders() bool {
	switch t.Status {
	case StatusTrial, StatusActive, StatusPastDue:
		return true
	}
	return false
}

// CanAccessDashboard retorna true si el panel debe estar habilitado.
func (t *Tenant) CanAccessDashboard() bool {
	return t.Status != StatusDeleted
}

func isValidSlug(s string) bool {
	if s == "" || len(s) > MaxSlugLen {
		return false
	}
	for i, r := range s {
		switch {
		case r >= 'a' && r <= 'z':
		case r >= '0' && r <= '9':
		case r == '-':
		default:
			return false
		}
		if i == 0 && r == '-' {
			return false
		}
	}
	return s[len(s)-1] != '-'
}
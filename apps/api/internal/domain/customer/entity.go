package customer

import (
	"errors"
	"strings"
	"time"
	"unicode/utf8"
)

var (
	ErrInvalidEmail = errors.New("customer: email inválido")
	ErrInvalidName  = errors.New("customer: nombre inválido")
	ErrNotFound     = errors.New("customer: no encontrado")
)

const MaxNameLen = 120

type Customer struct {
	ID               string
	Email            string
	Phone            string
	Name             string
	MarketingOptIn   bool
	MarketingOptInAt *time.Time
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

// UpsertInput es el payload canónico para crear/actualizar un customer.
// Email se normaliza a lowercase. MarketingOptIn solo pasa a true, nunca
// revierte (el endpoint de preferencias para revocar no existe aún).
type UpsertInput struct {
	Email          string
	Name           string
	Phone          string
	MarketingOptIn bool
}

func (in *UpsertInput) Normalize() error {
	in.Email = strings.TrimSpace(strings.ToLower(in.Email))
	in.Name = strings.TrimSpace(in.Name)
	in.Phone = strings.TrimSpace(in.Phone)

	if !isValidEmail(in.Email) {
		return ErrInvalidEmail
	}
	if in.Name == "" || utf8.RuneCountInString(in.Name) > MaxNameLen {
		return ErrInvalidName
	}
	return nil
}

func isValidEmail(s string) bool {
	if s == "" || len(s) > 254 {
		return false
	}
	at := strings.LastIndex(s, "@")
	if at <= 0 || at == len(s)-1 {
		return false
	}
	return strings.Contains(s[at+1:], ".")
}
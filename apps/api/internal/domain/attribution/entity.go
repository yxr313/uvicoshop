package attribution

import (
	"errors"
	"strings"
	"time"
	"unicode/utf8"
)

var (
	ErrInvalidField = errors.New("attribution: campo excede el límite")
)

const (
	MaxUTMLen     = 500
	MaxClickIDLen = 1000
	MaxURLLen     = 2048
	MaxUserAgent  = 512
	MaxLanguage   = 64
	MaxTimezone   = 64
	MaxSessionID  = 128
)

// Attribution captura los datos de marketing de una compra. Todos los
// campos son opcionales: un checkout sin attribution es válido.
type Attribution struct {
	SessionID      string
	UTMSource      string
	UTMMedium      string
	UTMCampaign    string
	UTMTerm        string
	UTMContent     string
	FBCLID         string
	GCLID          string
	TTCLID         string
	TWCLID         string
	Referrer       string
	LandingPage    string
	UserAgent      string
	AcceptLanguage string
	Timezone       string
	IPHash         string
}

// Validate aplica límites de tamaño. No rechaza el checkout por campos
// faltantes: la attribution es opcional.
func (a *Attribution) Validate() error {
	trim := func(s *string) string {
		*s = strings.TrimSpace(*s)
		return *s
	}

	if !within(trim(&a.SessionID), MaxSessionID) {
		return fieldErr("session_id", MaxSessionID)
	}
	for _, f := range []struct {
		name  string
		val   *string
		limit int
	}{
		{"utm_source", &a.UTMSource, MaxUTMLen},
		{"utm_medium", &a.UTMMedium, MaxUTMLen},
		{"utm_campaign", &a.UTMCampaign, MaxUTMLen},
		{"utm_term", &a.UTMTerm, MaxUTMLen},
		{"utm_content", &a.UTMContent, MaxUTMLen},
		{"fbclid", &a.FBCLID, MaxClickIDLen},
		{"gclid", &a.GCLID, MaxClickIDLen},
		{"ttclid", &a.TTCLID, MaxClickIDLen},
		{"twclid", &a.TWCLID, MaxClickIDLen},
		{"referrer", &a.Referrer, MaxURLLen},
		{"landing_page", &a.LandingPage, MaxURLLen},
		{"user_agent", &a.UserAgent, MaxUserAgent},
		{"accept_language", &a.AcceptLanguage, MaxLanguage},
		{"timezone", &a.Timezone, MaxTimezone},
	} {
		if !within(trim(f.val), f.limit) {
			return fieldErr(f.name, f.limit)
		}
	}
	return nil
}

func (a Attribution) IsEmpty() bool {
	return a.SessionID == "" &&
		a.UTMSource == "" && a.UTMMedium == "" && a.UTMCampaign == "" &&
		a.UTMTerm == "" && a.UTMContent == "" &&
		a.FBCLID == "" && a.GCLID == "" && a.TTCLID == "" && a.TWCLID == "" &&
		a.Referrer == "" && a.LandingPage == ""
}

func within(s string, limit int) bool {
	if s == "" {
		return true
	}
	return utf8.RuneCountInString(s) <= limit
}

func fieldErr(name string, limit int) error {
	return &FieldError{Field: name, Limit: limit}
}

type FieldError struct {
	Field string
	Limit int
}

func (e *FieldError) Error() string {
	return "attribution: " + e.Field + " excede " + itoa(e.Limit) + " caracteres"
}

func (e *FieldError) Unwrap() error { return ErrInvalidField }

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	var buf [12]byte
	i := len(buf)
	for n > 0 {
		i--
		buf[i] = byte('0' + n%10)
		n /= 10
	}
	return string(buf[i:])
}

// Record es el DTO que persiste el adapter Postgres.
type Record struct {
	OrderID        string
	TenantID       string
	SessionID      string
	UTMSource      string
	UTMMedium      string
	UTMCampaign    string
	UTMTerm        string
	UTMContent     string
	FBCLID         string
	GCLID          string
	TTCLID         string
	TWCLID         string
	Referrer       string
	LandingPage    string
	UserAgent      string
	AcceptLanguage string
	Timezone       string
	IPHash         string
	CreatedAt      time.Time
}
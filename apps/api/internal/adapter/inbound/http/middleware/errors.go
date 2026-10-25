package middleware

import (
	"errors"
	"net/http"
)

var (
	errMissingUser   = errors.New("middleware: user_id no presente en contexto")
	errMissingTenant = errors.New("middleware: tenant_id no presente en contexto")
)

func writeJSONError(w http.ResponseWriter, status int, code, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-store")
	w.WriteHeader(status)
	_, _ = w.Write([]byte(`{"error":"` + code + `","message":"` + message + `"}`))
}
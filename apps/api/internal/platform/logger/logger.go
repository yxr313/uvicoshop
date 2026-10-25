package logger

import (
	"context"
	"log/slog"
	"os"

	"go.opentelemetry.io/otel/trace"

	"uvico-api/internal/platform/config"
)

const (
	traceIDKey = "trace_id"
	spanIDKey  = "span_id"
)

// New construye el logger del servicio. En producción usa JSON estructurado
// (parseable por Datadog, Grafana Loki, etc.). En desarrollo usa texto legible.
// En ambos casos enriquece con trace_id/span_id cuando hay span OTel activo.
func New(cfg config.ObservabilityConfig) *slog.Logger {
	opts := &slog.HandlerOptions{Level: parseLevel(cfg.LogLevel)}

	var base slog.Handler
	if os.Getenv("APP_ENV") == "production" {
		base = slog.NewJSONHandler(os.Stdout, opts)
	} else {
		base = slog.NewTextHandler(os.Stdout, opts)
	}

	return slog.New(&traceHandler{Handler: base})
}

func parseLevel(s string) slog.Level {
	var l slog.Level
	if err := l.UnmarshalText([]byte(s)); err != nil {
		return slog.LevelInfo
	}
	return l
}

type traceHandler struct {
	slog.Handler
}

func (h *traceHandler) Handle(ctx context.Context, r slog.Record) error {
	if sc := trace.SpanContextFromContext(ctx); sc.IsValid() {
		r.AddAttrs(
			slog.String(traceIDKey, sc.TraceID().String()),
			slog.String(spanIDKey, sc.SpanID().String()),
		)
	}
	return h.Handler.Handle(ctx, r)
}

func (h *traceHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &traceHandler{Handler: h.Handler.WithAttrs(attrs)}
}

func (h *traceHandler) WithGroup(name string) slog.Handler {
	return &traceHandler{Handler: h.Handler.WithGroup(name)}
}
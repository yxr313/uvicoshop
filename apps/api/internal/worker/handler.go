package worker

import (
	"context"
	"fmt"
	"sort"
	"sync"
)

// Event es un registro de outbox_events deserializado.
// Payload es el JSONB crudo; cada handler decide cómo decodificarlo.
type Event struct {
	ID            string
	TenantID      string
	AggregateType string
	AggregateID   string
	EventType     string
	Payload       []byte
	Attempts      int32
	MaxAttempts   int32
}

// Handler procesa un evento. Si retorna error, el worker
// incrementa attempts y agenda el reintento con backoff.
type Handler interface {
	Handle(ctx context.Context, event Event) error
}

// HandlerFunc permite usar funciones sueltas como Handler.
type HandlerFunc func(ctx context.Context, event Event) error

func (f HandlerFunc) Handle(ctx context.Context, event Event) error {
	return f(ctx, event)
}

// Registry mapea event_type -> Handler. Es thread-safe.
// Fail-fast: panics en Register() por duplicados o nil.
// Preferible a silenciar conflictos de configuración.
type Registry struct {
	mu       sync.RWMutex
	handlers map[string]Handler
}

func NewRegistry() *Registry {
	return &Registry{handlers: make(map[string]Handler)}
}

func (r *Registry) Register(eventType string, h Handler) {
	if eventType == "" {
		panic("worker: event_type vacío")
	}
	if h == nil {
		panic("worker: handler nil para " + eventType)
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.handlers[eventType]; exists {
		panic("worker: handler duplicado para " + eventType)
	}
	r.handlers[eventType] = h
}

func (r *Registry) Lookup(eventType string) (Handler, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	h, ok := r.handlers[eventType]
	return h, ok
}

func (r *Registry) Names() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]string, 0, len(r.handlers))
	for k := range r.handlers {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}

// Verify asegura que todos los event_types de una lista tengan handler.
// Usar en startup para detectar config faltante antes de servir tráfico.
func (r *Registry) Verify(required []string) error {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var missing []string
	for _, et := range required {
		if _, ok := r.handlers[et]; !ok {
			missing = append(missing, et)
		}
	}
	if len(missing) > 0 {
		return fmt.Errorf("worker: handlers faltantes: %v", missing)
	}
	return nil
}
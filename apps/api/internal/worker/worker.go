package worker

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"uvico-api/internal/platform/database"
)

const (
	defaultBatchSize    = 100
	defaultPollInterval = 2 * time.Second
	defaultReapEvery    = time.Minute
	defaultLockTimeout  = 5 * time.Minute
	defaultHandlerTO    = 60 * time.Second
	maxBackoffShift     = 12 // 2^12 = 4096s ≈ 68 min
)

// Metrics es opcional. Si se provee, el worker emite un punto por
// cada ejecución de handler. Enlazar con Prometheus en main.go.
type Metrics interface {
	Observe(eventType, status string, duration time.Duration)
}

type Options struct {
	Pool         *database.Pool
	Logger       *slog.Logger
	Registry     *Registry
	WorkerID     string
	BatchSize    int
	PollInterval time.Duration
	ReapEvery    time.Duration
	LockTimeout  time.Duration
	HandlerTO    time.Duration
	Metrics      Metrics
}

type Worker struct {
	opts Options
}

func New(opts Options) *Worker {
	if opts.Registry == nil {
		panic("worker: Registry requerido")
	}
	if opts.Pool == nil {
		panic("worker: Pool requerido")
	}
	if opts.BatchSize <= 0 {
		opts.BatchSize = defaultBatchSize
	}
	if opts.PollInterval <= 0 {
		opts.PollInterval = defaultPollInterval
	}
	if opts.ReapEvery <= 0 {
		opts.ReapEvery = defaultReapEvery
	}
	if opts.LockTimeout <= 0 {
		opts.LockTimeout = defaultLockTimeout
	}
	if opts.HandlerTO <= 0 {
		opts.HandlerTO = defaultHandlerTO
	}
	if opts.WorkerID == "" {
		opts.WorkerID = "worker-" + uuid.Must(uuid.NewV7()).String()
	}
	if opts.Logger == nil {
		opts.Logger = slog.Default()
	}
	return &Worker{opts: opts}
}

func (w *Worker) WorkerID() string { return w.opts.WorkerID }

// Run arranca el poller y el reaper hasta que ctx se cancele.
// Retorna ctx.Err() para que main sepa si fue shutdown limpio.
func (w *Worker) Run(ctx context.Context) error {
	var wg sync.WaitGroup
	wg.Add(2)
	go func() { defer wg.Done(); w.pollLoop(ctx) }()
	go func() { defer wg.Done(); w.reapLoop(ctx) }()
	wg.Wait()
	return ctx.Err()
}

// ---------------------------------------------------------------------------
// Poll loop
// ---------------------------------------------------------------------------

func (w *Worker) pollLoop(ctx context.Context) {
	ticker := time.NewTicker(w.opts.PollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := w.pollOnce(ctx); err != nil && !errors.Is(err, context.Canceled) {
				w.opts.Logger.Error("worker: poll failed", "err", err)
			}
		}
	}
}

func (w *Worker) pollOnce(ctx context.Context) error {
	events, err := w.claimBatch(ctx, w.opts.BatchSize)
	if err != nil {
		return fmt.Errorf("claim: %w", err)
	}
	if len(events) == 0 {
		return nil
	}

	w.opts.Logger.Info("worker: batch claimed", "count", len(events), "worker_id", w.opts.WorkerID)

	for i := range events {
		if ctx.Err() != nil {
			return ctx.Err()
		}
		w.handleOne(ctx, &events[i])
	}
	return nil
}

// claimBatch hace dos cosas:
//  1. SELECT ... FOR UPDATE SKIP LOCKED (dentro de una transacción)
//  2. UPDATE locked_at, locked_by para los IDs obtenidos
//  3. COMMIT (suelta locks de Postgres, mantiene lock lógico vía locked_at)
//
// El handler procesa SIN locks de DB. Si el worker muere, el reaper
// libera locked_at después de 2*LockTimeout.
func (w *Worker) claimBatch(ctx context.Context, limit int) ([]Event, error) {
	tx, err := w.opts.Pool.BeginTx(ctx, pgx.TxOptions{
		IsoLevel:   pgx.ReadCommitted,
		AccessMode: pgx.ReadWrite,
	})
	if err != nil {
		return nil, fmt.Errorf("begin: %w", err)
	}
	defer func() {
		rbCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), 5*time.Second)
		defer cancel()
		_ = tx.Rollback(rbCtx)
	}()

	rows, err := tx.Query(ctx, `
		SELECT id, tenant_id, aggregate_type, aggregate_id, event_type, payload, attempts, max_attempts
		  FROM outbox_events
		 WHERE processed_at IS NULL
		   AND attempts < max_attempts
		   AND next_attempt_at <= NOW()
		   AND (locked_at IS NULL OR locked_at < NOW() - $1::interval)
		 ORDER BY next_attempt_at, created_at
		 LIMIT $2
		 FOR UPDATE SKIP LOCKED
	`, w.opts.LockTimeout.String(), limit)
	if err != nil {
		return nil, fmt.Errorf("select: %w", err)
	}
	defer rows.Close()

	var events []Event
	for rows.Next() {
		var e Event
		if err := rows.Scan(
			&e.ID, &e.TenantID, &e.AggregateType, &e.AggregateID,
			&e.EventType, &e.Payload, &e.Attempts, &e.MaxAttempts,
		); err != nil {
			return nil, fmt.Errorf("scan: %w", err)
		}
		events = append(events, e)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows: %w", err)
	}

	if len(events) == 0 {
		return nil, tx.Commit(ctx)
	}

	ids := make([]string, len(events))
	for i, e := range events {
		ids[i] = e.ID
	}
	if _, err := tx.Exec(ctx, `
		UPDATE outbox_events
		   SET locked_at = NOW(), locked_by = $1
		 WHERE id = ANY($2)
	`, w.opts.WorkerID, ids); err != nil {
		return nil, fmt.Errorf("lock: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("commit: %w", err)
	}
	return events, nil
}

// ---------------------------------------------------------------------------
// Dispatch
// ---------------------------------------------------------------------------

func (w *Worker) handleOne(ctx context.Context, e *Event) {
	start := time.Now()

	handler, ok := w.opts.Registry.Lookup(e.EventType)
	if !ok {
		w.opts.Logger.Error("worker: sin handler registrado",
			"event_id", e.ID, "event_type", e.EventType)
		w.markFailed(ctx, e, fmt.Errorf("worker: sin handler para %q", e.EventType))
		w.observe(e.EventType, "no_handler", time.Since(start))
		return
	}

	hctx, cancel := context.WithTimeout(ctx, w.opts.HandlerTO)
	defer cancel()

	err := handler.Handle(hctx, *e)
	dur := time.Since(start)

	if err != nil {
		w.opts.Logger.Warn("worker: handler falló",
			"event_id", e.ID,
			"event_type", e.EventType,
			"attempt", e.Attempts+1,
			"max_attempts", e.MaxAttempts,
			"duration_ms", dur.Milliseconds(),
			"err", err,
		)
		w.markFailed(ctx, e, err)
		w.observe(e.EventType, "failed", dur)
		return
	}

	w.markProcessed(ctx, e)
	w.opts.Logger.Info("worker: evento procesado",
		"event_id", e.ID,
		"event_type", e.EventType,
		"duration_ms", dur.Milliseconds(),
	)
	w.observe(e.EventType, "processed", dur)
}

func (w *Worker) observe(eventType, status string, d time.Duration) {
	if w.opts.Metrics != nil {
		w.opts.Metrics.Observe(eventType, status, d)
	}
}

// ---------------------------------------------------------------------------
// State transitions (siempre fuera de locks)
// ---------------------------------------------------------------------------

func (w *Worker) markProcessed(ctx context.Context, e *Event) {
	_, err := w.opts.Pool.Exec(ctx, `
		UPDATE outbox_events
		   SET processed_at = NOW(),
		       locked_at = NULL,
		       locked_by = NULL
		 WHERE id = $1
	`, e.ID)
	if err != nil {
		w.opts.Logger.Error("worker: mark processed failed", "event_id", e.ID, "err", err)
	}
}

// markFailed aplica backoff exponencial: 2^(attempts+1) segundos, cap 68 min.
// Con max_attempts=10, un evento fallido reintenta durante ~3.5h antes de DLQ.
func (w *Worker) markFailed(ctx context.Context, e *Event, cause error) {
	shift := e.Attempts + 1
	if shift > maxBackoffShift {
		shift = maxBackoffShift
	}
	backoffSeconds := int64(1) << shift

	_, err := w.opts.Pool.Exec(ctx, `
		UPDATE outbox_events
		   SET attempts        = attempts + 1,
		       next_attempt_at = NOW() + (INTERVAL '1 second' * $2),
		       last_error      = $3,
		       locked_at       = NULL,
		       locked_by       = NULL
		 WHERE id = $1
	`, e.ID, backoffSeconds, cause.Error())
	if err != nil {
		w.opts.Logger.Error("worker: mark failed failed", "event_id", e.ID, "err", err)
	}
}

// ---------------------------------------------------------------------------
// Reaper
// ---------------------------------------------------------------------------

func (w *Worker) reapLoop(ctx context.Context) {
	ticker := time.NewTicker(w.opts.ReapEvery)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			n, err := w.reapOnce(ctx)
			if err != nil && !errors.Is(err, context.Canceled) {
				w.opts.Logger.Error("worker: reap failed", "err", err)
				continue
			}
			if n > 0 {
				w.opts.Logger.Warn("worker: locks huérfanos liberados", "count", n)
			}
		}
	}
}

// reapOnce libera locks de workers que crashearon con kill -9.
// No interfiere con workers sanos: solo libera locked_at más viejos
// que 2*LockTimeout.
func (w *Worker) reapOnce(ctx context.Context) (int64, error) {
	tag, err := w.opts.Pool.Exec(ctx, `
		UPDATE outbox_events
		   SET locked_at = NULL,
		       locked_by = NULL
		 WHERE processed_at IS NULL
		   AND locked_at IS NOT NULL
		   AND locked_at < NOW() - $1::interval
	`, (2 * w.opts.LockTimeout).String())
	if err != nil {
		return 0, fmt.Errorf("reap: %w", err)
	}
	return tag.RowsAffected(), nil
}
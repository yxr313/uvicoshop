package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/exaring/otelpgx"

	"uvico-api/internal/platform/config"
	"uvico-api/internal/platform/database"
	"uvico-api/internal/platform/logger"
	"uvico-api/internal/platform/telemetry"
	"uvico-api/internal/worker"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "fatal: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	cfg := config.MustLoad()

	log := logger.New(cfg.Observability)
	slog.SetDefault(log)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	shutdown, err := telemetry.Setup(ctx, cfg.Observability)
	if err != nil {
		return fmt.Errorf("init telemetry: %w", err)
	}
	defer func() { _ = shutdown(context.Background()) }()

	pool, err := database.New(ctx, cfg.Database, database.Options{
		Logger: log,
		Tracer: otelpgx.NewTracer(),
	})
	if err != nil {
		return fmt.Errorf("init database: %w", err)
	}
	defer pool.Close()

	registry := worker.NewRegistry()
	registerHandlers(registry, log, pool)

	w := worker.New(worker.Options{
		Pool:         pool,
		Logger:       log,
		Registry:     registry,
		WorkerID:     cfg.Service.Name + "-" + cfg.Service.Region + "-" + cfg.Server.Port,
		BatchSize:    cfg.Messaging.OutboxBatchSize,
		PollInterval: cfg.Messaging.OutboxPollInterval,
	})

	log.Info("worker starting",
		"worker_id", w.WorkerID(),
		"batch_size", cfg.Messaging.OutboxBatchSize,
		"poll_interval", cfg.Messaging.OutboxPollInterval.String(),
		"handlers", registry.Names(),
	)

	if err := w.Run(ctx); err != nil && !errors.Is(err, context.Canceled) {
		return fmt.Errorf("worker: %w", err)
	}

	log.Info("worker stopped")
	return nil
}

func registerHandlers(r *worker.Registry, log *slog.Logger, pool *database.Pool) {
	// Fase 7: handlers reales.
	// Por ahora, un noop para que el worker arranque y el reaper esté activo:
	r.Register("product.created", worker.HandlerFunc(func(ctx context.Context, e worker.Event) error {
		log.Info("product.created", "product_id", e.AggregateID, "tenant_id", e.TenantID)
		return nil
	}))
}
package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/exaring/otelpgx"

	"uvico-api/internal/adapter/inbound/http/router"
	"uvico-api/internal/domain/billing"
	"uvico-api/internal/platform/config"
	"uvico-api/internal/platform/database"
	"uvico-api/internal/platform/logger"
	"uvico-api/internal/platform/redisclient"
	"uvico-api/internal/platform/telemetry"
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
	defer func() {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		_ = shutdown(shutdownCtx)
	}()

	pool, err := database.New(ctx, cfg.Database, database.Options{
		Logger: log,
		Tracer: otelpgx.NewTracer(),
	})
	if err != nil {
		return fmt.Errorf("init database: %w", err)
	}
	defer pool.Close()

	rdb, err := redisclient.New(ctx, cfg.Redis)
	if err != nil {
		return fmt.Errorf("init redis: %w", err)
	}
	defer func() {
		if err := rdb.Close(); err != nil {
			log.Error("redis close", "err", err)
		}
	}()

	plans := billing.Plans(
		cfg.Plans.FounderMonthlyCents,
		cfg.Plans.FounderTrialCents,
		cfg.Plans.EssentialMonthlyCents,
		cfg.Plans.ProfessionalMonthlyCents,
		cfg.Plans.EnterpriseMonthlyCents,
		cfg.Plans.AnnualDiscountBPS,
	)

	deps := &router.Dependencies{
		Pool:   pool,
		Redis:  rdb,
		Logger: log,
		Cfg:    cfg,
		Plans:  plans,
	}

	handler := router.New(deps)

	srv := &http.Server{
		Addr:              ":" + cfg.Server.Port,
		Handler:           handler,
		ReadTimeout:       cfg.Server.ReadTimeout,
		ReadHeaderTimeout: cfg.Server.ReadHeaderTimeout,
		WriteTimeout:      cfg.Server.WriteTimeout,
		IdleTimeout:       cfg.Server.IdleTimeout,
		MaxHeaderBytes:    cfg.Server.MaxHeaderBytes,
	}

	serverErr := make(chan error, 1)
	go func() {
		log.Info("http server listening",
			"addr", srv.Addr,
			"env", cfg.Env,
			"version", cfg.Service.Version,
			"region", cfg.Service.Region,
		)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErr <- err
		}
		close(serverErr)
	}()

	select {
	case err := <-serverErr:
		return fmt.Errorf("http server: %w", err)
	case <-ctx.Done():
		log.Info("shutdown signal received")
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.Server.ShutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Error("graceful shutdown failed, forcing close", "err", err)
		if closeErr := srv.Close(); closeErr != nil {
			return fmt.Errorf("shutdown: %w (close: %v)", err, closeErr)
		}
		return fmt.Errorf("shutdown: %w", err)
	}

	log.Info("http server stopped")
	return nil
}
package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/exaring/otelpgx"
	"github.com/redis/go-redis/v9"

	"uvico-api/internal/cron"
	bcvscraper "uvico-api/internal/cron/jobs/bcv"
	billingjobs "uvico-api/internal/cron/jobs/billing"
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
	defer func() { _ = shutdown(context.Background()) }()

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

	lock := cron.NewRedisLock(rdb, log)
	sched := cron.NewScheduler(cron.SchedulerOptions{
		Logger:   log,
		Lock:     lock,
		Location: time.UTC,
	})

	registerJobs(sched, cfg, pool, rdb, log)

	if err := sched.Start(); err != nil {
		return fmt.Errorf("start scheduler: %w", err)
	}

	<-ctx.Done()
	log.Info("shutdown signal received, stopping scheduler")

	stopCtx, stopCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer stopCancel()
	if err := sched.Stop(stopCtx); err != nil {
		return fmt.Errorf("stop scheduler: %w", err)
	}

	log.Info("scheduler stopped")
	return nil
}

func registerJobs(
	s *cron.Scheduler,
	cfg *config.Config,
	pool *database.Pool,
	rdb *redis.Client,
	log *slog.Logger,
) {
	s.MustAdd(cron.JobSpec{
		Name:       "bcv-scraper",
		Spec:       cfg.BCV.ScraperCron,
		RunOnStart: true,
		Job:        bcvscraper.New(cfg.BCV, rdb, log),
	})

	s.MustAdd(cron.JobSpec{
		Name: "billing-trial-expiration",
		Spec: "*/5 * * * *",
		Job:  billingjobs.NewTrialExpiration(pool, log),
	})

	s.MustAdd(cron.JobSpec{
		Name: "billing-grace-period",
		Spec: "*/5 * * * *",
		Job:  billingjobs.NewGracePeriod(pool, log),
	})

	s.MustAdd(cron.JobSpec{
		Name: "billing-deletion",
		Spec: "0 3 * * *",
		Job:  billingjobs.NewDeletion(pool, log),
	})

	// TODO Msg 5: outbox-partitions, inventory-reconcile, abandoned-carts,
	// billing-founder-upgrade, payout-weekly, payout-fortnightly.
}
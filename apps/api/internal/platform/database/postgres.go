package database

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"uvico-api/internal/platform/config"
)

type Pool struct {
	*pgxpool.Pool
	cfg    config.DatabaseConfig
	logger *slog.Logger
}

type Options struct {
	Logger *slog.Logger
	Tracer pgx.QueryTracer
}

type TxOptions struct {
	IsoLevel   pgx.TxIsoLevel
	AccessMode pgx.TxAccessMode
}

func New(ctx context.Context, cfg config.DatabaseConfig, opts Options) (*Pool, error) {
	if cfg.URL == "" {
		return nil, errors.New("database: DATABASE_URL vacío")
	}
	logger := opts.Logger
	if logger == nil {
		logger = slog.Default()
	}

	poolCfg, err := pgxpool.ParseConfig(cfg.URL)
	if err != nil {
		return nil, fmt.Errorf("database: parse config: %w", err)
	}

	poolCfg.MaxConns = cfg.MaxConns
	poolCfg.MinConns = cfg.MinConns
	poolCfg.MaxConnLifetime = cfg.MaxConnLifetime
	poolCfg.MaxConnIdleTime = cfg.MaxConnIdleTime
	poolCfg.HealthCheckPeriod = cfg.HealthCheck
	poolCfg.ConnConfig.ConnectTimeout = cfg.ConnectTimeout

	poolCfg.ConnConfig.RuntimeParams["statement_timeout"] =
		strconv.FormatInt(cfg.StatementTimeout.Milliseconds(), 10)
	poolCfg.ConnConfig.RuntimeParams["idle_in_transaction_session_timeout"] =
		strconv.FormatInt((cfg.StatementTimeout * 4).Milliseconds(), 10)
	poolCfg.ConnConfig.RuntimeParams["application_name"] = "uvico-api"
	poolCfg.ConnConfig.DefaultQueryExecMode = pgx.QueryExecModeExec

	if opts.Tracer != nil {
		poolCfg.ConnConfig.Tracer = opts.Tracer
	}

	pool, err := pgxpool.NewWithConfig(ctx, poolCfg)
	if err != nil {
		return nil, fmt.Errorf("database: create pool: %w", err)
	}

	pingCtx, cancel := context.WithTimeout(ctx, cfg.ConnectTimeout)
	defer cancel()
	if err := pool.Ping(pingCtx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("database: ping: %w", err)
	}

	logger.Info("database pool ready",
		"max_conns", cfg.MaxConns,
		"min_conns", cfg.MinConns,
		"statement_timeout", cfg.StatementTimeout.String(),
	)
	return &Pool{Pool: pool, cfg: cfg, logger: logger}, nil
}

func (p *Pool) Ping(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, p.cfg.ConnectTimeout)
	defer cancel()
	return p.Pool.Ping(ctx)
}

func (p *Pool) WithTenantTx(ctx context.Context, tenantID string, fn func(ctx context.Context, tx pgx.Tx) error) error {
	return p.WithTenantTxOpts(ctx, tenantID, TxOptions{
		IsoLevel:   pgx.ReadCommitted,
		AccessMode: pgx.ReadWrite,
	}, fn)
}

func (p *Pool) WithTenantTxOpts(
	ctx context.Context,
	tenantID string,
	opts TxOptions,
	fn func(ctx context.Context, tx pgx.Tx) error,
) error {
	if tenantID == "" {
		return errors.New("database: tenantID requerido")
	}

	tx, err := p.Pool.BeginTx(ctx, pgx.TxOptions{
		IsoLevel:   opts.IsoLevel,
		AccessMode: opts.AccessMode,
	})
	if err != nil {
		return fmt.Errorf("database: begin tx: %w", err)
	}

	defer func() {
		rbCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), 5*time.Second)
		defer cancel()
		if rbErr := tx.Rollback(rbCtx); rbErr != nil && !errors.Is(rbErr, pgx.ErrTxClosed) {
			p.logger.Error("database: rollback failed", "tenant_id", tenantID, "err", rbErr)
		}
	}()

	if _, err := tx.Exec(ctx, "SELECT set_config('app.tenant_id', $1, true)", tenantID); err != nil {
		return fmt.Errorf("database: set tenant: %w", err)
	}

	if err := fn(ctx, tx); err != nil {
		return err
	}
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("database: commit: %w", err)
	}
	return nil
}

func (p *Pool) Close() {
	p.Pool.Close()
	p.logger.Info("database pool closed")
}
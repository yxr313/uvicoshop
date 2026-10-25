package billing

import (
	"context"
	"log/slog"
	"time"

	"uvico-api/internal/platform/database"
)

const billingTimeout = 2 * time.Minute

// ---------------------------------------------------------------------------
// TrialExpiration: trial -> past_due al vencer trial_ends_at (T&C §5.5)
// ---------------------------------------------------------------------------

type TrialExpiration struct {
	pool   *database.Pool
	logger *slog.Logger
}

func NewTrialExpiration(pool *database.Pool, logger *slog.Logger) *TrialExpiration {
	if logger == nil {
		logger = slog.Default()
	}
	return &TrialExpiration{pool: pool, logger: logger}
}

func (j *TrialExpiration) Timeout() time.Duration { return billingTimeout }

func (j *TrialExpiration) Run(ctx context.Context) error {
	tag, err := j.pool.Exec(ctx, `
		UPDATE subscriptions
		   SET status = 'past_due',
		       past_due_since = NOW(),
		       updated_at = NOW()
		 WHERE status = 'trialing'
		   AND trial_ends_at <= NOW()
	`)
	if err != nil {
		return err
	}
	if n := tag.RowsAffected(); n > 0 {
		j.logger.Info("subscriptions moved to past_due", "count", n)
	}
	return nil
}

// ---------------------------------------------------------------------------
// GracePeriod: past_due -> paused al 6to dia (T&C §9.3)
// ---------------------------------------------------------------------------

type GracePeriod struct {
	pool   *database.Pool
	logger *slog.Logger
}

func NewGracePeriod(pool *database.Pool, logger *slog.Logger) *GracePeriod {
	if logger == nil {
		logger = slog.Default()
	}
	return &GracePeriod{pool: pool, logger: logger}
}

func (j *GracePeriod) Timeout() time.Duration { return billingTimeout }

func (j *GracePeriod) Run(ctx context.Context) error {
	tag, err := j.pool.Exec(ctx, `
		UPDATE subscriptions
		   SET status = 'paused',
		       paused_at = NOW(),
		       updated_at = NOW()
		 WHERE status = 'past_due'
		   AND past_due_since <= NOW() - INTERVAL '6 days'
	`)
	if err != nil {
		return err
	}
	if n := tag.RowsAffected(); n > 0 {
		j.logger.Info("subscriptions paused", "count", n)
	}
	return nil
}

// ---------------------------------------------------------------------------
// Deletion: paused -> pending_deletion a los 25 dias,
//           pending_deletion -> deleted a los 30 dias (DPA §8)
// ---------------------------------------------------------------------------

type Deletion struct {
	pool   *database.Pool
	logger *slog.Logger
}

func NewDeletion(pool *database.Pool, logger *slog.Logger) *Deletion {
	if logger == nil {
		logger = slog.Default()
	}
	return &Deletion{pool: pool, logger: logger}
}

func (j *Deletion) Timeout() time.Duration { return billingTimeout }

func (j *Deletion) Run(ctx context.Context) error {
	if err := j.markPending(ctx); err != nil {
		return err
	}
	return j.finalize(ctx)
}

func (j *Deletion) markPending(ctx context.Context) error {
	tag, err := j.pool.Exec(ctx, `
		UPDATE subscriptions
		   SET status = 'pending_deletion',
		       deletion_scheduled_at = paused_at + INTERVAL '30 days',
		       updated_at = NOW()
		 WHERE status = 'paused'
		   AND paused_at <= NOW() - INTERVAL '25 days'
	`)
	if err != nil {
		return err
	}
	if n := tag.RowsAffected(); n > 0 {
		j.logger.Warn("subscriptions marked pending_deletion (5-day notice)", "count", n)
	}
	return nil
}

func (j *Deletion) finalize(ctx context.Context) error {
	tag, err := j.pool.Exec(ctx, `
		UPDATE subscriptions
		   SET status = 'deleted',
		       updated_at = NOW()
		 WHERE status = 'pending_deletion'
		   AND deletion_scheduled_at <= NOW()
	`)
	if err != nil {
		return err
	}
	if n := tag.RowsAffected(); n > 0 {
		j.logger.Warn("subscriptions finalized (deletion executed)", "count", n)
	}
	return nil
}
package postgres

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"

	"uvico-api/internal/domain/paymentmethod"
	"uvico-api/internal/domain/shared"
	"uvico-api/internal/platform/database"
)

type PaymentMethodRepo struct {
	pool *database.Pool
}

func NewPaymentMethodRepo(pool *database.Pool) *PaymentMethodRepo {
	return &PaymentMethodRepo{pool: pool}
}

func (r *PaymentMethodRepo) Replace(ctx context.Context, tenantID shared.TenantID, methods []*paymentmethod.PaymentMethod) error {
	return r.pool.WithTenantTx(ctx, tenantID.String(), func(ctx context.Context, tx pgx.Tx) error {
		// Disable all first, then upsert the new list. This makes the
		// replace atomic and lets merchants reorder without losing config.
		if _, err := tx.Exec(ctx, `
			UPDATE tenant_payment_methods
			   SET enabled = false, updated_at = NOW()
			 WHERE tenant_id = $1
		`, tenantID.String()); err != nil {
			return fmt.Errorf("paymentmethod_repo: disable all: %w", err)
		}

		for _, m := range methods {
			cfg, err := json.Marshal(m.Config)
			if err != nil {
				return fmt.Errorf("paymentmethod_repo: marshal config: %w", err)
			}
			if _, err := tx.Exec(ctx, `
				INSERT INTO tenant_payment_methods
					(tenant_id, method, enabled, position, config, created_at, updated_at)
				VALUES ($1, $2, $3, $4, $5, $6, $7)
				ON CONFLICT (tenant_id, method) DO UPDATE
				   SET enabled    = EXCLUDED.enabled,
				       position   = EXCLUDED.position,
				       config     = EXCLUDED.config,
				       updated_at = EXCLUDED.updated_at
			`, tenantID.String(), string(m.Method), m.Enabled, m.Position,
				string(cfg), m.CreatedAt, m.UpdatedAt); err != nil {
				return fmt.Errorf("paymentmethod_repo: upsert %s: %w", m.Method, err)
			}
		}
		return nil
	})
}

func (r *PaymentMethodRepo) ListEnabled(ctx context.Context, tenantID shared.TenantID) ([]*paymentmethod.PaymentMethod, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT tenant_id, method, enabled, position, config, created_at, updated_at
		  FROM tenant_payment_methods
		 WHERE tenant_id = $1 AND enabled = true
		 ORDER BY position ASC, method ASC
	`, tenantID.String())
	if err != nil {
		return nil, fmt.Errorf("paymentmethod_repo: list enabled: %w", err)
	}
	defer rows.Close()
	return scanPaymentMethods(rows)
}

func (r *PaymentMethodRepo) ListAll(ctx context.Context, tenantID shared.TenantID) ([]*paymentmethod.PaymentMethod, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT tenant_id, method, enabled, position, config, created_at, updated_at
		  FROM tenant_payment_methods
		 WHERE tenant_id = $1
		 ORDER BY position ASC, method ASC
	`, tenantID.String())
	if err != nil {
		return nil, fmt.Errorf("paymentmethod_repo: list all: %w", err)
	}
	defer rows.Close()
	return scanPaymentMethods(rows)
}

func (r *PaymentMethodRepo) GetByMethod(ctx context.Context, tenantID shared.TenantID, method paymentmethod.Method) (*paymentmethod.PaymentMethod, error) {
	row := r.pool.QueryRow(ctx, `
		SELECT tenant_id, method, enabled, position, config, created_at, updated_at
		  FROM tenant_payment_methods
		 WHERE tenant_id = $1 AND method = $2
	`, tenantID.String(), string(method))
	return scanPaymentMethod(row)
}

func scanPaymentMethod(row rowScanner) (*paymentmethod.PaymentMethod, error) {
	var m paymentmethod.PaymentMethod
	var tid, method string
	var cfg []byte
	err := row.Scan(&tid, &method, &m.Enabled, &m.Position, &cfg, &m.CreatedAt, &m.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, paymentmethod.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("paymentmethod_repo: scan: %w", err)
	}
	m.TenantID = shared.TenantID(tid)
	m.Method = paymentmethod.Method(method)
	m.Config = paymentmethod.Config{}
	if len(cfg) > 0 {
		_ = json.Unmarshal(cfg, &m.Config)
	}
	return &m, nil
}

func scanPaymentMethods(rows pgx.Rows) ([]*paymentmethod.PaymentMethod, error) {
	out := make([]*paymentmethod.PaymentMethod, 0)
	for rows.Next() {
		m, err := scanPaymentMethod(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, m)
	}
	return out, rows.Err()
}
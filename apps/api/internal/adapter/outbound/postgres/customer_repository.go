package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"

	"uvico-api/internal/domain/customer"
	"uvico-api/internal/platform/database"
)

type CustomerRepo struct {
	pool *database.Pool
}

func NewCustomerRepo(pool *database.Pool) *CustomerRepo {
	return &CustomerRepo{pool: pool}
}

func (r *CustomerRepo) UpsertByEmail(ctx context.Context, in customer.UpsertInput) (string, error) {
	if err := in.Normalize(); err != nil {
		return "", err
	}

	var id string
	now := time.Now().UTC()

	var optInAt any
	if in.MarketingOptIn {
		optInAt = now
	}

	err := r.pool.QueryRow(ctx, `
		INSERT INTO customers (id, email, name, phone, marketing_opt_in, marketing_opt_in_at, created_at, updated_at)
		VALUES (gen_random_uuid(), $1, $2, $3, $4, $5, $6, $6)
		ON CONFLICT ((lower(email))) DO UPDATE
		   SET name = EXCLUDED.name,
		       phone = COALESCE(NULLIF(EXCLUDED.phone, ''), customers.phone),
		       marketing_opt_in = CASE
		           WHEN EXCLUDED.marketing_opt_in = true AND customers.marketing_opt_in = false
		           THEN true
		           ELSE customers.marketing_opt_in
		       END,
		       marketing_opt_in_at = CASE
		           WHEN EXCLUDED.marketing_opt_in = true AND customers.marketing_opt_in = false
		           THEN EXCLUDED.marketing_opt_in_at
		           ELSE customers.marketing_opt_in_at
		       END,
		       updated_at = EXCLUDED.updated_at
		RETURNING id
	`, in.Email, in.Name, in.Phone, in.MarketingOptIn, optInAt, now).Scan(&id)
	if err != nil {
		return "", fmt.Errorf("customer_repo: upsert: %w", err)
	}
	return id, nil
}

func (r *CustomerRepo) GetByEmail(ctx context.Context, email string) (*customer.Customer, error) {
	var c customer.Customer
	err := r.pool.QueryRow(ctx, `
		SELECT id, email, COALESCE(name, ''), COALESCE(phone, ''),
		       marketing_opt_in, marketing_opt_in_at, created_at, updated_at
		  FROM customers
		 WHERE lower(email) = lower($1)
	`, email).Scan(
		&c.ID, &c.Email, &c.Name, &c.Phone,
		&c.MarketingOptIn, &c.MarketingOptInAt, &c.CreatedAt, &c.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, customer.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("customer_repo: get: %w", err)
	}
	return &c, nil
}
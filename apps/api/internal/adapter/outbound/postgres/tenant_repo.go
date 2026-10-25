package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"uvico-api/internal/domain/shared"
	"uvico-api/internal/domain/tenant"
	"uvico-api/internal/platform/database"
)

type TenantRepo struct {
	pool *database.Pool
}

func NewTenantRepo(pool *database.Pool) *TenantRepo {
	return &TenantRepo{pool: pool}
}

func (r *TenantRepo) Create(ctx context.Context, t *tenant.Tenant) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO tenants (id, name, slug, owner_user_id, status, created_at, updated_at, version)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`, string(t.ID), t.Name, t.Slug, t.OwnerUserID, string(t.Status),
		t.CreatedAt, t.UpdatedAt, t.Version)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return tenant.ErrSlugTaken
		}
		return fmt.Errorf("tenant_repo: create: %w", err)
	}
	return nil
}

func (r *TenantRepo) Update(ctx context.Context, t *tenant.Tenant) error {
	tag, err := r.pool.Exec(ctx, `
		UPDATE tenants
		   SET name = $2,
		       status = $3,
		       updated_at = $4,
		       version = $5
		 WHERE id = $1 AND version = $6
	`, string(t.ID), t.Name, string(t.Status), t.UpdatedAt, t.Version, t.Version-1)
	if err != nil {
		return fmt.Errorf("tenant_repo: update: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return tenant.ErrVersionConflict
	}
	return nil
}

func (r *TenantRepo) GetByID(ctx context.Context, id shared.TenantID) (*tenant.Tenant, error) {
	var t tenant.Tenant
	var status string
	err := r.pool.QueryRow(ctx, `
		SELECT id, name, slug, owner_user_id, status, created_at, updated_at, version
		  FROM tenants
		 WHERE id = $1
	`, string(id)).Scan(
		(*string)(&t.ID), &t.Name, &t.Slug, &t.OwnerUserID,
		&status, &t.CreatedAt, &t.UpdatedAt, &t.Version,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, tenant.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("tenant_repo: get: %w", err)
	}
	t.Status = tenant.Status(status)
	return &t, nil
}

func (r *TenantRepo) SlugExists(ctx context.Context, slug string) (bool, error) {
	var exists bool
	err := r.pool.QueryRow(ctx, `
		SELECT EXISTS(SELECT 1 FROM tenants WHERE slug = $1)
	`, slug).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("tenant_repo: slug exists: %w", err)
	}
	return exists, nil
}
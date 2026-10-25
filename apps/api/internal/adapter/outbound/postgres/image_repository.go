package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"uvico-api/internal/domain/image"
	"uvico-api/internal/domain/shared"
	"uvico-api/internal/platform/database"
)

type ImageRepo struct {
	pool *database.Pool
}

func NewImageRepo(pool *database.Pool) *ImageRepo { return &ImageRepo{pool: pool} }

func (r *ImageRepo) Insert(ctx context.Context, img *image.Image) error {
	return r.pool.WithTenantTx(ctx, img.TenantID.String(), func(ctx context.Context, tx pgx.Tx) error {
		_, err := tx.Exec(ctx, `
			INSERT INTO product_images
				(id, tenant_id, product_id, storage_path, alt_text, position, created_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7)
		`, img.ID, img.TenantID.String(), img.ProductID.String(),
			img.StoragePath, img.AltText, img.Position, img.CreatedAt)
		if err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && pgErr.Code == "23505" {
				return fmt.Errorf("%w: storage_path duplicado", image.ErrInvalidPath)
			}
			return fmt.Errorf("image_repo: insert: %w", err)
		}
		return nil
	})
}

func (r *ImageRepo) ListByProduct(ctx context.Context, tenantID shared.TenantID, productID shared.ProductID) ([]*image.Image, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, tenant_id, product_id, storage_path, alt_text, position, created_at
		  FROM product_images
		 WHERE tenant_id = $1 AND product_id = $2
		 ORDER BY position ASC, created_at ASC
	`, tenantID.String(), productID.String())
	if err != nil {
		return nil, fmt.Errorf("image_repo: list: %w", err)
	}
	defer rows.Close()

	var out []*image.Image
	for rows.Next() {
		img, err := scanImage(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, img)
	}
	return out, rows.Err()
}

func (r *ImageRepo) GetByID(ctx context.Context, tenantID shared.TenantID, productID shared.ProductID, imageID string) (*image.Image, error) {
	row := r.pool.QueryRow(ctx, `
		SELECT id, tenant_id, product_id, storage_path, alt_text, position, created_at
		  FROM product_images
		 WHERE tenant_id = $1 AND product_id = $2 AND id = $3
	`, tenantID.String(), productID.String(), imageID)
	return scanImage(row)
}

func (r *ImageRepo) Delete(ctx context.Context, tenantID shared.TenantID, productID shared.ProductID, imageID string) error {
	return r.pool.WithTenantTx(ctx, tenantID.String(), func(ctx context.Context, tx pgx.Tx) error {
		tag, err := tx.Exec(ctx, `
			DELETE FROM product_images
			 WHERE tenant_id = $1 AND product_id = $2 AND id = $3
		`, tenantID.String(), productID.String(), imageID)
		if err != nil {
			return fmt.Errorf("image_repo: delete: %w", err)
		}
		if tag.RowsAffected() == 0 {
			return image.ErrNotFound
		}
		return nil
	})
}

func (r *ImageRepo) Reorder(ctx context.Context, tenantID shared.TenantID, productID shared.ProductID, order []string) error {
	return r.pool.WithTenantTx(ctx, tenantID.String(), func(ctx context.Context, tx pgx.Tx) error {
		for pos, id := range order {
			tag, err := tx.Exec(ctx, `
				UPDATE product_images
				   SET position = $4
				 WHERE tenant_id = $1 AND product_id = $2 AND id = $3
			`, tenantID.String(), productID.String(), id, pos)
			if err != nil {
				return fmt.Errorf("image_repo: reorder %s: %w", id, err)
			}
			if tag.RowsAffected() == 0 {
				return fmt.Errorf("%w: %s", image.ErrNotFound, id)
			}
		}
		return nil
	})
}

func (r *ImageRepo) CountByProduct(ctx context.Context, tenantID shared.TenantID, productID shared.ProductID) (int, error) {
	var count int
	err := r.pool.QueryRow(ctx, `
		SELECT COUNT(*) FROM product_images
		 WHERE tenant_id = $1 AND product_id = $2
	`, tenantID.String(), productID.String()).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("image_repo: count: %w", err)
	}
	return count, nil
}

func scanImage(row rowScanner) (*image.Image, error) {
	var img image.Image
	var id, tid, pid string
	err := row.Scan(&id, &tid, &pid, &img.StoragePath, &img.AltText, &img.Position, &img.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, image.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("image_repo: scan: %w", err)
	}
	img.ID = id
	img.TenantID = shared.TenantID(tid)
	img.ProductID = shared.ProductID(pid)
	return &img, nil
}
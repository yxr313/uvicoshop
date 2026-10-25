package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"

	"uvico-api/internal/domain/inventory"
	"uvico-api/internal/domain/shared"
	"uvico-api/internal/platform/database"
)

type InventoryRepo struct {
	pool *database.Pool
}

func NewInventoryRepo(pool *database.Pool) *InventoryRepo { return &InventoryRepo{pool: pool} }

func (r *InventoryRepo) Reserve(ctx context.Context, tenantID shared.TenantID, productID shared.ProductID, qty int64, referenceID string) error {
	if qty <= 0 {
		return inventory.ErrInvalidQuantity
	}
	return r.pool.WithTenantTx(ctx, tenantID.String(), func(ctx context.Context, tx pgx.Tx) error {
		tag, err := tx.Exec(ctx, `
			UPDATE inventory_items
			   SET reserved   = reserved + $1,
			       version    = version + 1,
			       updated_at = NOW()
			 WHERE tenant_id  = $2
			   AND product_id = $3
			   AND (on_hand - reserved) >= $1
		`, qty, tenantID.String(), productID.String())
		if err != nil {
			return fmt.Errorf("inventory_repo: reserve: %w", err)
		}
		if tag.RowsAffected() == 0 {
			return inventory.ErrInsufficientStock
		}
		return insertMovement(ctx, tx, tenantID, productID, -qty, "reserved", referenceID)
	})
}

func (r *InventoryRepo) Commit(ctx context.Context, tenantID shared.TenantID, productID shared.ProductID, qty int64, referenceID string) error {
	if qty <= 0 {
		return inventory.ErrInvalidQuantity
	}
	return r.pool.WithTenantTx(ctx, tenantID.String(), func(ctx context.Context, tx pgx.Tx) error {
		tag, err := tx.Exec(ctx, `
			UPDATE inventory_items
			   SET on_hand    = on_hand - $1,
			       reserved   = reserved - $1,
			       version    = version + 1,
			       updated_at = NOW()
			 WHERE tenant_id  = $2
			   AND product_id = $3
			   AND reserved >= $1
		`, qty, tenantID.String(), productID.String())
		if err != nil {
			return fmt.Errorf("inventory_repo: commit: %w", err)
		}
		if tag.RowsAffected() == 0 {
			return inventory.ErrNotFound
		}
		return insertMovement(ctx, tx, tenantID, productID, 0, "committed", referenceID)
	})
}

func (r *InventoryRepo) Release(ctx context.Context, tenantID shared.TenantID, productID shared.ProductID, qty int64, referenceID string) error {
	if qty <= 0 {
		return inventory.ErrInvalidQuantity
	}
	return r.pool.WithTenantTx(ctx, tenantID.String(), func(ctx context.Context, tx pgx.Tx) error {
		tag, err := tx.Exec(ctx, `
			UPDATE inventory_items
			   SET reserved   = reserved - $1,
			       version    = version + 1,
			       updated_at = NOW()
			 WHERE tenant_id  = $2
			   AND product_id = $3
			   AND reserved >= $1
		`, qty, tenantID.String(), productID.String())
		if err != nil {
			return fmt.Errorf("inventory_repo: release: %w", err)
		}
		if tag.RowsAffected() == 0 {
			return inventory.ErrNotFound
		}
		return insertMovement(ctx, tx, tenantID, productID, qty, "released", referenceID)
	})
}

func (r *InventoryRepo) Adjust(ctx context.Context, tenantID shared.TenantID, productID shared.ProductID, delta int64) error {
	if delta == 0 {
		return inventory.ErrInvalidQuantity
	}
	return r.pool.WithTenantTx(ctx, tenantID.String(), func(ctx context.Context, tx pgx.Tx) error {
		tag, err := tx.Exec(ctx, `
			UPDATE inventory_items
			   SET on_hand    = on_hand + $1,
			       version    = version + 1,
			       updated_at = NOW()
			 WHERE tenant_id  = $2
			   AND product_id = $3
			   AND (on_hand + $1) >= reserved
			   AND (on_hand + $1) >= 0
		`, delta, tenantID.String(), productID.String())
		if err != nil {
			return fmt.Errorf("inventory_repo: adjust: %w", err)
		}
		if tag.RowsAffected() == 0 {
			return inventory.ErrNotFound
		}
		return insertMovement(ctx, tx, tenantID, productID, delta, "adjustment", "")
	})
}

func (r *InventoryRepo) GetByProduct(ctx context.Context, tenantID shared.TenantID, productID shared.ProductID) (*inventory.Item, error) {
	var it inventory.Item
	var tid, pid string
	err := r.pool.QueryRow(ctx, `
		SELECT tenant_id, product_id, on_hand, reserved, version, created_at, updated_at
		  FROM inventory_items
		 WHERE tenant_id = $1 AND product_id = $2
	`, tenantID.String(), productID.String()).Scan(
		&tid, &pid, &it.OnHand, &it.Reserved, &it.Version, &it.CreatedAt, &it.UpdatedAt,
	)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return nil, inventory.ErrNotFound
		}
		return nil, fmt.Errorf("inventory_repo: get: %w", err)
	}
	it.TenantID = shared.TenantID(tid)
	it.ProductID = shared.ProductID(pid)
	return &it, nil
}

func insertMovement(ctx context.Context, tx pgx.Tx, tenantID shared.TenantID, productID shared.ProductID, delta int64, reason, referenceID string) error {
	var ref any
	if referenceID != "" {
		ref = referenceID
	}
	_, err := tx.Exec(ctx, `
		INSERT INTO inventory_movements
			(id, tenant_id, product_id, delta, reason, reference_id, created_at)
		VALUES (gen_random_uuid(), $1, $2, $3, $4, $5, $6)
	`, tenantID.String(), productID.String(), delta, reason, ref, time.Now().UTC())
	if err != nil {
		return fmt.Errorf("inventory_repo: movement: %w", err)
	}
	return nil
}

// ReleaseTx es la variante que corre dentro de una tx externa.
// Usada por el webhook service cuando libera stock tras un pago fallido.
func (r *InventoryRepo) ReleaseTx(
	ctx context.Context,
	tx pgx.Tx,
	tenantID shared.TenantID,
	productID shared.ProductID,
	qty int64,
	referenceID string,
) error {
	if qty <= 0 {
		return inventory.ErrInvalidQuantity
	}
	tag, err := tx.Exec(ctx, `
		UPDATE inventory_items
		SET reserved   = reserved - $1,
		    version    = version + 1,
		    updated_at = NOW()
		WHERE tenant_id  = $2
		AND product_id = $3
		AND reserved >= $1
	`, qty, tenantID.String(), productID.String())
	if err != nil {
		return fmt.Errorf("inventory_repo: release: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return inventory.ErrNotFound
	}
	return insertMovement(ctx, tx, tenantID, productID, qty, "released", referenceID)
}
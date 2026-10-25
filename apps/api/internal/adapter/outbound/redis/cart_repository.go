package redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	goredis "github.com/redis/go-redis/v9"

	"uvico-api/internal/domain/cart"
	"uvico-api/internal/domain/shared"
)

const cartTTL = 15 * time.Minute

type CartRepo struct {
	client *goredis.Client
}

func NewCartRepo(client *goredis.Client) *CartRepo {
	return &CartRepo{client: client}
}

func (r *CartRepo) key(tenantID shared.TenantID, sessionID string) string {
	return "cart:" + tenantID.String() + ":" + sessionID
}

func (r *CartRepo) Get(ctx context.Context, tenantID shared.TenantID, sessionID string) (*cart.Cart, error) {
	raw, err := r.client.Get(ctx, r.key(tenantID, sessionID)).Bytes()
	if errors.Is(err, goredis.Nil) {
		return &cart.Cart{
			TenantID:  tenantID,
			SessionID: sessionID,
			Items:     []cart.Item{},
			UpdatedAt: time.Now().UTC(),
		}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("cart: get: %w", err)
	}

	var c cart.Cart
	if err := json.Unmarshal(raw, &c); err != nil {
		return nil, fmt.Errorf("cart: unmarshal: %w", err)
	}
	c.RecalculateSubtotal()
	return &c, nil
}

func (r *CartRepo) AddItem(ctx context.Context, tenantID shared.TenantID, sessionID string, item cart.Item) (*cart.Cart, error) {
	// Leer, modificar, escribir con WATCH para atomicidad optimista.
	key := r.key(tenantID, sessionID)

	txf := func(tx *goredis.Tx) error {
		raw, err := tx.Get(ctx, key).Bytes()
		if err != nil && !errors.Is(err, goredis.Nil) {
			return err
		}

		var c cart.Cart
		if len(raw) > 0 {
			if err := json.Unmarshal(raw, &c); err != nil {
				return err
			}
		}
		c.TenantID = tenantID
		c.SessionID = sessionID
		if c.Items == nil {
			c.Items = []cart.Item{}
		}

		found := false
		for i := range c.Items {
			if c.Items[i].ProductID == item.ProductID {
				c.Items[i].Quantity += item.Quantity
				found = true
				break
			}
		}
		if !found {
			if len(c.Items) >= cart.MaxItemsPerCart {
				return fmt.Errorf("cart: excede %d items", cart.MaxItemsPerCart)
			}
			c.Items = append(c.Items, item)
		}
		c.UpdatedAt = time.Now().UTC()
		c.RecalculateSubtotal()

		payload, err := json.Marshal(c)
		if err != nil {
			return err
		}
		_, err = tx.TxPipelined(ctx, func(pipe goredis.Pipeliner) error {
			pipe.Set(ctx, key, payload, cartTTL)
			return nil
		})
		return err
	}

	if err := r.client.Watch(ctx, txf, key); err != nil {
		return nil, fmt.Errorf("cart: add: %w", err)
	}
	return r.Get(ctx, tenantID, sessionID)
}

func (r *CartRepo) SetQuantity(ctx context.Context, tenantID shared.TenantID, sessionID string, productID shared.ProductID, qty int32) (*cart.Cart, error) {
	key := r.key(tenantID, sessionID)

	txf := func(tx *goredis.Tx) error {
		raw, err := tx.Get(ctx, key).Bytes()
		if err != nil {
			return err
		}
		var c cart.Cart
		if err := json.Unmarshal(raw, &c); err != nil {
			return err
		}
		found := false
		for i := range c.Items {
			if c.Items[i].ProductID == productID {
				c.Items[i].Quantity = qty
				found = true
				break
			}
		}
		if !found {
			return cart.ErrItemNotFound
		}
		c.UpdatedAt = time.Now().UTC()
		c.RecalculateSubtotal()
		payload, err := json.Marshal(c)
		if err != nil {
			return err
		}
		return tx.Set(ctx, key, payload, cartTTL).Err()
	}

	if err := r.client.Watch(ctx, txf, key); err != nil {
		return nil, err
	}
	return r.Get(ctx, tenantID, sessionID)
}

func (r *CartRepo) RemoveItem(ctx context.Context, tenantID shared.TenantID, sessionID string, productID shared.ProductID) (*cart.Cart, error) {
	key := r.key(tenantID, sessionID)

	txf := func(tx *goredis.Tx) error {
		raw, err := tx.Get(ctx, key).Bytes()
		if err != nil {
			return err
		}
		var c cart.Cart
		if err := json.Unmarshal(raw, &c); err != nil {
			return err
		}
		found := false
		filtered := make([]cart.Item, 0, len(c.Items))
		for _, it := range c.Items {
			if it.ProductID == productID {
				found = true
				continue
			}
			filtered = append(filtered, it)
		}
		if !found {
			return cart.ErrItemNotFound
		}
		c.Items = filtered
		c.UpdatedAt = time.Now().UTC()
		c.RecalculateSubtotal()
		payload, err := json.Marshal(c)
		if err != nil {
			return err
		}
		return tx.Set(ctx, key, payload, cartTTL).Err()
	}

	if err := r.client.Watch(ctx, txf, key); err != nil {
		return nil, err
	}
	return r.Get(ctx, tenantID, sessionID)
}

func (r *CartRepo) Clear(ctx context.Context, tenantID shared.TenantID, sessionID string) error {
	empty := cart.Cart{
		TenantID:  tenantID,
		SessionID: sessionID,
		Items:     []cart.Item{},
		UpdatedAt: time.Now().UTC(),
	}
	payload, err := json.Marshal(empty)
	if err != nil {
		return err
	}
	return r.client.Set(ctx, r.key(tenantID, sessionID), payload, cartTTL).Err()
}

func (r *CartRepo) Delete(ctx context.Context, tenantID shared.TenantID, sessionID string) error {
	return r.client.Del(ctx, r.key(tenantID, sessionID)).Err()
}
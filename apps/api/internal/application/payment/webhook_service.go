package payment

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5"

	"uvico-api/internal/domain/inventory"
	orderdomain "uvico-api/internal/domain/order"
	"uvico-api/internal/platform/database"
)

type WebhookRepository interface {
	InsertIfNewTx(ctx context.Context, tx pgx.Tx, eventID, provider, orderID string) error
}

type WebhookService struct {
	pool      *database.Pool
	orderRepo orderdomain.TxRepository
	invRepo   inventory.TxRepository
	payRepo   WebhookRepository
	logger    *slog.Logger
	now       func() time.Time
}

func NewWebhookService(
	pool *database.Pool,
	orderRepo orderdomain.TxRepository,
	invRepo inventory.TxRepository,
	payRepo WebhookRepository,
	logger *slog.Logger,
) *WebhookService {
	return &WebhookService{
		pool: pool, orderRepo: orderRepo, invRepo: invRepo, payRepo: payRepo,
		logger: logger, now: time.Now,
	}
}

type ProcessResult struct {
	AlreadyProcessed bool
	OrderID          string
	Applied          string
}

func (s *WebhookService) Process(
	ctx context.Context,
	provider, eventID, paymentRef string,
	targetStatus orderdomain.Status,
) (*ProcessResult, error) {
	if eventID == "" {
		return nil, errors.New("payment: event_id requerido")
	}
	if paymentRef == "" {
		return nil, errors.New("payment: payment_ref requerido")
	}

	// 1. Encontrar la orden para saber el TenantID
	existing, err := s.orderRepo.GetByPaymentRefForUpdate(ctx, paymentRef)
	if err != nil {
		return nil, err
	}

	result := &ProcessResult{OrderID: existing.ID, Applied: string(targetStatus)}
	now := s.now()

	// 2. Abrir transacción con el tenant correcto
	err = s.pool.WithTenantTx(ctx, existing.TenantID.String(), func(ctx context.Context, tx pgx.Tx) error {
		o, err := s.orderRepo.GetByPaymentRefForUpdateTx(ctx, tx, paymentRef)
		if err != nil {
			return err
		}

		prevStatus := o.Status
		var applyErr error
		switch targetStatus {
		case orderdomain.StatusPaid:
			applyErr = o.MarkPaid(paymentRef, now)
		case orderdomain.StatusFailed:
			applyErr = o.MarkFailed(now)
		case orderdomain.StatusCancelled:
			applyErr = o.Cancel(now)
		default:
			return fmt.Errorf("payment: target_status no soportado %q", targetStatus)
		}
		if applyErr != nil {
			return applyErr
		}

		if o.Status == prevStatus {
			if err := s.payRepo.InsertIfNewTx(ctx, tx, eventID, provider, o.ID); err != nil {
				return err
			}
			result.AlreadyProcessed = true
			return nil
		}

		// 3. Actualizar orden
		if err := s.orderRepo.UpdateTx(ctx, tx, o); err != nil {
			return err
		}

		// 4. Liberar stock si falló/canceló
		if targetStatus == orderdomain.StatusFailed || targetStatus == orderdomain.StatusCancelled {
			for _, item := range o.Items {
				if err := s.invRepo.ReleaseTx(ctx, tx, o.TenantID, item.ProductID, int64(item.Quantity), o.ID); err != nil {
					return err
				}
			}
		}

		// 5. Registrar idempotencia
		if err := s.payRepo.InsertIfNewTx(ctx, tx, eventID, provider, o.ID); err != nil {
			return err
		}

		// 6. Outbox
		eventType := "payment.received"
		if targetStatus != orderdomain.StatusPaid {
			eventType = "payment.failed"
		}
		raw, _ := json.Marshal(map[string]any{
			"order_id": o.ID,
			"status":   string(o.Status),
			"provider": provider,
			"event_id": eventID,
		})

		_, err = tx.Exec(ctx, `
			INSERT INTO outbox_events
				(id, tenant_id, aggregate_type, aggregate_id, event_type, payload, created_at)
			VALUES (gen_random_uuid(), $1, 'order', $2, $3, $4, NOW())
		`, o.TenantID.String(), o.ID, eventType, string(raw))
		
		return err
	})

	if err != nil {
		return nil, err
	}
	return result, nil
}
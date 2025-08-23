package paymentusecases

import (
	"context"
	"errors"
	"net/http"
	"time"

	redisdb "github.com/guatom999/ecommerce-payment-api/databases/redisdb"
	"github.com/guatom999/ecommerce-payment-api/modules"
	paymentrepositories "github.com/guatom999/ecommerce-payment-api/modules/paymentRepositories"
)

type (
	paymentUsecase struct {
		paymentRepo paymentrepositories.PaymentRepositoryInterface
		redis       *redisdb.Store
	}
)

func NewPaymentUsecase(paymentRepo paymentrepositories.PaymentRepositoryInterface) PaymentUsecaseInterface {
	return &paymentUsecase{
		paymentRepo: paymentRepo,
	}
}

func (u *paymentUsecase) CreatePayment(ctx context.Context, cmd modules.CreatePaymentCommand) (*modules.PaymentRow, error) {
	if cmd.OrderID == "" || cmd.Amount < 0 || cmd.Currency == "" {
		return nil, errors.New("invalid input")
	}

	if cmd.IdempotencyKey != "" && u.redis != nil {
		started, rec, err := u.redis.TryStart(ctx, cmd.IdempotencyKey, 10*time.Minute)
		if err == nil && !started {
			if rec != nil && rec.Status == "done" {
				if pid, ok := rec.Response["id"].(string); ok && pid != "" {
					return u.paymentRepo.Get(ctx, pid)
				}
			}
			return nil, errors.New("request in progress")
		}
	}

	paymentID, err := u.paymentRepo.CreateProcessing(ctx, modules.CreatePaymentRequest{
		OrderID:  cmd.OrderID,
		UserID:   cmd.UserID,
		Amount:   cmd.Amount,
		Currency: cmd.Currency,
	})
	if err != nil {
		return nil, err
	}

	p, err := u.paymentRepo.Get(ctx, paymentID)
	if err != nil {
		return nil, err
	}

	if cmd.IdempotencyKey != "" && u.redis != nil {
		resp := map[string]any{
			"id":         p.ID,
			"order_id":   p.OrderID,
			"user_id":    p.UserID,
			"amount":     p.Amount,
			"currency":   p.Currency,
			"status":     p.Status,
			"created_at": p.CreatedAt,
			"updated_at": p.UpdatedAt,
		}
		_ = u.redis.Complete(ctx, cmd.IdempotencyKey, p.ID, http.StatusOK, resp, 24*time.Hour)
	}

	return p, nil
}
func (u *paymentUsecase) GetPayment(ctx context.Context, id string) (*modules.PaymentRow, error) {
	return u.paymentRepo.Get(ctx, id)
}
func (u *paymentUsecase) UpdateStatus(ctx context.Context, id, status string) error {
	return u.paymentRepo.UpdateStatus(ctx, id, status)
}

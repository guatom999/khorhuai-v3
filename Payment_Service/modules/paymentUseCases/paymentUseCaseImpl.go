package paymentusecases

import (
	"context"
	"errors"
	"net/http"
	"time"

	redisdb "github.com/guatom999/ecommerce-payment-api/databases/redisDb"
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
				// คุณอาจแปลง rec.Response → PaymentRow ได้ แต่ในที่นี้ให้ไปอ่าน DB ตาม payment_id จะชัวร์กว่า
				if pid, ok := rec.Response["id"].(string); ok && pid != "" {
					return u.paymentRepo.Get(ctx, pid)
				}
			}
			// ยัง processing อยู่
			return nil, errors.New("request in progress")
		}
		// ถ้า Redis error หรือ started=true → ดำเนินการต่อ (soft lock)
	}

	// --- สร้าง payment (processing)
	paymentID, err := u.paymentRepo.CreateProcessing(ctx, modules.CreatePaymentInput{
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

	// --- เก็บผลลัพธ์ลง Redis เพื่อ fast repeat
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

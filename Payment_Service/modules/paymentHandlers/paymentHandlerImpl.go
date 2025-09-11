package paymenthandlers

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/guatom999/ecommerce-payment-api/modules"
	"github.com/guatom999/ecommerce-payment-api/modules/paymentusecases"
	"github.com/guatom999/ecommerce-payment-api/request"
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
	"go.temporal.io/sdk/client"
)

type (
	paymenthandler struct {
		paymentUsease paymentusecases.PaymentUsecaseInterface
	}

	paymentWebhookHandler struct {
		Temporal client.Client
		Rdb      *redis.Client
	}
)

func NewPaymenthandler(paymentUsease paymentusecases.PaymentUsecaseInterface) PaymenthandlerInterface {
	return &paymenthandler{paymentUsease: paymentUsease}
}

func NewWebhookHandler(Temporal client.Client, rdb *redis.Client) PaymentWebhookHandlerInterface {
	return &paymentWebhookHandler{Temporal: Temporal, Rdb: rdb}
}

func (h *paymenthandler) CreatePayment(c echo.Context) error {

	ctx := context.Background()

	req := new(modules.CreatePaymentRequest)

	wrapper := request.NewContextWrapper(c)

	if err := wrapper.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid input"})
	}

	paymentCommand := modules.CreatePaymentCommand{
		OrderID:        req.OrderID,
		UserID:         req.UserID,
		Amount:         req.Amount,
		Currency:       req.Currency,
		IdempotencyKey: c.Request().Header.Get("Idempotency-Key"),
	}

	result, err := h.paymentUsease.CreatePayment(ctx, paymentCommand)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, result)
}
func (h *paymenthandler) GetPayment(c echo.Context) error {

	ctx := context.Background()

	id := c.Param("id")

	result, err := h.paymentUsease.GetPayment(ctx, id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, result)
}
func (h *paymenthandler) UpdatePaymentStatus(c echo.Context) error {

	ctx := context.Background()

	id := c.Param("id")

	wrapper := request.NewContextWrapper(c)

	req := new(modules.UpdatePaymentStatusRequest)

	if err := wrapper.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid input"})
	}

	if err := h.paymentUsease.UpdateStatus(ctx, id, req.Status); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "status updated"})
}

func (h *paymentWebhookHandler) PaymentWebhook(c echo.Context) error {

	ctx := context.Background()

	req := new(modules.WebhookEvent)

	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid input"})
	}

	req.Status = strings.ToLower(strings.TrimSpace(req.Status))
	if req.OrderId == "" || (req.Status != "succeeded" && req.Status != "failed") {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "missing inpuit"})
	}

	idemKey := h.buildIdemKey(req)
	ok, err := h.rdbSetOnce(ctx, idemKey, time.Duration(time.Minute*1))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "redis error" + err.Error()})
	}

	if !ok {
		return c.JSON(http.StatusOK, echo.Map{"ok": true, "dup": true})
	}

	workflowId := "checkout_" + req.OrderId
	payload := map[string]any{
		"status":     req.Status,
		"payment_id": req.PaymentId,
	}

	if err := retry(3, time.Second, func() error {
		return h.Temporal.SignalWorkflow(c.Request().Context(), workflowId, "", "payment_signal", payload)
	}); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "signal failed: " + err.Error()})
	}

	return c.JSON(http.StatusOK, echo.Map{"ok": true})
}

func (h *paymentWebhookHandler) buildIdemKey(req *modules.WebhookEvent) string {
	if req.EventId != "" {
		return fmt.Sprintf("webhook:payment:%s", req.EventId)
	}

	return fmt.Sprintf("webhook:payment:%s:%s:%s", req.OrderId, req.PaymentId, req.Status)
}

func (h *paymentWebhookHandler) rdbSetOnce(ctx context.Context, key string, ttl time.Duration) (bool, error) {
	return h.Rdb.SetNX(ctx, key, "1", ttl).Result()
}

func retry(n int, base time.Duration, fn func() error) error {
	var err error
	for i := 0; i < n; i++ {
		if err = fn(); err == nil {
			return nil
		}
		time.Sleep(base * time.Duration(1<<i)) // 1x,2x,4x
	}
	return errors.New("after retries: " + err.Error())
}

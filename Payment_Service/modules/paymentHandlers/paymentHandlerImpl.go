package paymenthandlers

import (
	"context"
	"net/http"

	"github.com/guatom999/ecommerce-payment-api/modules"
	paymentusecases "github.com/guatom999/ecommerce-payment-api/modules/paymentUseCases"
	"github.com/guatom999/ecommerce-payment-api/request"
	"github.com/labstack/echo/v4"
)

type (
	paymenthandler struct {
		paymentUsease paymentusecases.PaymentUsecaseInterface
	}
)

func NewPaymenthandler(paymentUsease paymentusecases.PaymentUsecaseInterface) PaymenthandlerInterface {
	return &paymenthandler{paymentUsease: paymentUsease}
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

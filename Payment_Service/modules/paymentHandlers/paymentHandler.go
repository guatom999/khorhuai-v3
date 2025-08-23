package paymenthandlers

import "github.com/labstack/echo/v4"

type (
	PaymenthandlerInterface interface {
		CreatePayment(c echo.Context) error
		GetPayment(c echo.Context) error
		UpdatePaymentStatus(c echo.Context) error
	}
)

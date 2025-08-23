package orderHandlers

import "github.com/labstack/echo/v4"

type (
	OrderHandlerInterface interface {
		CreateOrder(c echo.Context) error
		GetOrder(c echo.Context) error
		ListOrdersByUser(c echo.Context) error
		UpdateStatus(c echo.Context) error
	}
)

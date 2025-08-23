package cartHandlers

import "github.com/labstack/echo/v4"

type (
	CartHandlerInterface interface {
		GetCart(c echo.Context) error
		UpsertItem(c echo.Context) error
		Merge(c echo.Context) error
	}
)

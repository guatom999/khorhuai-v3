package producthandlers

import "github.com/labstack/echo/v4"

type (
	ProducthandlerInterface interface {
		Reserve(c echo.Context) error
		Release(c echo.Context) error
		Commit(c echo.Context) error
	}
)

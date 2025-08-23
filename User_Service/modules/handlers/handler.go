package handlers

import "github.com/labstack/echo/v4"

type (
	UserHandlerInterface interface {
		Register(c echo.Context) error
		Login(c echo.Context) error
		EditUser(c echo.Context) error
	}
)

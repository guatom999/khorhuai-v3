package notihandlers

import "github.com/labstack/echo/v4"

type (
	NotiHandlerInterface interface {
		Create(c echo.Context) error
		Get(c echo.Context) error
		List(c echo.Context) error
		AttemptSend(c echo.Context) error
		UpdateStatus(c echo.Context) error
	}
)

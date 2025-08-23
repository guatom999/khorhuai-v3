package request

import (
	"log"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
)

type (
	ContextWrapperInterface interface {
		Bind(data any) error
	}

	contextWrapper struct {
		echoCtx   echo.Context
		validator *validator.Validate
	}
)

func NewContextWrapper(echoCtx echo.Context) ContextWrapperInterface {
	return &contextWrapper{
		echoCtx:   echoCtx,
		validator: validator.New(),
	}
}

func (c *contextWrapper) Bind(data any) error {
	if err := c.echoCtx.Bind(data); err != nil {
		log.Printf("Error: Bind data failed: %v", err)
	}

	if err := c.validator.Struct(data); err != nil {
		log.Printf("Error: Validate failed: %v", err)
	}

	return nil
}

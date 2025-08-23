package producthandlers

import "github.com/labstack/echo/v4"

type (
	ProductHandlerInterface interface {
		GetAllProductWithCategory(c echo.Context) error
		GetAllProduct(c echo.Context) error
		CreateProduct(c echo.Context) error
		UpdateProduct(c echo.Context) error
		DeleteProduct(c echo.Context) error
	}
)

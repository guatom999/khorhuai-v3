package producthandlers

import (
	"context"
	"errors"
	"net/http"

	"github.com/guatom999/ecommerce-product-api/modules"
	productusecases "github.com/guatom999/ecommerce-product-api/modules/productUsecases"
	"github.com/guatom999/ecommerce-product-api/utils/request"
	"github.com/labstack/echo/v4"
)

type (
	productHandler struct {
		productUsecase productusecases.ProductUsecaseInterface
	}
)

func NewProductHandler(productUsecase productusecases.ProductUsecaseInterface) ProductHandlerInterface {
	return &productHandler{
		productUsecase: productUsecase,
	}
}

func (h *productHandler) GetAllProductWithCategory(c echo.Context) error {

	ctx := context.Background()

	products, err := h.productUsecase.GetAllProductWithCategory(ctx)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, products)

}

func (h *productHandler) GetAllProduct(c echo.Context) error {

	ctx := context.Background()

	products, err := h.productUsecase.GetAllProduct(ctx)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, products)
}

func (h *productHandler) CreateProduct(c echo.Context) error {

	ctx := context.Background()

	wrapper := request.NewContextWrapper(c)

	req := new(modules.CreateProductCategoryRequest)

	if err := wrapper.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	if err := h.productUsecase.CreateProduct(ctx, req); err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, "created success")

}

func (h *productHandler) UpdateProduct(c echo.Context) error {

	ctx := context.Background()

	wrapper := request.NewContextWrapper(c)

	req := new(modules.ProductPatchReq)

	if err := wrapper.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, errors.New("bad request"))
	}

	if err := h.productUsecase.UpdateProduct(ctx, req.Id, req); err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, "updated success")
}

func (h *productHandler) DeleteProduct(c echo.Context) error {

	ctx := context.Background()

	wrapper := request.NewContextWrapper(c)

	req := new(modules.DeleteProductReq)

	if err := wrapper.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, errors.New("bad request"))
	}

	if err := h.productUsecase.DeleteProduct(ctx, req.Id); err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, "deleted success")

}

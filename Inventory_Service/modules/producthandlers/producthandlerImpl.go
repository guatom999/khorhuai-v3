package producthandlers

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/guatom999/ecommerce-product-api/modules"
	"github.com/guatom999/ecommerce-product-api/modules/productusecases"
	"github.com/labstack/echo/v4"
)

type (
	producthandler struct {
		productUsecase productusecases.ProductusecaseInterface
	}
)

func NewProductHandler(productUsecase productusecases.ProductusecaseInterface) ProducthandlerInterface {
	return &producthandler{productUsecase: productUsecase}
}

func (h *producthandler) Reserve(c echo.Context) error {

	ctx := context.Background()

	req := new(modules.ReserveStockRequest)

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid body"})
	}

	ttl := time.Duration(req.TTLSeconds) * time.Second
	cmd := modules.ReserveInput{OrderID: req.OrderID, UserID: req.UserID, TTL: ttl}

	for _, item := range req.Items {
		cmd.Items = append(cmd.Items, modules.ReserveItem{ProductID: item.ProductID, Quantity: item.Quantity})
	}

	reserveId, err := h.productUsecase.Reserve(ctx, cmd)
	if err != nil {
		log.Printf("error reserving stock: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, echo.Map{"reservation_id": reserveId})
}
func (h *producthandler) Release(c echo.Context) error {

	ctx := context.Background()

	req := new(modules.ReleaseStockRequest)

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid body"})
	}

	if err := h.productUsecase.Release(ctx, req.ReservationID); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, echo.Map{"status": "released"})
}
func (h *producthandler) Commit(c echo.Context) error {

	ctx := context.Background()

	req := new(modules.CommitStockRequest)

	if err := c.Bind(&req); err != nil || req.ReservationID == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid body"})
	}
	if err := h.productUsecase.Commit(ctx, req.ReservationID); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, echo.Map{"status": "committed"})
}

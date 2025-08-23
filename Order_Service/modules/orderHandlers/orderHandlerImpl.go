package orderHandlers

import (
	"context"
	"net/http"
	"strconv"

	"github.com/guatom999/ecommerce-order-api/modules"
	"github.com/guatom999/ecommerce-order-api/modules/orderUsecases"
	"github.com/guatom999/ecommerce-order-api/utils/request"
	"github.com/labstack/echo/v4"
)

type (
	orderHandler struct {
		orderUsecase orderUsecases.OrderUsecaseInterface
	}
)

func NewOrderHandler(orderUsecase orderUsecases.OrderUsecaseInterface) OrderHandlerInterface {
	return &orderHandler{orderUsecase: orderUsecase}
}

func (h *orderHandler) CreateOrder(c echo.Context) error {

	ctx := context.Background()

	req := new(modules.CreateOrderRequest)

	wrapper := request.NewContextWrapper(c)

	if err := wrapper.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid body"})
	}

	if len(req.Items) == 0 {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "items required"})
	}

	input := modules.CreateOrderInput{
		UserID:       req.UserID,
		Currency:     req.Currency,
		Items:        make([]modules.OrderItemInput, 0, len(req.Items)),
		ShippingAddr: req.BillingAddr,
		BillingAddr:  req.BillingAddr,
	}

	for _, it := range req.Items {
		input.Items = append(input.Items, modules.OrderItemInput{
			ProductID: it.ProductID,
			Title:     it.Title,
			SKU:       it.SKU,
			Quantity:  it.Quantity,
			UnitPrice: it.UnitPrice,
		})
	}

	orderId, err := h.orderUsecase.Create(ctx, input)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()}) // TODO: change error message
	}

	return c.JSON(http.StatusOK, echo.Map{"order_id": orderId})
}

func (h *orderHandler) GetOrder(c echo.Context) error {

	ctx := context.Background()

	orderId := c.Param("id")

	order, err := h.orderUsecase.Get(ctx, orderId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()}) // TODO: change error message
	}

	return c.JSON(http.StatusOK, echo.Map{"order": order})
}

func (h *orderHandler) ListOrdersByUser(c echo.Context) error {

	ctx := context.Background()

	userID := c.Param("user_id")
	limit := queryInt(c, "limit", 20)
	offset := queryInt(c, "offset", 0)

	rows, err := h.orderUsecase.ListByUser(ctx, userID, limit, offset)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()}) // TODO: change error message
	}

	return c.JSON(http.StatusOK, rows)
}

func (h *orderHandler) UpdateStatus(c echo.Context) error {

	ctx := context.Background()

	id := c.Param("id")

	req := new(modules.UpdateStatusRequest)

	wrapper := request.NewContextWrapper(c)

	if err := wrapper.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid body"})
	}

	if err := h.orderUsecase.UpdateStatus(ctx, id, req.Status); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, echo.Map{"status": "updated"})
}

func queryInt(c echo.Context, key string, def int) int {
	if v := c.QueryParam(key); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return def
}

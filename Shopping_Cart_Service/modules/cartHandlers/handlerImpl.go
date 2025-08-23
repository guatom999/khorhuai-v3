package cartHandlers

import (
	"context"
	"net/http"

	"github.com/guatom999/ecommerce-shopping-cart-api/modules"
	"github.com/guatom999/ecommerce-shopping-cart-api/modules/cartUsecases"
	"github.com/guatom999/ecommerce-shopping-cart-api/utils/request"
	"github.com/labstack/echo/v4"
)

type (
	cartHandler struct {
		cartUsecase cartUsecases.CartUsecaseInterface
	}
)

func NewCartHandler(cartUsecase cartUsecases.CartUsecaseInterface) CartHandlerInterface {
	return &cartHandler{cartUsecase: cartUsecase}
}

func (h *cartHandler) GetCart(c echo.Context) error {

	ctx := context.Background()

	userID := headerPtr(c, "X-User-ID")
	sessionID := headerPtr(c, "X-Session-ID")

	cv, err := h.cartUsecase.GetCart(ctx, userID, sessionID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, cv)
}
func (h *cartHandler) UpsertItem(c echo.Context) error {

	ctx := context.Background()

	wrapper := request.NewContextWrapper(c)

	req := new(modules.UpsertCartItemRequest)
	if err := wrapper.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, "bad request")
	}
	if req.Quantity <= 0 || req.Currency == "" {
		return c.JSON(http.StatusBadRequest, "bad request")
	}

	userID := headerPtr(c, "X-User-ID")
	sessionID := headerPtr(c, "X-Session-ID")
	if userID == "" && sessionID == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "require X-User-ID or X-Session-ID"})
	}

	cv, err := h.cartUsecase.UpsertItem(ctx, userID, sessionID, req.ProductID, req.Quantity, req.UnitPriceCents, req.Currency)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, cv)
}
func (h *cartHandler) Merge(c echo.Context) error {

	ctx := context.Background()

	wrapper := request.NewContextWrapper(c)

	req := new(modules.MergeCartRequest)

	if err := wrapper.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid body"})
	}
	if req.UserID == "" || req.SessionID == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "user_id and session_id required"})
	}

	if err := h.cartUsecase.Merge(ctx, req.UserID, req.SessionID); err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, echo.Map{"status": "merged"})
}

func headerPtr(c echo.Context, key string) string {
	v := c.Request().Header.Get(key)
	if v == "" {
		return ""
	}
	return v
}

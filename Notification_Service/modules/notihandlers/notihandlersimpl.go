package notihandlers

import (
	"context"
	"net/http"
	"strconv"

	"github.com/guatom999/ecommerce-notification-api/modules"
	"github.com/guatom999/ecommerce-notification-api/modules/notiusecases"
	"github.com/guatom999/ecommerce-notification-api/request"
	"github.com/labstack/echo/v4"
)

type (
	notiHandler struct {
		notiUsecase notiusecases.NotiusecaseInterface
	}
)

func NewNotiHandler(notiUsecase notiusecases.NotiusecaseInterface) NotiHandlerInterface {
	return &notiHandler{notiUsecase: notiUsecase}
}

func (h *notiHandler) Create(c echo.Context) error {

	ctx := context.Background()

	req := new(modules.CreateNotificationRequest)

	wrapper := request.NewContextWrapper(c)

	if err := wrapper.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "bad request"})
	}

	id, err := h.notiUsecase.Create(ctx, modules.CreateInput{
		UserID:       req.UserID,
		Channel:      req.Channel,
		Recipient:    req.Recipient,
		TemplateName: req.TemplateName,
		Data:         req.Data,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "internal server error"})
	}

	return c.JSON(http.StatusCreated, echo.Map{"notification_id": id})
}
func (h *notiHandler) Get(c echo.Context) error {

	ctx := context.Background()

	userId := c.Param("user_id")

	result, err := h.notiUsecase.Get(ctx, userId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "internal server error"})
	}

	return c.JSON(http.StatusOK, result)
}
func (h *notiHandler) List(c echo.Context) error {

	ctx := context.Background()

	userId := c.QueryParam("user_id")
	status := c.QueryParam("status")
	limit := func() int {
		limit, err := strconv.Atoi(c.QueryParam("limit"))
		if err != nil {
			return 20
		}
		return limit
	}()
	offset := func() int {
		offset, err := strconv.Atoi(c.QueryParam("offset"))
		if err != nil {
			return 0
		}
		return offset
	}()

	result, err := h.notiUsecase.List(ctx, userId, status, limit, offset)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, result)
}
func (h *notiHandler) AttemptSend(c echo.Context) error {

	ctx := context.Background()

	id := c.Param("id")

	wrapper := request.NewContextWrapper(c)

	req := new(modules.SendAttemptRequest)

	if err := wrapper.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid body"})
	}

	if err := h.notiUsecase.AttemptSend(ctx, id, req.Status, req.ErrorMessage, req.ProviderRaw); err != nil {
		c.JSON(http.StatusOK, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "updated"})
}
func (h *notiHandler) UpdateStatus(c echo.Context) error {

	ctx := context.Background()

	id := c.Param("id")

	wrapper := request.NewContextWrapper(c)

	req := new(modules.UpdateNotificationStatusRequest)

	if err := wrapper.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid body"})
	}

	if err := h.notiUsecase.UpdateStatus(ctx, id, req.Status); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "updated"})
}

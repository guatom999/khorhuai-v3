package handlers

import (
	"context"
	"errors"
	"net/http"

	"github.com/guatom999/ecommerce-user-api/modules"
	"github.com/guatom999/ecommerce-user-api/modules/usecases"
	"github.com/labstack/echo/v4"
)

var (
	ErrDuplicateEmail = errors.New("email already exists")
)

type (
	userHandler struct {
		userUsecase usecases.UserUsecaseInterface
	}
)

func NewUserHandler(userUsecase usecases.UserUsecaseInterface) UserHandlerInterface {
	return &userHandler{userUsecase: userUsecase}
}

func (h *userHandler) Register(c echo.Context) error {

	ctx := context.Background()

	req := new(modules.CreateUserReq)

	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, "invalid request")
	}

	if err := h.userUsecase.Register(ctx, req); err != nil {
		if errors.Is(err, ErrDuplicateEmail) {
			return c.JSON(http.StatusBadRequest, "user already exists")
		}

		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusAccepted, "user created")
}

func (h *userHandler) Login(c echo.Context) error {

	ctx := context.Background()

	req := new(modules.LoginReq)

	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, "invalid request")
	}

	user, err := h.userUsecase.Login(ctx, req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "something weng wrong")
	}

	return c.JSON(http.StatusAccepted, user)

}

func (h *userHandler) EditUser(c echo.Context) error {

	ctx := context.Background()

	req := new(modules.EditUserReq)

	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, "invalid request")
	}

	if err := h.userUsecase.EditUser(ctx, req); err != nil {
		return c.JSON(http.StatusInternalServerError, "something went wrong with server")
	}

	return c.JSON(http.StatusOK, "edit user success")
}

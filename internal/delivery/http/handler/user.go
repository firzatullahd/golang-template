package handler

import (
	"net/http"

	"github.com/firzatullahd/cats-social-api/internal/model"
	"github.com/labstack/echo/v4"
)

func (h *Handler) Register(c echo.Context) error {
	ctx := c.Request().Context()
	var payload model.RegisterRequest
	if err := c.Bind(&payload); err != nil {
		return err
	}
	resp, err := h.Usecase.Register(ctx, &payload)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusCreated, resp)
}

func (h *Handler) Login(c echo.Context) error {
	ctx := c.Request().Context()
	var payload model.LoginRequest
	if err := c.Bind(&payload); err != nil {
		return err
	}
	resp, err := h.Usecase.Login(ctx, &payload)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, resp)
}

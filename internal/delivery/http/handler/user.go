package handler

import (
	"net/http"

	"github.com/firzatullahd/cats-social-api/internal/model"
	error_envelope "github.com/firzatullahd/cats-social-api/internal/model/error"
	"github.com/labstack/echo/v4"
)

func (h *Handler) Register(c echo.Context) error {
	ctx := c.Request().Context()
	var payload model.RegisterRequest
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, model.ErrorResponse{Message: err.Error()})
	}
	resp, err := h.Usecase.Register(ctx, &payload)
	if err != nil {
		code, errMsg := error_envelope.ParseError(err)
		return c.JSON(code, model.ErrorResponse{Message: errMsg})
	}
	return c.JSON(http.StatusCreated, model.Response[*model.AuthResponse]{Data: resp, Message: "User registered successfully"})
}

func (h *Handler) Login(c echo.Context) error {
	ctx := c.Request().Context()
	var payload model.LoginRequest
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, model.ErrorResponse{Message: err.Error()})
	}
	resp, err := h.Usecase.Login(ctx, &payload)
	if err != nil {
		code, errMsg := error_envelope.ParseError(err)
		return c.JSON(code, model.ErrorResponse{Message: errMsg})
	}
	return c.JSON(http.StatusOK, model.Response[*model.AuthResponse]{Data: resp, Message: "User logged successfully"})
}

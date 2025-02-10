package handler

import (
	"net/http"

	"github.com/firzatullahd/golang-template/internal/user/model"
	customerror "github.com/firzatullahd/golang-template/internal/user/model/error"
	"github.com/firzatullahd/golang-template/utils/constant"
	"github.com/firzatullahd/golang-template/utils/logger"
	"github.com/firzatullahd/golang-template/utils/response"
	"github.com/labstack/echo/v4"
)

func (h *Handler) Register(c echo.Context) error {
	ctx := c.Request().Context()
	var payload model.RegisterRequest
	if err := c.Bind(&payload); err != nil {
		return response.ErrorResponse(c, http.StatusBadRequest, err)
	}
	resp, err := h.Service.Register(ctx, payload)
	if err != nil {
		code, err := customerror.ParseError(err)
		logger.Log.Error(err)
		return response.ErrorResponse(c, code, err)
	}
	return response.SuccessResponse(c, http.StatusCreated, resp, nil)
}

func (h *Handler) Login(c echo.Context) error {
	ctx := c.Request().Context()
	var payload model.AuthRequest
	if err := c.Bind(&payload); err != nil {
		return response.ErrorResponse(c, http.StatusBadRequest, err)
	}
	resp, err := h.Service.Login(ctx, payload)
	if err != nil {
		code, err := customerror.ParseError(err)
		return response.ErrorResponse(c, code, err)
	}

	return response.SuccessResponse(c, http.StatusOK, resp, nil)
}

func (h *Handler) InitiateVerification(c echo.Context) error {
	ctx := c.Request().Context()
	data, ok := c.Get(constant.UserDataKey).(model.UserData)
	if !ok {
		return response.ErrorResponse(c, http.StatusUnauthorized, nil)
	}

	err := h.Service.InitiateVerification(ctx, data.Username)
	if err != nil {
		code, err := customerror.ParseError(err)
		return response.ErrorResponse(c, code, err)
	}

	return response.SuccessResponse(c, http.StatusOK, nil, nil)
}

func (h *Handler) Verification(c echo.Context) error {
	ctx := c.Request().Context()
	data, ok := c.Get(constant.UserDataKey).(model.UserData)
	if !ok {
		return response.ErrorResponse(c, http.StatusUnauthorized, nil)
	}
	code := c.Param("code")
	err := h.Service.Verification(ctx, data.Username, code)
	if err != nil {
		code, err := customerror.ParseError(err)
		return response.ErrorResponse(c, code, err)
	}

	return response.SuccessResponse(c, http.StatusOK, nil, nil)
}

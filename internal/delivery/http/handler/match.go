package handler

import (
	"net/http"
	"strconv"

	"github.com/firzatullahd/cats-social-api/internal/model"
	error_envelope "github.com/firzatullahd/cats-social-api/internal/model/error"
	"github.com/firzatullahd/cats-social-api/internal/utils/constant"
	"github.com/labstack/echo/v4"
)

func (h *Handler) CreateMatch(c echo.Context) error {
	ctx := c.Request().Context()

	data, ok := c.Get(constant.UserDataKey).(model.UserData)
	if !ok {
		return c.JSON(http.StatusUnauthorized, model.ErrorResponse{Message: error_envelope.ErrUnauthorized.Error()})
	}

	var payload model.CreateMatchRequest
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, model.ErrorResponse{Message: err.Error()})
	}

	payload.UserId = data.ID
	err := h.Usecase.CreateMatch(ctx, &payload)
	if err != nil {
		code, errMsg := error_envelope.ParseError(err)
		return c.JSON(code, model.ErrorResponse{Message: errMsg})
	}

	return c.JSON(http.StatusCreated, model.Response[*model.GeneralResponse]{Message: "success"})
}

func (h *Handler) DeleteMatch(c echo.Context) error {
	ctx := c.Request().Context()

	data, ok := c.Get(constant.UserDataKey).(model.UserData)
	if !ok {
		return c.JSON(http.StatusUnauthorized, model.ErrorResponse{Message: error_envelope.ErrUnauthorized.Error()})
	}

	id := c.Param("id")
	matchId, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, model.ErrorResponse{Message: err.Error()})
	}

	err = h.Usecase.DeleteMatch(ctx, matchId, data.ID)
	if err != nil {
		code, errMsg := error_envelope.ParseError(err)
		return c.JSON(code, model.ErrorResponse{Message: errMsg})
	}

	return c.JSON(http.StatusCreated, model.Response[*model.GeneralResponse]{Message: "successfully delete match"})
}

func (h *Handler) ApproveMatch(c echo.Context) error {
	ctx := c.Request().Context()

	data, ok := c.Get(constant.UserDataKey).(model.UserData)
	if !ok {
		return c.JSON(http.StatusUnauthorized, model.ErrorResponse{Message: error_envelope.ErrUnauthorized.Error()})
	}

	var payload model.UpdateMatchRequest
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, model.ErrorResponse{Message: err.Error()})
	}

	err := h.Usecase.ApproveMatch(ctx, payload.MatchId, data.ID)
	if err != nil {
		code, errMsg := error_envelope.ParseError(err)
		return c.JSON(code, model.ErrorResponse{Message: errMsg})
	}

	return c.JSON(http.StatusCreated, model.Response[*model.GeneralResponse]{Message: "successfully approve match"})
}

func (h *Handler) RejectMatch(c echo.Context) error {
	ctx := c.Request().Context()

	data, ok := c.Get(constant.UserDataKey).(model.UserData)
	if !ok {
		return c.JSON(http.StatusUnauthorized, model.ErrorResponse{Message: error_envelope.ErrUnauthorized.Error()})
	}

	var payload model.UpdateMatchRequest
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, model.ErrorResponse{Message: err.Error()})
	}

	err := h.Usecase.RejectMatch(ctx, payload.MatchId, data.ID)
	if err != nil {
		code, errMsg := error_envelope.ParseError(err)
		return c.JSON(code, model.ErrorResponse{Message: errMsg})
	}

	return c.JSON(http.StatusCreated, model.Response[*model.GeneralResponse]{Message: "successfully reject match"})
}

func (h *Handler) FindMatch(c echo.Context) error {
	ctx := c.Request().Context()

	data, ok := c.Get(constant.UserDataKey).(model.UserData)
	if !ok {
		return c.JSON(http.StatusUnauthorized, model.ErrorResponse{Message: error_envelope.ErrUnauthorized.Error()})
	}

	var payload model.UpdateMatchRequest
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, model.ErrorResponse{Message: err.Error()})
	}

	resp, err := h.Usecase.FindMatch(ctx, data.ID)
	if err != nil {
		code, errMsg := error_envelope.ParseError(err)
		return c.JSON(code, model.ErrorResponse{Message: errMsg})
	}

	return c.JSON(http.StatusCreated, model.Response[[]model.FindMatchResponse]{Data: resp, Message: "successfully reject match"})
}

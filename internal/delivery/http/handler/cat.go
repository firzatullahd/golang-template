package handler

import (
	"net/http"
	"strconv"

	"github.com/firzatullahd/cats-social-api/internal/model"
	error_envelope "github.com/firzatullahd/cats-social-api/internal/model/error"
	"github.com/firzatullahd/cats-social-api/internal/utils/constant"
	"github.com/labstack/echo/v4"
)

func (h *Handler) CreateCat(c echo.Context) error {
	ctx := c.Request().Context()

	data, ok := c.Get(constant.UserDataKey).(model.UserData)
	if !ok {
		return c.JSON(http.StatusUnauthorized, model.ErrorResponse{Message: error_envelope.ErrUnauthorized.Error()})
	}

	var payload model.CreateCatRequest
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, model.ErrorResponse{Message: err.Error()})
	}

	resp, err := h.Usecase.CreateCat(ctx, &payload, data.ID)
	if err != nil {
		code, errMsg := error_envelope.ParseError(err)
		return c.JSON(code, model.ErrorResponse{Message: errMsg})
	}

	return c.JSON(http.StatusCreated, model.Response[*model.CreateCatResponse]{Data: resp, Message: "success"})
}

func (h *Handler) DeleteCat(c echo.Context) error {
	ctx := c.Request().Context()

	data, ok := c.Get(constant.UserDataKey).(model.UserData)
	if !ok {
		return c.JSON(http.StatusUnauthorized, model.ErrorResponse{Message: error_envelope.ErrUnauthorized.Error()})
	}

	id := c.Param("id")
	catId, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, model.ErrorResponse{Message: err.Error()})
	}

	err = h.Usecase.DeleteCat(ctx, catId, data.ID)
	if err != nil {
		code, errMsg := error_envelope.ParseError(err)
		return c.JSON(code, model.ErrorResponse{Message: errMsg})
	}

	return c.JSON(http.StatusCreated, model.Response[*model.GeneralResponse]{Message: "successfully delete cat"})
}

func (h *Handler) UpdateCat(c echo.Context) error {
	ctx := c.Request().Context()

	data, ok := c.Get(constant.UserDataKey).(model.UserData)
	if !ok {
		return c.JSON(http.StatusUnauthorized, model.ErrorResponse{Message: error_envelope.ErrUnauthorized.Error()})
	}

	id := c.Param("id")
	catId, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, model.ErrorResponse{Message: err.Error()})
	}

	var payload model.UpdateCatRequest
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, model.ErrorResponse{Message: err.Error()})
	}

	payload.ID = catId
	payload.UserID = data.ID

	err = h.Usecase.UpdateCat(ctx, &payload)
	if err != nil {
		code, errMsg := error_envelope.ParseError(err)
		return c.JSON(code, model.ErrorResponse{Message: errMsg})
	}

	return c.JSON(http.StatusCreated, model.Response[*model.GeneralResponse]{Message: "successfully update cat"})
}

func (h *Handler) FindCat(c echo.Context) error {
	ctx := c.Request().Context()

	data, ok := c.Get(constant.UserDataKey).(model.UserData)
	if !ok {
		return c.JSON(http.StatusUnauthorized, model.ErrorResponse{Message: error_envelope.ErrUnauthorized.Error()})
	}

	resp, err := h.Usecase.FindCat(ctx, &model.FindCatRequest{
		Limit:      c.Param("limit"),
		Offset:     c.Param("offset"),
		ID:         c.Param("id"),
		Sex:        c.Param("sex"),
		Race:       c.Param("race"),
		HasMatched: c.Param("hasMatched"),
		Age:        c.Param("ageInMonth"),
		SearchName: c.Param("search"),
		Owned:      c.Param("owned"),
		UserId:     data.ID,
	})
	if err != nil {
		code, errMsg := error_envelope.ParseError(err)
		return c.JSON(code, model.ErrorResponse{Message: errMsg})
	}

	return c.JSON(http.StatusCreated, model.Response[[]model.FindCatResponse]{Data: resp, Message: "successfully find cat"})
}

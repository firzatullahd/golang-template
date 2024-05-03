package handler

import (
	"net/http"

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

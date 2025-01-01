package response

import (
	"strconv"

	echo "github.com/labstack/echo/v4"
)

type (
	Response struct {
		StatusCode int         `json:"status_code"`
		Message    string      `json:"message"`
		Data       interface{} `json:"data"`
		Pagination *Pagination `json:"pagination,omitempty"`
		Error      *Error      `json:"error,omitempty"`
	}
	Error struct {
		Code   string        `json:"code"`
		Errors []interface{} `json:"details,omitempty"`
	}
	Pagination struct {
		CurrentPage   uint64 `json:"current_page"`
		LastPage      uint64 `json:"last_page"`
		Count         uint64 `json:"count"`
		RecordPerPage uint64 `json:"record_per_page"`
	}
)

func SuccessResponse(c echo.Context, statusCode int, data interface{}, pagination *Pagination) error {
	return c.JSON(statusCode, &Response{
		Data:       data,
		StatusCode: statusCode,
		Message:    "success",
		Pagination: pagination,
	})
}

func ErrorResponse(c echo.Context, statusCode int, err error, details ...interface{}) error {
	var (
		message = err.Error()
	)

	// ctx := e.Request().Context()
	resp := Response{
		StatusCode: statusCode,
		Message:    message,
		Error: &Error{
			Code:   strconv.Itoa(statusCode),
			Errors: details,
		},
	}
	// log.Error(ctx, "got response with error", log.Any("response", resp), log.Any("error", *resp.Error))

	return c.JSON(statusCode, resp)
}

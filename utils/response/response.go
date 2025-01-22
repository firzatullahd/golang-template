package response

import (
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
		Code   int           `json:"code"`
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
	return c.JSON(statusCode, Response{
		StatusCode: statusCode,
		Message:    err.Error(),
	})
}

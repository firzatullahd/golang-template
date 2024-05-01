package middleware

import (
	"context"

	"github.com/google/uuid"
	echo "github.com/labstack/echo/v4"
)

const (
	CorrelationIDKey string = "X-Correlation-ID"
)

func LogContext() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ctx := context.WithValue(c.Request().Context(), CorrelationIDKey, uuid.New().String())
			req := c.Request().WithContext(ctx)
			c.SetRequest(req)
			return next(c)
		}
	}
}

package route

import (
	"net/http"

	"github.com/firzatullahd/golang-template/config"
	"github.com/firzatullahd/golang-template/internal/user/delivery/http/handler"
	"github.com/firzatullahd/golang-template/internal/user/delivery/http/middleware"

	echo "github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
)

func Serve(conf *config.Config, h *handler.Handler) {
	e := echo.New()
	m := middleware.NewMiddleware(conf.JWTSecretKey)
	e.Pre(echoMiddleware.RemoveTrailingSlash())
	// e.Use(echoMiddleware.Logger(), echoMiddleware.Recover())
	e.Use(m.LogContext())

	e.GET("/health", HealthCheck)

	userApi := e.Group("/v1/user")
	userApi.POST("/register", h.Register)
	userApi.POST("/login", h.Login)
	userApi.POST("/verification/:username", h.InitialVerification)
	userApi.POST("/verify/:username/:code", h.Verify)

	_ = e.Group("/v1/", m.Auth())

	// approval verification

	e.Logger.Fatal(e.Start(":" + conf.AppPort))
}

func HealthCheck(c echo.Context) error {
	return c.String(http.StatusOK, "Welcome to golang-templates")
}

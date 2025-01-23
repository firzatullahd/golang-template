package route

import (
	"net/http"
	"time"

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
	e.Use(echoMiddleware.Logger(), echoMiddleware.Recover())
	e.Use(m.LogContext())

	e.GET("/health", HealthCheck)

	userApi := e.Group("/v1/user")
	userApi.POST("/register", h.Register)
	userApi.POST("/login", h.Login)

	authApi := e.Group("/v1/user", m.Auth)

	authApi.POST("/verification", h.InitialVerification)
	authApi.POST("/verification/:code", h.Verification)

	// approval verification
	server := &http.Server{
		Addr:         ":" + conf.AppPort,
		ReadTimeout:  3600 * time.Second,
		WriteTimeout: 3600 * time.Second,
	}
	e.Logger.Fatal(e.StartServer(server))
}

func HealthCheck(c echo.Context) error {
	return c.String(http.StatusOK, "Welcome to golang-templates")
}

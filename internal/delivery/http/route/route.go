package route

import (
	"fmt"
	"net/http"

	"github.com/firzatullahd/cats-social-api/internal/config"
	"github.com/firzatullahd/cats-social-api/internal/delivery/http/handler"
	"github.com/firzatullahd/cats-social-api/internal/delivery/http/middleware"

	echo "github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
)

func Serve(conf *config.Config, h *handler.Handler) {
	e := echo.New()
	e.Pre(echoMiddleware.AddTrailingSlash())
	e.Use(echoMiddleware.Logger(), echoMiddleware.Recover())
	e.Use(middleware.LogContext())

	e.GET("/", hello)

	userApi := e.Group("/v1/user")
	userApi.POST("/register", h.Register)
	userApi.POST("/login", h.Login)

	// todo
	catApi := e.Group("/v1/cat", middleware.Auth())
	catApi.POST("/", h.Login)
	catApi.GET("/", h.Login)
	catApi.PUT("/:id", h.Login)
	catApi.DELETE("/:id", h.Login)

	// todo
	matchApi := e.Group("/v1/match", middleware.Auth())
	matchApi.POST("/", h.Login)
	matchApi.GET("/", h.Login)
	matchApi.POST("/approve", h.Login)
	matchApi.POST("/reject", h.Login)
	matchApi.DELETE("/:id", h.Login)

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", conf.Port)))
}

func hello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}

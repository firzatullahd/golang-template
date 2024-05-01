package route

import (
	"fmt"
	"net/http"

	"github.com/firzatullahd/cats-social-api/internal/config"
	"github.com/firzatullahd/cats-social-api/internal/delivery/http/handler"
	"github.com/labstack/echo"
)

func Serve(conf *config.Config, h *handler.Handler) {
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	e.POST("/user/register", h.Register)
	e.POST("/user/login", h.Login)

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", conf.Port)))
}

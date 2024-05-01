package route

import (
	"net/http"

	"github.com/firzatullahd/cats-social-api/internal/config"
	"github.com/labstack/echo"
)

func Start(conf *config.Config) {
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.Logger.Fatal(e.Start(":1323"))
}

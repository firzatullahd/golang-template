package middleware

import (
	"errors"
	"strings"
	"time"

	"github.com/firzatullahd/golang-template/internal/user/model"
	"github.com/firzatullahd/golang-template/utils/constant"
	"github.com/golang-jwt/jwt/v5"
	echo "github.com/labstack/echo/v4"
)

func (m *Middleware) Auth() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			accessToken := c.Request().Header.Get("Authorization")
			if strings.Contains(accessToken, "Bearer") {
				accessToken = strings.Replace(accessToken, "Bearer ", "", -1)
			}

			claims := model.MyClaim{}
			token, err := jwt.ParseWithClaims(accessToken, &claims, func(token *jwt.Token) (interface{}, error) {
				return []byte(m.JWTSecretKey), nil
			})

			if err != nil {
				return errors.New("token invalid")
			}

			if !token.Valid {
				return errors.New("token invalid")
			}

			timeExp, err := token.Claims.GetExpirationTime()
			if err != nil {
				return errors.New("token invalid")
			}

			res := timeExp.Compare(time.Now())
			if res == -1 {
				return errors.New("token expired")
			}

			c.Set(constant.UserDataKey, claims.UserData)
			return nil
		}
	}
}

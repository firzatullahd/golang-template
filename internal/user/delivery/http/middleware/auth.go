package middleware

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/firzatullahd/golang-template/internal/user/model"
	"github.com/firzatullahd/golang-template/utils/constant"
	"github.com/firzatullahd/golang-template/utils/response"
	"github.com/golang-jwt/jwt/v5"
	echo "github.com/labstack/echo/v4"
)

func (m *Middleware) Auth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		accessToken := c.Request().Header.Get("Authorization")
		if strings.Contains(accessToken, "Bearer") {
			accessToken = strings.Replace(accessToken, "Bearer ", "", -1)
		}

		claims := model.MyClaim{}
		token, err := jwt.ParseWithClaims(accessToken, &claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(m.JWTSecretKey), nil
		})

		if err != nil || !token.Valid {
			return response.ErrorResponse(c, http.StatusUnauthorized, errors.New("token invalid"))
		}

		timeExp, err := token.Claims.GetExpirationTime()
		if err != nil {
			return response.ErrorResponse(c, http.StatusUnauthorized, errors.New("token invalid"))
		}

		res := timeExp.Compare(time.Now())
		if res == -1 {
			return response.ErrorResponse(c, http.StatusUnauthorized, errors.New("token expired"))
		}

		c.Set(constant.UserDataKey, claims.UserData)
		return next(c)
	}
}

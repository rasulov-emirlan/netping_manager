package rest

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

func GetToken(accessToken string, jwtkey []byte) (*Claims, error) {
	token, err := jwt.ParseWithClaims(accessToken, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return jwtkey, nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("users: Invalid Access Token")
}

const UserInfoFromContext = "userinfo"

func CheckRole(jwtKey []byte, isAdmin bool) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")

			if authHeader == "" {
				return c.NoContent(http.StatusUnauthorized)
			}

			authHeaders := strings.Split(authHeader, " ")
			if len(authHeaders) != 2 {
				return c.NoContent(http.StatusUnauthorized)
			}

			if authHeaders[0] != "Bearer" {
				return c.NoContent(http.StatusUnauthorized)
			}

			claims, err := GetToken(authHeaders[1], jwtKey)
			if err != nil {
				return c.NoContent(http.StatusUnauthorized)
			}

			if claims.IsAdmin != isAdmin {
				return c.NoContent(http.StatusUnauthorized)
			}

			c.Set(UserInfoFromContext, claims)
			return next(c)
		}
	}
}

package rest

import (
	"errors"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

func GetToken(c echo.Context, jwtkey []byte) (*Claims, error) {
	cookie, err := c.Cookie("AccessToken")
	if err != nil {
		return nil, err
	}
	tknStr := cookie.Value
	claims := &Claims{}

	tkn, err := jwt.ParseWithClaims(tknStr, claims, func(t *jwt.Token) (interface{}, error) {
		return jwtkey, nil
	})
	if err != nil {
		return nil, err
	}
	if !tkn.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}

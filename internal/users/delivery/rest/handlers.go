package rest

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/rasulov-emirlan/netping-manager/internal/users"
)

const jwtTokenExpirationTime = time.Hour * 8

type Claims struct {
	UserID  int  `json:"userId"`
	IsAdmin bool `json:"isAdmin"`
	jwt.StandardClaims
}

type handler struct {
	service users.Service
	jwtKey  []byte
}

func NewHandler(s users.Service, jwtKey []byte) (*handler, error) {
	if s == nil {
		return nil, errors.New("users: delivery: rest: service cannot be nil")
	}
	return &handler{
		service: s,
		jwtKey:  jwtKey,
	}, nil
}

func (h *handler) Register(router *echo.Group) error {
	router.GET("/config/users", h.getUsers(), CheckRole(h.jwtKey, true))
	router.POST("/config/users", h.registerUser(), CheckRole(h.jwtKey, true))
	router.DELETE("/config/users/:id", h.deleteUser(), CheckRole(h.jwtKey, true))
	router.POST("/config/users/login", h.login())
	router.POST("/config/users/logout", h.logout())
	return nil
}

func (h *handler) getUsers() echo.HandlerFunc {
	return func(c echo.Context) error {
		u, err := h.service.ReadAll(c.Request().Context())
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, u)
	}
}

func (h *handler) login() echo.HandlerFunc {
	type Request struct {
		Name     string `json:"username"`
		Password string `json:"password"`
	}
	type Response struct {
		AccessToken string `json:"accessToken"`
	}
	return func(c echo.Context) error {
		req := &Request{}
		if err := c.Bind(req); err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		u, err := h.service.ReadByName(c.Request().Context(), req.Name)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		if u.Password != req.Password {
			return c.NoContent(http.StatusUnauthorized)
		}

		expTime := time.Now().Add(jwtTokenExpirationTime)

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, &Claims{
			UserID:  u.ID,
			IsAdmin: u.IsAdmin,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: expTime.Unix(),
			},
		})
		tokenString, err := token.SignedString(h.jwtKey)
		if err != nil {
			return c.NoContent(http.StatusInternalServerError)
		}

		c.SetCookie(&http.Cookie{
			Name:     "AccessToken",
			Value:    tokenString,
			Expires:  expTime,
			HttpOnly: true,
		})
		return c.JSON(http.StatusOK, Response{AccessToken: tokenString})
	}
}

func (h *handler) logout() echo.HandlerFunc {
	return func(c echo.Context) error {
		c.SetCookie(&http.Cookie{
			Name:     "AccessToken",
			Value:    "",
			Expires:  time.Now(),
			HttpOnly: true,
		})
		return c.NoContent(http.StatusOK)
	}
}

func (h *handler) registerUser() echo.HandlerFunc {
	type Request struct {
		Name     string `json:"username"`
		Password string `json:"password"`
	}
	return func(c echo.Context) error {
		req := &Request{}
		if err := c.Bind(req); err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		u, err := h.service.Create(
			c.Request().Context(),
			req.Name, req.Password,
			false,
		)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, u)
	}
}

func (h *handler) deleteUser() echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")
		userID, err := strconv.Atoi(id)
		if err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		if err := h.service.Delete(
			c.Request().Context(),
			userID,
		); err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		return c.NoContent(http.StatusOK)
	}
}

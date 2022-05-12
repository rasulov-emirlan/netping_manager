package rest

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/rasulov-emirlan/netping-manager/internal/users"
	"go.uber.org/zap"
)

const (
	jwtTokenExpirationTime     = time.Minute * 20
	refreshTokenExpirationTime = time.Hour * 24
	refreshCookieName          = "RefreshToken"
)

type Claims struct {
	UserID  int  `json:"userId"`
	IsAdmin bool `json:"isAdmin"`
	jwt.StandardClaims
}

type RefreshClaims struct {
	UserID int `json:"userid"`
	jwt.StandardClaims
}

type handler struct {
	service users.Service
	log     *zap.SugaredLogger
	jwtKey  []byte
}

func NewHandler(s users.Service, jwtKey []byte, log *zap.SugaredLogger) (*handler, error) {
	if s == nil {
		return nil, errors.New("users: delivery: rest: service cannot be nil")
	}
	return &handler{
		service: s,
		jwtKey:  jwtKey,
		log:     log,
	}, nil
}

func (h *handler) Register(router *echo.Group) error {
	router.GET("/config/users", h.getUsers(), CheckRole(h.jwtKey, true))
	router.POST("/config/users", h.registerUser(), CheckRole(h.jwtKey, true))
	router.PATCH("/config/users/:id", h.update(), CheckRole(h.jwtKey, true))
	router.DELETE("/config/users/:id", h.deleteUser(), CheckRole(h.jwtKey, true))
	router.POST("/config/users/login", h.login())
	router.POST("/config/users/refresh", h.refresh())
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

func (h *handler) refresh() echo.HandlerFunc {
	type Response struct {
		AccessToken  string `json:"accessToken"`
		RefreshToken string `json:"refreshToken"`
	}
	return func(c echo.Context) error {
		// We validate that user has a cookie with our refresh token
		cookie, err := c.Cookie(refreshCookieName)
		if err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		if err := cookie.Valid(); err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		claims, err := GetRefresh(cookie.Value, h.jwtKey)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}

		// Then we give the user a new access token
		// If user has been banned by this time he will not get a new access token
		// cause we will get an error from our service ;)
		u, err := h.service.Read(c.Request().Context(), claims.UserID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}

		// Here we generate the actual access token
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

		// And just return it to the user
		return c.JSON(http.StatusOK, Response{AccessToken: tokenString, RefreshToken: cookie.Value})
	}
}

func (h *handler) login() echo.HandlerFunc {
	type Request struct {
		Name     string `json:"username"`
		Password string `json:"password"`
	}
	type Response struct {
		AccessToken  string `json:"accessToken"`
		RefreshToken string `json:"refreshToken"`
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
		refreshTime := time.Now().Add(refreshTokenExpirationTime)

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

		refresh := jwt.NewWithClaims(jwt.SigningMethodHS256, &RefreshClaims{
			UserID: u.ID,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: refreshTime.Unix(),
			},
		})

		refreshString, err := refresh.SignedString(h.jwtKey)
		if err != nil {
			return c.NoContent(http.StatusInternalServerError)
		}

		c.SetCookie(&http.Cookie{
			Name:     refreshCookieName,
			Value:    refreshString,
			Expires:  expTime,
			HttpOnly: true,
		})
		return c.JSON(http.StatusOK, Response{AccessToken: tokenString, RefreshToken: refreshString})
	}
}

func (h *handler) logout() echo.HandlerFunc {
	return func(c echo.Context) error {
		c.SetCookie(&http.Cookie{
			Name:     refreshCookieName,
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

func (h *handler) update() echo.HandlerFunc {
	type Request struct {
		Name     string `json:"name"`
		Passwrod string `json:"password"`
		IsAdmin  bool   `json:"isAdmin"`
	}
	return func(c echo.Context) error {
		claims, ok := c.Get(UserInfoFromContext).(*Claims)
		if !ok {
			return c.JSON(http.StatusUnauthorized, "could not cast Claims to access token interface{}")
		}

		// We will log in this handler so we will have
		// to sync our logger at the end of this function
		defer h.log.Sync()
		userID, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}

		// Here we take json from request body
		// And convert it onto go struct Request{}
		req := Request{}
		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		u := users.User{
			Name:     req.Name,
			Password: req.Passwrod,
			IsAdmin:  req.IsAdmin,
		}

		if err := h.service.Update(c.Request().Context(), userID, u); err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		h.log.Infow("Users state in database has been updated", zap.Int("updater-userid", claims.UserID), zap.Int("updated-userid", userID))
		return c.NoContent(http.StatusOK)
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

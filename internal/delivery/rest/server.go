package rest

import (
	"context"
	"embed"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Registrator interface {
	Register(router *echo.Group) error
}

type server struct {
	port   string
	router *echo.Echo

	websiteFS          *embed.FS
	managerRegistrator Registrator
	usersRegistrator   Registrator
}

func NewServer(port string, websiteFS *embed.FS, tw, tr time.Duration, m, u Registrator) (*server, error) {
	e := echo.New()
	e.Server.ReadTimeout = tr
	e.Server.WriteTimeout = tw
	return &server{
		router:             e,
		port:               port,
		managerRegistrator: m,
		usersRegistrator:   u,
		websiteFS:          websiteFS,
	}, nil
}

func (s *server) Start() error {
	s.router.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		// This application will not be used outside of inner network
		// so we can allow all origins
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET", "POST", "DELETE", "PUT", "PATCH", "OPTIONS"},
		// We will have to save some cookies so we have to allow credentials
		AllowCredentials: true,
	}))

	s.router.Group("/website").Use(middleware.StaticWithConfig(middleware.StaticConfig{
		// This middleware is for serving SPA websites
		Index:      "index.html",
		HTML5:      true,
		IgnoreBase: false,
		Browse:     true,
		Filesystem: http.FS(echo.MustSubFS(s.websiteFS, "dist")),
	}))
	
	api := s.router.Group("/api")
	if err := s.managerRegistrator.Register(api); err != nil {
		return err
	}
	if err := s.usersRegistrator.Register(api); err != nil {
		return err
	}
	return s.router.Start(s.port)
}

func (s *server) Shutdown(ctx context.Context) error {
	return s.router.Shutdown(ctx)
}

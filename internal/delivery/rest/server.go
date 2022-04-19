package rest

import (
	"context"
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

	managerRegistrator Registrator
}

func NewServer(port string, tw, tr time.Duration, m Registrator) (*server, error) {
	e := echo.New()
	e.Server.ReadTimeout = tr
	e.Server.WriteTimeout = tw
	return &server{
		router:             e,
		port:               port,
		managerRegistrator: m,
	}, nil
}

func (s *server) Start() error {
	s.router.Use(middleware.CORS())
	manager := s.router.Group("")
	if err := s.managerRegistrator.Register(manager); err != nil {
		return err
	}
	return s.router.Start(s.port)
}

func (s *server) Shutdown(ctx context.Context) error {
	return s.router.Shutdown(ctx)
}

package server

import (
	"time"

	"github.com/labstack/echo/v4"
	"github.com/rasulov-emirlan/netping-manager/internal/pkg/watcher"
)

type Registrator interface {
	Register(router *echo.Group) error
}

type server struct {
	port    string
	router  *echo.Echo
	watcher *watcher.Watcher

	managerRegistrator Registrator
}

func NewServer(w *watcher.Watcher, port string, tw, tr time.Duration, m Registrator) (*server, error) {
	e := echo.New()
	e.Server.ReadTimeout = tr
	e.Server.WriteTimeout = tw
	return &server{
		router:             e,
		watcher:            w,
		port:               port,
		managerRegistrator: m,
	}, nil
}

func (s *server) Start() error {
	// s.router.POST("/config/location", addLocation(s.watcher))
	// s.router.DELETE("/config/location", removeLocation(s.watcher))
	// s.router.POST("/config/socket", addSocket(s.watcher))
	// s.router.DELETE("/config/socket", removeSocket(s.watcher))

	// s.router.POST("/control", setValue(s.watcher))
	// s.router.GET("/control", getAll(s.watcher))

	manager := s.router.Group("")
	if err := s.managerRegistrator.Register(manager); err != nil {
		return err
	}

	return s.router.Start(s.port)
}
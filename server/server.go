package server

import (
	"time"

	"github.com/labstack/echo/v4"
	"github.com/rasulov-emirlan/netping-manager/watcher"
)

type server struct {
	port    string
	router  *echo.Echo
	watcher *watcher.Watcher
}

func NewServer(w *watcher.Watcher, port string, tw, tr time.Duration) (*server, error) {
	e := echo.New()
	e.Server.ReadTimeout = tr
	e.Server.WriteTimeout = tw
	return &server{
		router:  e,
		watcher: w,
		port:    port,
	}, nil
}

func (s *server) Start() error {
	s.router.POST("/", setValue(s.watcher))
	s.router.GET("/", getAll(s.watcher))
	return s.router.Start(s.port)
}

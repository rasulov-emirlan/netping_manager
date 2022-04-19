package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/rasulov-emirlan/netping-manager/config"
	"github.com/rasulov-emirlan/netping-manager/internal/compositors"

	"github.com/rasulov-emirlan/netping-manager/internal/delivery/rest"
)

func main() {
	var cfg *config.Config
	var err error
	if len(os.Args) > 1 {
		cfg, err = config.NewConfig(os.Args[1:]...)
	} else {
		cfg, err = config.NewConfig()
	}
	if err != nil {
		log.Fatal(err)
	}

	h, close, err := compositors.NewManager(*cfg)
	if err != nil {
		log.Fatal(err)
	}
	server, err := rest.NewServer(cfg.Server.Port, cfg.Server.TimeoutWrite, cfg.Server.TimeoutRead, h)
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		if err := server.Start(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	if err := server.Shutdown(context.TODO()); err != nil {
		log.Fatal(err)
	}
	if err := close(); err != nil {
		log.Fatal(err)
	}
	log.Println("Gracefully shutting down :)")
}

package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/rasulov-emirlan/netping-manager/config"
	"github.com/rasulov-emirlan/netping-manager/internal/manager"

	"github.com/rasulov-emirlan/netping-manager/internal/delivery/rest"
	managerH "github.com/rasulov-emirlan/netping-manager/internal/manager/delivery/rest"
	"github.com/rasulov-emirlan/netping-manager/internal/pkg/watcher"
	"github.com/rasulov-emirlan/netping-manager/pkg/logger"
	"go.uber.org/zap"
)

func main() {
	cfg := &config.Config{}
	var err error
	if len(os.Args) > 1 {
		cfg, err = config.NewConfig(os.Args[1:]...)
	} else {
		cfg, err = config.NewConfig()
	}
	if err != nil {
		log.Fatal(err)
	}

	l := []*manager.Location{{
		ID:            1,
		Name:          "Ошская станция",
		RealLocation:  "Город Ош ул.Бакаева",
		SNMPaddress:   "192.168.0.100",
		SNMPcommunity: "SWITCH",
		SNMPport:      161,
		Sockets: []*manager.Socket{{
			ID:      1,
			Name:    "Кондиционер",
			SNMPmib: ".1.3.6.1.4.1.25728.8900.1.1.3.4",
			IsON:    false,
		}},
	}}

	w, err := watcher.NewWatcher(l)
	if err != nil {
		log.Fatal(err)
	}
	level := zap.NewAtomicLevelAt(zap.InfoLevel)
	z, err := logger.NewZap("logs.log", true, level)
	if err != nil {
		log.Fatal(err)
	}
	s, err := manager.NewService(w, z)
	if err != nil {
		log.Fatal(err)
	}
	h, err := managerH.NewHandler(s)
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
	log.Println("Gracefully shutting down :)")
}

package main

import (
	"context"
	"embed"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/rasulov-emirlan/netping-manager/config"
	"github.com/rasulov-emirlan/netping-manager/internal/compositors"
	"github.com/rasulov-emirlan/netping-manager/pkg/db"
	"github.com/rasulov-emirlan/netping-manager/pkg/logger"
	"go.uber.org/zap"

	"github.com/rasulov-emirlan/netping-manager/internal/delivery/rest"
)

//go:embed dist
var website embed.FS

func main() {
	var cfg *config.Config
	var err error
	// This application uses arguments as filenames for configs
	// config files have to be .env files
	// If you do not provide any filenames, this app will try to
	// find configs from your enviorment variables
	if len(os.Args) > 1 {
		cfg, err = config.NewConfig(os.Args[1:]...)
	} else {
		cfg, err = config.NewConfig()
	}
	if err != nil {
		log.Fatal(err)
	}

	// Craeting logger
	level := zap.NewAtomicLevelAt(zap.InfoLevel)
	z, logCloser, err := logger.NewZap(cfg.LogFilename, false, level)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := logCloser(); err != nil {
			log.Fatal(err)
		}
	}()

	// Connecting to database
	dbConn, err := db.NewMySQL(cfg.Database)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := dbConn.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	// Initializing all the handlers and services
	managerh, err := compositors.NewManager(*cfg, z, dbConn)
	if err != nil {
		log.Fatal(err)
	}
	usersH, err := compositors.NewUsers(cfg, z, dbConn)
	if err != nil {
		log.Fatal(err)
	}

	// Creating a new http rest server
	server, err := rest.NewServer(cfg.Server.Port, &website, cfg.Server.TimeoutWrite, cfg.Server.TimeoutRead, managerh, usersH)
	if err != nil {
		log.Fatal(err)
	}

	// Here our apiserver starts working in a separate goroutine
	go func() {
		if err := server.Start(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	// Here we block the main goroutine untile we get a SIGTERM from OS
	// And then we gracufully shutdown our server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	if err := server.Shutdown(context.TODO()); err != nil {
		log.Fatal(err)
	}
	log.Println("Gracefully shutting down :)")
}

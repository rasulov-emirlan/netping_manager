package main

import (
	"log"
	"os"

	"github.com/rasulov-emirlan/netping-manager/config"
	"github.com/rasulov-emirlan/netping-manager/server"
	"github.com/rasulov-emirlan/netping-manager/watcher"
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

	watcher, err := watcher.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	_, err = watcher.AddLocation("Ошская станция", "192.168.0.100", "SWITCH", 161)
	if err != nil {
		log.Fatal(err)
	}
	_, err = watcher.AddSocket("Ошская станция", "Кондиционер", ".1.3.6.1.4.1.25728.8900.1.1.3.4")
	if err != nil {
		log.Fatal(err)
	}
	v, err := watcher.Walk()
	if err != nil {
		log.Fatal(err)
	}
	log.Println(v)
	s, err := server.NewServer(watcher, cfg.Server.Port, cfg.Server.TimeoutWrite, cfg.Server.TimeoutRead)
	if err != nil {
		log.Fatal(err)
	}
	s.Start()
}

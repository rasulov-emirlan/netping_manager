package main

import (
	"log"
	"os"

	"github.com/gosnmp/gosnmp"
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
	conn := gosnmp.Default
	conn.Community = "SWITCH"
	conn.Target = "192.168.0.100"
	if err := conn.Connect(); err != nil {
		log.Fatal(err)
	}
	locations := map[string]watcher.Location{
		"Default": {
			Address: "192.168.0.100",
			Conn:    conn,
			Sockets: []watcher.Socket{
				{
					Name:    "Розетка",
					Address: ".1.3.6.1.4.1.25728.8900.1.1.3.4",
					Warning: "Are you sure? This will affect 4th line",
				},
			},
		}}
	watcher, err := watcher.NewWatcher(locations)
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

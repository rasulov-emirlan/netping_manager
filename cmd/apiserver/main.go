package main

import (
	"log"
	"os"

	"github.com/rasulov-emirlan/netping-manager/config"
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
}

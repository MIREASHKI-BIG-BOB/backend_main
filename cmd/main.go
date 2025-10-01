package main

import (
	"flag"
	"log"

	"github.com/MIREASHKI-BIG-BOB/backend_main/config"
	"github.com/MIREASHKI-BIG-BOB/backend_main/internal/server"
)

func main() {
	cfgPath := flag.String("c", "config/config.yaml", "path to config file")
	flag.Parse()

	cfg, err := config.ReadConfig(*cfgPath)
	if err != nil {
		log.Fatal(err)
	}

	srv, err := server.New(cfg, nil)
	if err != nil {
		log.Fatal(err)
	}

	srv.Run()
}

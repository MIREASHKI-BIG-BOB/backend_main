package main

import (
	"log"

	"github.com/MIREASHKI-BIG-BOB/backend_main/config"
	"github.com/MIREASHKI-BIG-BOB/backend_main/internal/server"
)

func main() {
	cfg, err := config.ReadConfig("config/config.yaml")
	if err != nil {
		log.Fatal(err)
	}

	srv, err := server.New(cfg, nil)
	if err != nil {
		log.Fatal(err)
	}

	srv.Run()
}

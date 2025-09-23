package main

import (
	"github.com/MIREASHKI-BIG-BOB/backend_main/config"
	"github.com/MIREASHKI-BIG-BOB/backend_main/internal/server"
)

func main() {
	cfg, err := config.ReadConfig("config/config.yaml")
	if err != nil {
		panic(err)
	}

	srv, err := server.New(cfg)
	if err != nil {
		panic(err)
	}

	if err = srv.Run(); err != nil {
		panic(err)
	}
}

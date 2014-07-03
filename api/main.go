package main

import (
	"log"

	cfg "github.com/bitnel/bitnel/api/config"
	"github.com/bitnel/bitnel/api/server"
)

func main() {
	config, err := cfg.LoadConfig("config.json")
	if err != nil {
		log.Fatal(err)
	}

	svr := server.New(config)
	svr.Start()
}

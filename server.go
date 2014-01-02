package main

import (
	"database/sql"
	"github.com/bitnel/bitnel-api/config"
	"log"
	"net/http"
)

var appConfig *config.Config
var db *sql.DB

func main() {
	var err error

	appConfig, err = config.LoadConfig("config.json")
	if err != nil {
		log.Fatal(err)
	}

	db, err = sql.Open("postgres", appConfig.Database)

	hd := routerHandler(router())
	log.Fatal(http.ListenAndServe(appConfig.ListenAddr, hd))
}

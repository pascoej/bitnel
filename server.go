package main

import (
	"database/sql"
	"github.com/bitnel/bitnel-api/config"
	_ "github.com/lib/pq"
	"log"
	"net/http"
)

var appConfig *config.Config
var db *sql.DB
var globalMatchingEngine *MatchingEngine

func main() {
	var err error

	appConfig, err = config.LoadConfig("config.json")
	if err != nil {
		log.Fatal(err)
	}

	db, err = sql.Open("postgres", appConfig.Database)
	if err != nil {
		log.Fatalln(err)
	}

	globalMatchingEngine = &MatchingEngine{}
	globalMatchingEngine.start()

	hd := routerHandler(router())
	log.Fatal(http.ListenAndServe(appConfig.ListenAddr, hd))
}

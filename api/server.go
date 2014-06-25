package main

import (
	"database/sql"
	"github.com/bitnel/bitnel/api/config"
	"github.com/bitnel/bitnel/api/matching"
	_ "github.com/lib/pq"
	"log"
	"net/http"
)

var appConfig *config.Config
var db *sql.DB
var globalMatchingEngine *matching.Engine

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

	globalMatchingEngine = matching.NewEngine(db, 10000)
	globalMatchingEngine.Start()

	hd := routerHandler(router())
	log.Fatal(http.ListenAndServe(appConfig.ListenAddr, hd))
}

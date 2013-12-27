package main

import (
	"database/sql"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"log"
	"net/http"
)

var appConfig *config
var db *sql.DB

func main() {
	var err error
	appConfig, err = loadConfig("config.json")
	if err != nil {
		log.Fatal(err)
	}

	// driver does not check if there's an actual connection made
	// do something about this later
	db, err = sql.Open("postgres", appConfig.Database)
	if err != nil {
		log.Fatal(err)
	}

	r := mux.NewRouter()
	r.HandleFunc("/", homeHandler)
	http.Handle("/", r)

	log.Fatal(http.ListenAndServe(appConfig.ListenAddr, nil))
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome to the Bitnel API!111!11!111")
}

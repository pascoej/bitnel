package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"io/ioutil"
	"log"
	"net/http"
)

var config map[string]string
var db *sql.DB

func main() {
	var err error
	config, err = loadConfig("config.json")
	if err != nil {
		log.Fatal(err)
	}

	// driver does not check if there's an actual connection made
	// do something about this later
	db, err = sql.Open("postgres", config["database"])
	if err != nil {
		log.Fatal(err)
	}

	r := mux.NewRouter()
	r.HandleFunc("/", homeHandler)
	http.Handle("/", r)

	log.Fatal(http.ListenAndServe(config["listenAddr"], nil))
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome to the Bitnel API!111!11!111")
}

// TODO: Don't relay same errors
func loadConfig(filename string) (map[string]string, error) {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var m map[string]string
	err = json.Unmarshal(b, &m)

	if err != nil {
		return nil, err
	}

	return m, nil
}

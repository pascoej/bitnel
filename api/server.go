package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/bitnel/bitnel/api/config"
	"github.com/bitnel/bitnel/api/matching"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
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

func router() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/", index)

	r.Handle("/user", apiHandler(createUser)).Methods("POST")
	r.Handle("/user", apiHandler(updateUser)).Methods("PUT")
	r.Handle("/user/accounts", use(getAccounts, oauthAuth)).Methods("GET")

	r.Handle("/markets/{currencyPair}/orders", use(listOrder, findMarket)).Methods("GET")

	r.Handle("/accounts/{accountUuid}/orders/{orderUuid}", use(getOrder, oauthAuth, findAccount, findOrder)).Methods("GET")
	r.Handle("/accounts/{accountUuid}/orders", use(createOrder, findAccount, oauthAuth)).Methods("POST")
	r.Handle("/accounts/{accountUuid}/orders/{orderUuid}", use(deleteOrder, oauthAuth, findAccount, findOrder)).Methods("DELETE")

	r.NotFoundHandler = http.Handler(apiHandler(notFound))

	return r
}

func routerHandler(rtr *mux.Router) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rtr.ServeHTTP(w, r)
	}
}

func use(h apiHandler, middleware ...func(apiHandler) apiHandler) apiHandler {
	for _, m := range middleware {
		h = m(h)
	}

	return h
}

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Welcome to the Bitnel API!")
}

func notFound(w http.ResponseWriter, r *http.Request) *serverError {
	return writeError(w, errNotFound)
}

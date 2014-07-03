package server

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"

	cfg "github.com/bitnel/bitnel/api/config"
	"github.com/bitnel/bitnel/api/matching"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	_ "github.com/lib/pq"
)

type server struct {
	config         *cfg.Config
	db             *sql.DB
	matchingEngine *matching.Engine
	decoder        *schema.Decoder
}

func New(config *cfg.Config) *server {
	return &server{
		config: config,
	}
}

type middleware func(apiHandler) apiHandler

func (s *server) Start() error {
	var err error
	s.db, err = sql.Open("postgres", s.config.Database)
	if err != nil {
		return errors.New("server: could not connect to database")
	}

	s.matchingEngine.Start()

	if err = http.ListenAndServe(s.config.ListenAddr, router(s)); err != nil {
		return errors.New("server: could not start http server")
	}

	return nil
}

var routes = []struct {
	method  string
	path    string
	handler apiHandler
	mw      []middleware
}{
	{
		method:  "POST",
		path:    "/users",
		handler: createUser,
		mw:      nil},
	{
		method:  "PUT",
		path:    "/user",
		handler: updateUser,
		mw:      nil},
	{
		method:  "GET",
		path:    "/user/accounts",
		handler: listUserAccounts,
		mw:      []middleware{findMarket}},
	{
		method:  "GET",
		path:    "/markets",
		handler: listMarkets,
		mw:      []middleware{findMarket}},
	{
		method:  "GET",
		path:    "/markets/{currencyPair}",
		handler: getMarket,
		mw:      []middleware{findMarket}},
	{
		method:  "GET",
		path:    "/markets/{currencyPair}/orders",
		handler: listMarketOrders,
		mw:      []middleware{findMarket}},
	{
		method:  "GET",
		path:    "/accounts/{accountUuid}",
		handler: getAccount,
		mw:      []middleware{findAccount}},
	{
		method:  "GET",
		path:    "/accounts/{accountUuid}/orders",
		handler: listAccountOrders,
		mw:      []middleware{findAccount, oauthAuth}},
	{
		method:  "POST",
		path:    "/accounts/{accountUuid}/orders",
		handler: createAccountOrder,
		mw:      []middleware{findAccount, oauthAuth}},
	{
		method:  "GET",
		path:    "/accounts/{accountUuid}/orders/{orderUuid}",
		handler: getAccountOrder,
		mw:      []middleware{oauthAuth, findAccount, findOrder}},
	{
		method:  "DELETE",
		path:    "/accounts/{accountUuid}/orders/{orderUuid}",
		handler: cancelAccountOrder,
		mw:      []middleware{oauthAuth, findAccount, findOrder}},
}

func router(s *server) *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/", index)
	r.NotFoundHandler = http.Handler(makeHandler(s, notFound))

	for _, e := range routes {
		r.HandleFunc(e.path, makeHandler(s, e.handler, e.mw...)).Methods(e.method)
	}

	return r
}

// mfw mfw
func makeHandler(s *server, h apiHandler, mfw ...middleware) http.HandlerFunc {
	for _, m := range mfw {
		h = m(h)
	}

	return func(w http.ResponseWriter, r *http.Request) {
		if err := h(s, w, r); err != nil {
			log.Println(err.Error())
			writeError(w, errServerError)
		}
	}
}

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Welcome to the Bitnel API!")
}

func notFound(s *server, w http.ResponseWriter, r *http.Request) *serverError {
	return writeError(w, errNotFound)
}

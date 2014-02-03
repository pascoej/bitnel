package main

import (
	"github.com/gorilla/mux"
	"net/http"
)

// A http handler must be enclosed in an `apiHandler` so that the handlers are
// compatible with `http.Handler`.
func router() *mux.Router {
	r := mux.NewRouter()

	r.Handle("/", apiHandler(useMiddleware(indexHandler)))

	//su := r.PathPrefix("/users").Subrouter()
	r.Handle("/users", apiHandler(createUserHandler)).Methods("POST")
	r.Handle("/users", apiHandler(updateUserHandler)).Methods("PUT")
	r.Handle("/accounts", apiHandler(useMiddleware(getAccountsHandler, oauthTokenUserFinder))).Methods("GET")
	oauthr := r.PathPrefix("/oauth").Subrouter()
	oauthr.Handle("/token", apiHandler(oauthTokenHandler)).Methods("POST")

	sm := r.PathPrefix("/markets").Subrouter()
	sm.Handle("/{currencyPair}/accounts/{accountUuid}/orders/{orderUuid}", apiHandler(useMiddleware(getOrderHandler, marketFinder, oauthTokenUserFinder, accountFinder))).Methods("GET")
	sm.Handle("/{currencyPair}/accounts/{accountUuid}/orders", apiHandler(useMiddleware(createOrderHandler, marketFinder, oauthTokenUserFinder, accountFinder))).Methods("POST")
	sm.Handle("/{currencyPair}/orders", apiHandler(useMiddleware(listOrderHandler, marketFinder))).Methods("GET")

	r.NotFoundHandler = http.Handler(apiHandler(notFoundHandler))

	return r
}

// For testing purposes (coming later).
func routerHandler(rtr *mux.Router) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rtr.ServeHTTP(w, r)
	}
}

// Goes through each middleware
func useMiddleware(handler apiHandler, middleware ...func(apiHandler) apiHandler) apiHandler {
	for _, m := range middleware {
		handler = m(handler)
	}

	return handler
}

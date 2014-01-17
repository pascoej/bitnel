package main

import (
	"github.com/gorilla/mux"
	"net/http"
)

// A http handler must be enclosed in an `apiHandler` so that the handlers are
// compatible with `http.Handler`.
func router() *mux.Router {
	r := mux.NewRouter()

	r.Handle("/", apiHandler(indexHandler))

	su := r.PathPrefix("/users").Subrouter()
	su.Handle("/", apiHandler(createUserHandler)).Methods("POST")
	su.Handle("/", apiHandler(updateUserHandler)).Methods("PUT")

	sm := r.PathPrefix("/markets").Subrouter()
	sm.Handle("/{currencyPair}/orders", apiHandler(useMiddleware(createOrderHandler, marketFinder))).Methods("POST", "GET")

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

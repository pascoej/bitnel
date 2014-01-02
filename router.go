package main

import (
	"github.com/gorilla/mux"
	"net/http"
)

func router() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/", indexHandler)
	r.HandleFunc("/users", createUserHandler).Methods("POST")
	r.HandleFunc("/users", updateUserHandler).Methods("PUT")
	r.HandleFunc("/orders", createOrderHandler).Methods("POST")

	return r
}

func routerHandler(rtr *mux.Router) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rtr.ServeHTTP(w, r)
	}
}

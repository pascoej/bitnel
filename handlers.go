package main

import (
	"fmt"
	"net/http"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome to the Bitnel API!111!11!111")
}

func createOrderHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "hi")
}

func createUserHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "hi")
}

func updateUserHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "hi")
}

package main

import (
	"fmt"
	"net/http"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome to the Bitnel API!")
}

func createOrderHandler(w http.ResponseWriter, r *http.Request) {}

func createUserHandler(w http.ResponseWriter, r *http.Request) {}

func updateUserHandler(w http.ResponseWriter, r *http.Request) {}

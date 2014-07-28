package main

import (
	"log"
	"net/http"
)

func main() {
	http.Handle("/js/", http.FileServer(http.Dir("./static/")))
	http.Handle("/css/", http.FileServer(http.Dir("./static/")))
	http.Handle("/views/", http.FileServer(http.Dir("./static/")))
	http.HandleFunc("/", all)
	log.Fatal(http.ListenAndServe(":8337", nil))
}

func all(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/index.html")
}

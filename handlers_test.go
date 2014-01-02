package main

import (
	"net/http"
	"net/http/httptest"
)

func makeRequest(method string, handler func(w http.ResponseWriter, r *http.Request)) (*httptest.ResponseRecorder, error) {
	req, err := http.NewRequest("GET", "", nil)
	if err != nil {
		return nil, err
	}

	w := httptest.NewRecorder()
	handler(w, req)

	return w, nil
}

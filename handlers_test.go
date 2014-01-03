package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"runtime"
	"testing"
)

var handlerTests = []struct {
	handler func(http.ResponseWriter, *http.Request)
	method  string
	path    string

	// body is meant to handle a json request
	body       map[string]interface{}
	exStatus   int
	exResponse interface{}
}{
	{indexHandler, "GET", "", nil, 200, "Welcome to the Bitnel API!\n"},
	{createOrderHandler, "GET", "", nil, 200, nil},
	{createUserHandler, "GET", "", nil, 200, nil},
	{updateUserHandler, "GET", "", nil, 200, nil},
}

func TestHandlers(t *testing.T) {
	for _, tt := range handlerTests {
		// marshal our request body
		b, err := json.Marshal(tt.body)
		if err != nil {
			t.Error(err)
		}

		req, err := http.NewRequest(tt.method, tt.path, bytes.NewReader(b))
		if err != nil {
			t.Error(err)
		}

		w := httptest.NewRecorder()
		tt.handler(w, req)

		// check if response code matches that of expected
		if tt.exStatus != w.Code {
			t.Errorf("%s: expected status code %s, got %s", funcName(tt.handler), tt.exStatus, w.Code)
		}

		// this getto solution should work for now

		var eb []byte

		switch tty := tt.exResponse.(type) {
		case map[string]interface{}:
			eb, err = json.Marshal(tty)
		case string:
			eb = []byte(tty)
		}

		bs := w.Body.Bytes()

		if !bytes.Equal(eb, bs) {
			t.Errorf("%s: expected body %s, got %s", funcName(tt.handler), tt.exResponse, bs)
		}
	}
}

func funcName(fn interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name()
}

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// API handler is a custom type that gives us the chance to log server errors
// when handlers return a server error.
type apiHandler func(http.ResponseWriter, *http.Request) *serverError

func (fn apiHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := fn(w, r); err != nil {
		log.Println(err.Error())
		writeError(w, errServerError)
	}
}

// A server error is returned by handlers signaling that some error the handler
// encountered should be logged for investigation.
type serverError struct {
	Err error
	Msg string
}

func (s *serverError) Error() string {
	return fmt.Sprint("%s: %s", s.Msg, s.Err)
}

type AccountUuid string // for middleware

// API error codes describe the specifics of the error. The user already knows
// if the error was caused by them or us, but error codes allow the user to know
// if the resource was `not_found` or they failed the `input_validation`.
type apiErrorCode int

const (
	errCodeServerError apiErrorCode = iota
	errCodeInputValidation
	errCodeNotFound
	errCodeAuth
	errCodeTooBusy
)

func (e apiErrorCode) String() string {
	switch e {
	case errCodeServerError:
		return "internal_server_error"
	case errCodeInputValidation:
		return "input_validation"
	case errCodeNotFound:
		return "not_found"
	case errCodeAuth:
		return "err_auth"
	case errCodeTooBusy:
		return "too_busy"
	}

	return ""
}

// API error types describe the nature of the error returned to the user.
// Currently there is only `server` and `request`. This tells if the error was
// caused by the `request` (user), or `server` (us).
type apiErrorType int

const (
	// Something the user screwed up. Go away!
	errTypeRequest apiErrorType = iota

	// Something we screwed up. We're soooo sorry. /s (jk, we really are D:)
	errTypeServer
)

func (e apiErrorType) String() string {
	switch e {
	case errTypeRequest:
		return "request"
	case errTypeServer:
		return "server"
	}

	return ""
}

// API errors are predefined error types that are immediately available to
// return to the user. They must be the only thing returned to the user in case
// of an unprocessible error.
type apiError struct {
	Type   apiErrorType
	Code   apiErrorCode
	Status int
	Msg    string
}

var (
	errInputValidation = &apiError{errTypeRequest, errCodeInputValidation, http.StatusBadRequest, "Your input could not be validated"}
	errServerError     = &apiError{errTypeServer, errCodeServerError, http.StatusInternalServerError, "Something went wrong on our side"}
	errNotFound        = &apiError{errTypeRequest, errCodeNotFound, http.StatusNotFound, "We cannot find that resource"}
	errAuth            = &apiError{errTypeRequest, errCodeAuth, http.StatusForbidden, "You are not authenticated"}
	errTooBusy = 	&apiError{errTypeServer, errCodeTooBusy, http.StatusServiceUnavailable, "The server is too busy"}
)

// `writeError` notifies the user of an API error. This calls `writeJson` too.
func writeError(w http.ResponseWriter, err *apiError) *serverError {
	w.WriteHeader(err.Status)

	m := map[string]interface{}{
		"message": err.Msg,
		"type":    err.Type.String(),
		"code":    err.Code.String(),
	}

	return writeJson(w, m)
}

// `writeJson` should be used to communicate with user.
func writeJson(w http.ResponseWriter, m interface{}) *serverError {
	if err := json.NewEncoder(w).Encode(m); err != nil {
		return &serverError{err, "could not encode json"}
	}

	return nil
}

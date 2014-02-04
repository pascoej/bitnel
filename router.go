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

	// CREATE A NEW USER
	// POST /users
	r.Handle("/users", apiHandler(createUserHandler)).Methods("POST")

	// UPDATE A USER'S INFORMATION
	// PUT /users
	r.Handle("/users", apiHandler(updateUserHandler)).Methods("PUT")

	oauthr := r.PathPrefix("/oauth").Subrouter()

	// POST /oauth/token
	oauthr.Handle("/token", apiHandler(oauthTokenHandler)).Methods("POST")

	// GET /markets/BTCUSD/orders
	r.Handle("/markets/{currencyPair}/orders", apiHandler(useMiddleware(listOrderHandler, marketFinder))).Methods("GET")

	sm := r.PathPrefix("/accounts").Subrouter()

	// GET ORDER BY ORDER UUID ASSCOCIATED WITH AN ACCOUNT
	// GET /accounts/0564bdb5-c35f-4f9f-b1bb-b574d201fa90/orders/a564bdd2-c35f-4f9f-b1bd-b574d201fa90
	sm.Handle("/{accountUuid}/orders/{orderUuid}", apiHandler(useMiddleware(getOrderHandler, oauthTokenUserFinder, accountFinder))).Methods("GET")

	// CREATE AN ORDER ASSOCIATED WITH AN ACCOUNT
	// POST /accounts/0564bdb5-c35f-4f9f-b1bb-b574d201fa90/orders
	sm.Handle("/{accountUuid}/orders", apiHandler(useMiddleware(createOrderHandler, oauthTokenUserFinder, accountFinder))).Methods("POST")

	// SINGULAR RESOURCE -- USER NEEDS TO BE AUTHENTICATED
	um := r.PathPrefix("/user").Subrouter()

	// GET ACCOUNTS OF A USER
	// GET /user/accounts
	um.Handle("/accounts", apiHandler(useMiddleware(getAccountsHandler, oauthTokenUserFinder))).Methods("GET")

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

package rest

import (
	"log"
	"net/http"

	"github.com/buraksekili/store-service/db"
	"github.com/gorilla/mux"
)

// ServerREST serves a REST API for the product service.
// It takes addr to indicate which address should the server listen,
// and database handler which will operate on the database.
func ServerREST(addr string, dh db.DBHandler) error {
	r := mux.NewRouter()

	h := newProductServiceHandler(dh)

	// obtain sub-router for product router
	pr := r.PathPrefix("/products").Subrouter()

	pr.Methods("GET").Path("").HandlerFunc(h.getProducts)

	log.Printf("Listening on %s\n", addr)
	return http.ListenAndServe(addr, r)
}

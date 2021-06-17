package rest

import (
	"log"
	"net/http"

	"github.com/buraksekili/store-service/db"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

// ServerREST serves a REST API for the product service.
// It takes addr for which address should the server listen,
// and database handler which will operate on the database.
func ServerREST(addr string, dh db.DBHandler) error {
	r := mux.NewRouter()

	h := newProductServiceHandler(dh)

	// obtain sub-router for product router
	pr := r.PathPrefix("/products").Subrouter()
	pr.Methods("GET").Path("").HandlerFunc(h.getProducts)
	pr.Methods("GET").Path("/{product_id}").HandlerFunc(h.findProduct)
	pr.Methods("GET").Path("/vendor/{vendor_id}").HandlerFunc(h.findProductsByVendor)
	pr.Methods("POST").Path("").HandlerFunc(h.addProduct)

	cr := r.PathPrefix("/comments").Subrouter()
	cr.Methods("GET").Path("").HandlerFunc(h.getComments)
	cr.Methods("POST").Path("/{product_id}").HandlerFunc(h.addComment)

	vr := r.PathPrefix("/vendors").Subrouter()
	vr.Methods("GET").Path("").HandlerFunc(h.getVendors)
	vr.Methods("POST").Path("").HandlerFunc(h.addVendor)

	log.Printf("Productservice listening on %s\n", addr)

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
	})

	handler := c.Handler(r)
	return http.ListenAndServe(addr, handler)
}

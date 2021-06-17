package rest

import (
	"log"
	"net/http"

	"github.com/buraksekili/store-service/db"
	"github.com/gorilla/mux"
)

func ServerREST(addr string, dh db.DBHandler) error {
	r := mux.NewRouter()

	h := newUserServiceHandler(dh)

	// obtain sub-router for product router
	ur := r.PathPrefix("/users").Subrouter()
	ur.Methods("GET").Path("").HandlerFunc(h.getUsers)
	ur.Methods("GET").Path("/{user_id}").HandlerFunc(h.findUser)
	ur.Methods("POST").Path("/signup").HandlerFunc(h.addUser)
	ur.Methods("POST").Path("/login").HandlerFunc(h.login)

	log.Printf("Userservice listening on %s\n", addr)
	return http.ListenAndServe(addr, r)
}

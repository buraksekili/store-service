package rest

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/buraksekili/store-service/db"
)

type productServiceHandler struct {
	dbHandler db.DBHandler
}

func newProductServiceHandler(dh db.DBHandler) *productServiceHandler {
	return &productServiceHandler{dh}
}

func (h *productServiceHandler) getProducts(w http.ResponseWriter, r *http.Request) {
	products, err := h.dbHandler.GetProducts()
	if err != nil {
		// TODO: check for a proper status code
		w.WriteHeader(400)
		fmt.Fprintf(w, "cannot get products, err: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json;charset=utf8")
	err = json.NewEncoder(w).Encode(&products)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "Error occured while trying encode events to JSON %s", err)
	}
}

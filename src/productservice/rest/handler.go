package rest

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"gopkg.in/mgo.v2/bson"

	"github.com/buraksekili/store-service/db"
)

type productServiceHandler struct {
	dbHandler db.DBHandler
}

func newProductServiceHandler(dh db.DBHandler) *productServiceHandler {
	return &productServiceHandler{dh}
}

func (h *productServiceHandler) addProduct(w http.ResponseWriter, r *http.Request) {
	p := db.Product{}
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		w.WriteHeader(400)
		fmt.Fprintf(w, "cannot add product, invalid r.Body")
		return
	}

	byteID, err := h.dbHandler.AddProduct(p)
	if err != nil {
		w.WriteHeader(400)
		fmt.Fprintf(w, "cannot add product, err: %v", err)
		return
	}

	p.ID = bson.ObjectId(byteID)
	w.Header().Set("Content-Type", "application/json;charset=utf8")
	w.WriteHeader(201)
	if err := json.NewEncoder(w).Encode(p); err != nil {
		log.Printf("cannot encode product, err: %v", err)
	}
}

func (h *productServiceHandler) addVendor(w http.ResponseWriter, r *http.Request) {
	v := db.Vendor{}
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		w.WriteHeader(400)
		fmt.Fprintf(w, "cannot add vendor, invalid r.Body")
		return
	}

	byteID, err := h.dbHandler.AddVendor(v)
	if err != nil {
		w.WriteHeader(400)
		fmt.Fprintf(w, "cannot add vendor, err: %v", err)
		return
	}
	v.ID = bson.ObjectId(byteID)

	w.Header().Set("Content-Type", "application/json;charset=utf8")
	w.WriteHeader(201)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Printf("cannot encode vendor, err: %v", err)
	}
}

func (h *productServiceHandler) addComment(w http.ResponseWriter, r *http.Request) {
	c := db.Comment{}
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		w.WriteHeader(400)
		fmt.Fprintf(w, "cannot obtain comment, invalid r.Body, err: %v", err)
		return
	}

	vars := mux.Vars(r)
	pID, ok := vars["product_id"]
	if !ok {
		w.WriteHeader(400)
		fmt.Fprintf(w, "cannot obtain 'product_id' parameter")
		return
	}

	c.Date = time.Now()
	byteID, err := h.dbHandler.AddComment(c, pID)
	if err != nil {
		w.WriteHeader(400)
		fmt.Fprintf(w, "cannot add comment, err: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf8")
	w.WriteHeader(201)
	c.ID = bson.ObjectId(byteID)
	if err := json.NewEncoder(w).Encode(c); err != nil {
		log.Printf("cannot encode comment, err: %v", err)
	}
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
		fmt.Fprintf(w, "cannot encode products %s", err)
	}
}

func (h *productServiceHandler) findProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pID, ok := vars["product_id"]
	if !ok {
		w.WriteHeader(400)
		fmt.Fprintf(w, "cannot obtain 'product_id' parameter")
		return
	}

	p, err := h.dbHandler.FindProduct(pID)
	if err != nil {
		w.WriteHeader(400)
		fmt.Fprintf(w, "cannot find product %s, err: %v", pID, err)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf8")
	w.WriteHeader(201)
	json.NewEncoder(w).Encode(p)
}

func (h *productServiceHandler) findProductsByVendor(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	vID, ok := vars["vendor_id"]
	if !ok {
		w.WriteHeader(400)
		fmt.Fprintf(w, "cannot obtain 'vendor_id' parameter")
		return
	}

	products, err := h.dbHandler.FindProductsByVendor(vID)
	if err != nil {
		w.WriteHeader(400)
		fmt.Fprintf(w, "cannot find product %s, err: %v", vID, err)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf8")
	w.WriteHeader(201)
	json.NewEncoder(w).Encode(products)
}

func (h *productServiceHandler) getComments(w http.ResponseWriter, r *http.Request) {
	comments, err := h.dbHandler.GetComments()
	if err != nil {
		w.WriteHeader(400)
		fmt.Fprintf(w, "cannot get comments, err: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json;charset=utf8")
	err = json.NewEncoder(w).Encode(&comments)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "cannot encode comments %s", err)
	}
}

func (h *productServiceHandler) getVendors(w http.ResponseWriter, r *http.Request) {
	vendors, err := h.dbHandler.GetVendors()
	if err != nil {
		w.WriteHeader(400)
		fmt.Fprintf(w, "cannot get vendors, err: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json;charset=utf8")
	err = json.NewEncoder(w).Encode(&vendors)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "cannot encode vendors %s", err)
	}
}

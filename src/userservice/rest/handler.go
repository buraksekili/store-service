package rest

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"gopkg.in/mgo.v2/bson"

	"github.com/buraksekili/store-service/db"
)

type userServieHandler struct {
	dbHandler db.DBHandler
}

func (h userServieHandler) addUser(w http.ResponseWriter, r *http.Request) {
	u := db.User{}
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		w.WriteHeader(400)
		fmt.Fprintf(w, "cannot add user, invalid r.Body")
		return
	}

	byteID, err := h.dbHandler.AddUser(u)
	if err != nil {
		w.WriteHeader(400)
		fmt.Fprintf(w, "cannot add product, err: %v", err)
		return
	}

	u.ID = bson.ObjectId(byteID)

	w.Header().Set("Content-Type", "application/json;charset=utf8")
	w.WriteHeader(201)
	if err := json.NewEncoder(w).Encode(u); err != nil {
		log.Printf("cannot encode product, err: %v", err)
	}
}

func (h userServieHandler) getUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.dbHandler.GetUsers()
	if err != nil {
		w.WriteHeader(400)
		fmt.Fprintf(w, "cannot add product, err: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json;charset=utf8")
	w.WriteHeader(201)
	if err := json.NewEncoder(w).Encode(users); err != nil {
		log.Printf("cannot encode product, err: %v", err)
	}
}

func (h userServieHandler) findUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	uID, ok := vars["user_id"]
	if !ok {
		w.WriteHeader(400)
		fmt.Fprintf(w, "cannot obtain 'user_id' parameter")
		return
	}

	user, err := h.dbHandler.FindUser(uID)
	if err != nil {
		w.WriteHeader(400)
		fmt.Fprintf(w, "cannot find user %s, err: %v", user.ID, err)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf8")
	w.WriteHeader(201)
	json.NewEncoder(w).Encode(user)
}

func (h userServieHandler) login(w http.ResponseWriter, r *http.Request) {

	u := db.User{}
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		w.WriteHeader(400)
		fmt.Fprintf(w, "cannot decode r.Body %s, %v", r.Body, err)
		return
	}

	user, err := h.dbHandler.FindUserByName(u.Username)
	if err != nil {
		w.WriteHeader(400)
		fmt.Fprintf(w, "cannot find username %s, err: %v", user, err)
		return
	}

	if user.Password != u.Password {
		w.WriteHeader(400)
		fmt.Fprintf(w, "cannot login to %s, err: invalid password", user.Username)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf8")
	w.WriteHeader(201)
	json.NewEncoder(w).Encode(user)
}

func newUserServiceHandler(dh db.DBHandler) *userServieHandler {
	return &userServieHandler{dh}
}

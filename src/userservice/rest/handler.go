package rest

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	amqp2 "github.com/buraksekili/store-service/amqp"

	"github.com/gorilla/mux"

	"gopkg.in/mgo.v2/bson"

	"github.com/buraksekili/store-service/db"
)

type userServiceHandler struct {
	dbHandler db.DBHandler
	publisher amqp2.AMQPPublisher
}

func (h *userServiceHandler) addUser(w http.ResponseWriter, r *http.Request) {
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

	userEvent := &amqp2.AddUserEvent{
		ID:       u.ID,
		Username: u.Username,
		Email:    u.Email,
		Password: u.Password,
	}
	if err := h.publisher.Publish(userEvent); err != nil {
		log.Printf("cannot publish event %#v, err: %v", userEvent, err)
	}

	w.Header().Set("Content-Type", "application/json;charset=utf8")
	w.WriteHeader(201)
	if err := json.NewEncoder(w).Encode(u); err != nil {
		log.Printf("cannot encode product, err: %v", err)
	}
}

func (h *userServiceHandler) getUsers(w http.ResponseWriter, r *http.Request) {
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

func (h *userServiceHandler) findUser(w http.ResponseWriter, r *http.Request) {
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

func (h *userServiceHandler) login(w http.ResponseWriter, r *http.Request) {

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

func newUserServiceHandler(dh db.DBHandler, p amqp2.AMQPPublisher) *userServiceHandler {
	return &userServiceHandler{dh, p}
}

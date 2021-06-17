package amqp

import "gopkg.in/mgo.v2/bson"

type Event interface {
	Name() string
}

type AddUserEvent struct {
	ID       bson.ObjectId `bson:"_id"`
	Username string        `json:"username"`
	Email    string        `json:"email"`
	Password string        `json:"password"`
}

func (ue *AddUserEvent) Name() string {
	return "add_user"
}

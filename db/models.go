// Package db contains model structs for databases.

package db

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

type Product struct {
	ID          bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Name        string
	Category    string
	Description string
	ImageURL    string
	Price       float32
	Stock       int
	Vendor      Vendor
	Comments    []Comment
}

type Vendor struct {
	ID          bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Name        string
	Description string
}

type Comment struct {
	ID      bson.ObjectId `bson:"_id"`
	Owner   string
	Content string
	Date    time.Time
}

type User struct {
	ID       bson.ObjectId `bson:"_id"`
	Username string        `json:"username"`
	Email    string        `json:"email"`
	Password string        `json:"password"`
}

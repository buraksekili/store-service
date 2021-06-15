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
	Price       float32
	Stock       int
	Vendor      Vendor
	Comments    []Comment
}

type Vendor struct {
	Name        string
	Description string
}

type Comment struct {
	ID      bson.ObjectId `bson:"_id"`
	Owner   string
	Content string
	Date    time.Time
}

package mongo

import (
	"github.com/buraksekili/store-service/db"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	DB       = "store"
	PRODUCTS = "products"
)

type MongoDBLayer struct {
	session *mgo.Session
}

func NewMongoDBLayer(url string) (db.DBHandler, error) {
	s, err := mgo.Dial(url)
	return &MongoDBLayer{s}, err
}

func (m *MongoDBLayer) AddProduct(product db.Product) ([]byte, error) {
	s := m.session.Copy()
	defer s.Close()

	if !product.ID.Valid() {
		product.ID = bson.NewObjectId()
	}

	return []byte(product.ID), s.DB(DB).C(PRODUCTS).Insert(&product)
}

func (m *MongoDBLayer) GetProducts() ([]db.Product, error) {
	s := m.session.Copy()
	defer s.Close()
	products := []db.Product{}
	err := s.DB(DB).C(PRODUCTS).Find(nil).All(&products)
	return products, err
}

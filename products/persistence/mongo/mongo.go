package mongo

import (
	"context"
	"fmt"

	"github.com/buraksekili/store-service/products"
	"github.com/buraksekili/store-service/users"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	DB       = "store"
	PRODUCTS = "products"
	VENDORS  = "vendors"
	COMMENTS = "comments"
)

type MongoDBLayer struct {
	session *mgo.Session
}

func NewMongoDBLayer(url string) (products.ProductRepository, error) {
	s, err := mgo.Dial(url)
	return &MongoDBLayer{s}, err
}

func (m *MongoDBLayer) Save(_ context.Context, product products.Product) (string, error) {
	s := m.session.Copy()
	defer s.Close()

	if !product.ID.Valid() {
		product.ID = bson.NewObjectId()
	}

	v := users.Vendor{}
	err := s.DB(DB).C(VENDORS).Find(bson.M{"_id": product.Vendor.ID}).One(&v)
	if err != nil {
		return "", fmt.Errorf("cannot find vendor %s: %v", product.Vendor.ID, err)
	}
	product.Vendor = v
	if !product.Vendor.ID.Valid() {
		return "", fmt.Errorf("invalid ID for the vendor, %v", product.Vendor.ID)
	}

	return product.ID.String(), s.DB(DB).C(PRODUCTS).Insert(&product)
}

func (m *MongoDBLayer) FindByID(_ context.Context, productID string) (products.Product, error) {
	s := m.session.Copy()
	defer s.Close()

	p := products.Product{}
	if err := s.DB(DB).C(PRODUCTS).Find(bson.M{"_id": bson.ObjectIdHex(productID)}).One(&p); err != nil {
		return p, fmt.Errorf("cannot find product %s, err: %v", productID, err)
	}

	return p, nil
}

func (m *MongoDBLayer) FindByVendorID(_ context.Context, offset, limit int, vendorID string) ([]products.Product, error) {
	s := m.session.Copy()
	defer s.Close()

	p := []products.Product{}
	if !bson.ObjectId(vendorID).Valid() {
		return p, fmt.Errorf("invalid vendorID as %s", vendorID)
	}

	mID := bson.M{"$match": bson.M{"vendor._id": bson.ObjectIdHex(vendorID)}}
	pipeLine := []bson.M{mID}

	if err := s.DB(DB).C(PRODUCTS).Pipe(pipeLine).All(&p); err != nil {
		return p, fmt.Errorf("cannot find products %s, err: %v", vendorID, err)
	}

	return p, nil
}

func (m *MongoDBLayer) GetProducts(_ context.Context, offset, limit int) ([]products.Product, error) {
	// TODO: implement offset and limit.
	s := m.session.Copy()
	defer s.Close()
	products := []products.Product{}
	err := s.DB(DB).C(PRODUCTS).Find(nil).All(&products)
	return products, err
}

func (m *MongoDBLayer) GetComments() ([]products.Comment, error) {
	s := m.session.Copy()
	defer s.Close()
	comments := []products.Comment{}
	err := s.DB(DB).C(COMMENTS).Find(nil).All(&comments)
	return comments, err
}

func (m *MongoDBLayer) AddComment(comment products.Comment, id string) ([]byte, error) {
	s := m.session.Copy()
	defer s.Close()

	if !comment.ID.Valid() {
		comment.ID = bson.NewObjectId()
	}

	err := s.DB(DB).C(PRODUCTS).Update(
		bson.M{"_id": bson.ObjectIdHex(id)},
		bson.M{"$push": bson.M{"comments": comment}},
	)
	if err != nil {
		return nil, fmt.Errorf("cannot update product for comment %s, err: %v", id, err)
	}

	if err = s.DB(DB).C(COMMENTS).Insert(comment); err != nil {
		return nil, fmt.Errorf("cannot save comment %s, err: %v", id, err)
	}

	return []byte(comment.ID), nil
}

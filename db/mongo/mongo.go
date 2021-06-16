package mongo

import (
	"fmt"

	"github.com/buraksekili/store-service/db"
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

	v := db.Vendor{}
	err := s.DB(DB).C(VENDORS).Find(bson.M{"name": product.Vendor.Name}).One(&v)
	if err != nil {
		return nil, fmt.Errorf("cannot find vendor %s: %v", product.Vendor.ID, err)
	}
	product.Vendor = v

	if !product.Vendor.ID.Valid() {
		return nil, fmt.Errorf("invalid ID for the vendor, %v", product.Vendor.ID)
	}

	return []byte(product.ID), s.DB(DB).C(PRODUCTS).Insert(&product)
}

func (m *MongoDBLayer) AddVendor(vendor db.Vendor) ([]byte, error) {
	s := m.session.Copy()
	defer s.Close()

	if !vendor.ID.Valid() {
		vendor.ID = bson.NewObjectId()
	}

	return []byte(vendor.ID), s.DB(DB).C(VENDORS).Insert(&vendor)
}

func (m *MongoDBLayer) AddComment(comment db.Comment, id string) ([]byte, error) {
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

func (m *MongoDBLayer) FindProduct(productID string) (db.Product, error) {

	s := m.session.Copy()
	defer s.Close()

	p := db.Product{}
	if err := s.DB(DB).C(PRODUCTS).Find(bson.M{"_id": bson.ObjectIdHex(productID)}).One(&p); err != nil {
		return p, fmt.Errorf("cannot find product %s, err: %v", productID, err)
	}

	return p, nil
}

func (m *MongoDBLayer) FindProductsByVendor(vendorID string) ([]db.Product, error) {
	s := m.session.Copy()
	defer s.Close()

	p := []db.Product{}
	if !bson.ObjectId(vendorID).Valid() {
		return p, fmt.Errorf("invalid vendorID as %s", vendorID)
	}

	mID := bson.M{
		"$match": bson.M{"vendor._id": bson.ObjectIdHex(vendorID)},
	}
	pipeLine := []bson.M{mID}

	if err := s.DB(DB).C(PRODUCTS).Pipe(pipeLine).All(&p); err != nil {
		return p, fmt.Errorf("cannot find products %s, err: %v", vendorID, err)
	}

	return p, nil
}

func (m *MongoDBLayer) GetProducts() ([]db.Product, error) {
	s := m.session.Copy()
	defer s.Close()
	products := []db.Product{}
	err := s.DB(DB).C(PRODUCTS).Find(nil).All(&products)
	return products, err
}

func (m *MongoDBLayer) GetComments() ([]db.Comment, error) {
	s := m.session.Copy()
	defer s.Close()
	comments := []db.Comment{}
	err := s.DB(DB).C(COMMENTS).Find(nil).All(&comments)
	return comments, err
}

func (m *MongoDBLayer) GetVendors() ([]db.Vendor, error) {
	s := m.session.Copy()
	defer s.Close()
	vendors := []db.Vendor{}
	err := s.DB(DB).C(VENDORS).Find(nil).All(&vendors)
	return vendors, err
}

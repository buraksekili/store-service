package mongo

import (
	"context"
	"fmt"

	"gopkg.in/mgo.v2/bson"

	"github.com/buraksekili/store-service/products"
	"github.com/pkg/errors"
	"gopkg.in/mgo.v2"
)

const (
	DB       = "store"
	PRODUCTS = "products"
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

	product.ID = bson.NewObjectId().Hex()
	return product.ID, s.DB(DB).C(PRODUCTS).Insert(&product)
}

func (m *MongoDBLayer) GetProduct(_ context.Context, productID string) (products.Product, error) {
	s := m.session.Copy()
	defer s.Close()

	p := products.Product{}
	if err := s.DB(DB).C(PRODUCTS).Find(bson.M{"_id": productID}).One(&p); err != nil {
		return p, errors.Wrap(fmt.Errorf("cannot find product %s", productID), err.Error())
	}
	return p, nil
}

func (m *MongoDBLayer) GetProductByVendorID(_ context.Context, vendorID string) ([]products.Product, error) {
	s := m.session.Copy()
	defer s.Close()

	p := []products.Product{}
	mID := bson.M{"$match": bson.M{"vendor._id": vendorID}}
	pipeLine := []bson.M{mID}
	if err := s.DB(DB).C(PRODUCTS).Pipe(pipeLine).All(&p); err != nil {
		return []products.Product{}, errors.Wrap(fmt.Errorf("cannot find products of %s", vendorID), err.Error())
	}
	return p, nil
}

func (m *MongoDBLayer) GetProducts(_ context.Context, offset, limit int) (products.ProductPage, error) {
	s := m.session.Copy()
	defer s.Close()

	listProducts := []products.Product{}
	err := s.DB(DB).C(PRODUCTS).Find(nil).Sort("name").Skip(offset * limit).Limit(limit).All(&listProducts)
	if err != nil {
		return products.ProductPage{}, errors.Wrap(fmt.Errorf("cannot run a query to obtain a list of product"), err.Error())
	}
	count, err := s.DB(DB).C(PRODUCTS).Find(nil).Count()
	if err != nil {
		return products.ProductPage{}, errors.Wrap(fmt.Errorf("cannot run a query to obtain the count of products"), err.Error())
	}
	return products.ProductPage{
		Total:    count,
		Offset:   offset,
		Limit:    limit,
		Products: listProducts,
	}, err
}

func (m *MongoDBLayer) GetComments(ctx context.Context, offset, limit int) (products.CommentPage, error) {
	s := m.session.Copy()
	defer s.Close()
	comments := []products.Comment{}
	err := s.DB(DB).C(COMMENTS).Find(nil).Sort("created_at").Skip(offset * limit).Limit(limit).All(&comments)
	if err != nil {
		return products.CommentPage{}, errors.Wrap(fmt.Errorf("cannot run a query to fetch all comments"), err.Error())
	}
	count, err := s.DB(DB).C(COMMENTS).Find(nil).Count()
	if err != nil {
		return products.CommentPage{}, errors.Wrap(fmt.Errorf("cannot run a query to fetch the count of the comments"), err.Error())
	}
	return products.CommentPage{
		Total:    count,
		Offset:   offset,
		Limit:    limit,
		Comments: comments,
	}, err
}

func (m *MongoDBLayer) AddComment(ctx context.Context, comment products.Comment, productID string) (string, error) {
	s := m.session.Copy()
	defer s.Close()

	err := s.DB(DB).C(PRODUCTS).Update(
		bson.M{"_id": productID},
		bson.M{"$push": bson.M{"comments": comment}},
	)
	if err != nil {
		return "", errors.Wrap(fmt.Errorf("cannot update product %s for comment", productID), err.Error())
	}
	if err = s.DB(DB).C(COMMENTS).Insert(comment); err != nil {
		return "", errors.Wrap(fmt.Errorf("cannot save comment %s", comment), err.Error())
	}
	return comment.ID, nil
}

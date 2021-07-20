package products

import (
	"context"
	"time"

	"github.com/buraksekili/store-service/users"

	"gopkg.in/mgo.v2/bson"
)

type Product struct {
	ID          bson.ObjectId `bson:"_id"`
	Name        string
	Category    string
	Description string
	ImageURL    string
	Price       float32
	Stock       int
	Comments    []Comment
	Vendor      users.Vendor
}

// ProductRepository is an interface for persistence.
type ProductRepository interface {
	// Save adds a new Product object.
	Save(ctx context.Context, product Product) (string, error)

	// FindByID finds a Product based on its string ID.
	FindByID(ctx context.Context, productID string) (Product, error)

	// FindByVendorID finds products of the Vendor
	// having a given string ID.
	FindByVendorID(ctx context.Context, offset, limit int, vendorID string) ([]Product, error)

	// GetProducts fetches all the products.
	GetProducts(ctx context.Context, offset, limit int) ([]Product, error)
}

type Comment struct {
	ID      bson.ObjectId `bson:"_id"`
	Owner   string
	Content string
	Date    time.Time
}

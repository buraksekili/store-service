package products

import (
	"context"
	"time"

	"github.com/buraksekili/store-service/users"
)

type Product struct {
	ID          string       `json:"id" bson:"_id"`
	Name        string       `json:"name" bson:"name"`
	Category    string       `json:"category" bson:"category"`
	Description string       `json:"description" bson:"description"`
	ImageURL    string       `json:"image_url" bson:"image_url"`
	Price       float32      `json:"price" bson:"price"`
	Stock       int          `json:"stock" bson:"stock"`
	Comments    []Comment    `json:"comments" bson:"comments"`
	Vendor      users.Vendor `json:"vendor" bson:"vendor"`
}

// ProductRepository is an interface for persistence.
type ProductRepository interface {
	// Save adds a new Product object.
	Save(ctx context.Context, product Product) (string, error)

	// GetProduct fetches a Product based on its string ID.
	GetProduct(ctx context.Context, productID string) (Product, error)

	// GetProductByVendorID finds products of the Vendor
	// having a given string ID.
	GetProductByVendorID(ctx context.Context, vendorID string) ([]Product, error)

	// GetProducts fetches all the products.
	GetProducts(ctx context.Context, offset, limit int) (ProductPage, error)

	// GetComments fetches all the comments on the system.
	GetComments(ctx context.Context, offset, limit int) (CommentPage, error)

	// AddComment adds a new comment to the given product.
	AddComment(ctx context.Context, comment Comment, productID string) (string, error)
}

type Comment struct {
	ID        string    `json:"id" bson:"_id"`
	Owner     string    `json:"owner" bson:"owner"`
	Content   string    `json:"content" bson:"content"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
}

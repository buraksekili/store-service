package products

import (
	"context"
	"errors"

	"github.com/buraksekili/store-service/users"
)

var (
	// ErrMalformedEntity shows that given entity through HTTP or any other
	// method is invalid to perform on it.
	ErrMalformedEntity = errors.New("malformed entity specification")

	// ErrVendorExistence shows that vendor is not found.
	ErrVendorExistence = errors.New("vendor does not exist")
)

// ProductPage represents a page of products to help navigation.
type ProductPage struct {
	Total    int
	Offset   int
	Limit    int
	Products []Product
}

// CommentPage represents a page of comments to help navigation.
type CommentPage struct {
	Total    int
	Offset   int
	Limit    int
	Comments []Comment
}

// ProductService represents an interface for the products service domain
// related functionalities.
type ProductService interface {
	// AddProduct adds a new product. Returns new product's
	// ID and the error if exists.
	AddProduct(ctx context.Context, product Product) (string, error)

	// GetProduct fetches a product based on ID of the product.
	GetProduct(ctx context.Context, productID string) (Product, error)

	// ListProducts lists all products.
	ListProducts(ctx context.Context, offset, limit int) (ProductPage, error)

	// ListVendorProducts lists all products of the given Vendor.
	ListVendorProducts(ctx context.Context, vendorID string) ([]Product, error)

	// GetComments fetches all comments on the system.
	GetComments(ctx context.Context, offset, limit int) (CommentPage, error)
}

type productsService struct {
	products ProductRepository
	users    users.UserRepository
}

// New returns a new products service.
func New(products ProductRepository, users users.UserRepository) ProductService {
	return productsService{products, users}
}

// AddProduct adds a new product. Returns new product's ID and the error if exists.
func (ps productsService) AddProduct(ctx context.Context, product Product) (string, error) {
	v, err := ps.users.GetVendorByID(ctx, product.Vendor.ID)
	if err != nil {
		return "", err
	}
	if v.ID == "" {
		return "", ErrVendorExistence
	}
	return ps.products.Save(ctx, product)
}

// GetProduct fetches a product based on ID of the product.
func (ps productsService) GetProduct(ctx context.Context, productID string) (Product, error) {
	return ps.products.GetProduct(ctx, productID)
}

// ListProducts lists all products.
func (ps productsService) ListProducts(ctx context.Context, offset, limit int) (ProductPage, error) {
	return ps.products.GetProducts(ctx, offset, limit)
}

// ListVendorProducts lists all products of the given Vendor.
func (ps productsService) ListVendorProducts(ctx context.Context, vendorID string) ([]Product, error) {
	return ps.products.GetProductByVendorID(ctx, vendorID)
}

// GetComments fetches all comments on the system.
func (ps productsService) GetComments(ctx context.Context, offset, limit int) (CommentPage, error) {
	return ps.products.GetComments(ctx, offset, limit)
}

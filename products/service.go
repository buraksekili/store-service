package products

import (
	"context"
	"errors"
)

var (
	ErrMalformedEntity = errors.New("malformed entity specification")
)

type ProductPage struct {
	Total    int
	Offset   int
	Limit    int
	Products []Product
}

type CommentPage struct {
	Total    int
	Offset   int
	Limit    int
	Comments []Comment
}

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
}

// New returns a new products service.
func New(products ProductRepository) ProductService {
	return productsService{products}
}

func (ps productsService) AddProduct(ctx context.Context, product Product) (string, error) {
	return ps.products.Save(ctx, product)
}

func (ps productsService) GetProduct(ctx context.Context, productID string) (Product, error) {
	return ps.products.GetProduct(ctx, productID)
}

func (ps productsService) ListProducts(ctx context.Context, offset, limit int) (ProductPage, error) {
	return ps.products.GetProducts(ctx, offset, limit)
}

func (ps productsService) ListVendorProducts(ctx context.Context, vendorID string) ([]Product, error) {
	return ps.products.GetProductByVendorID(ctx, vendorID)
}

func (ps productsService) GetComments(ctx context.Context, offset, limit int) (CommentPage, error) {
	return ps.products.GetComments(ctx, offset, limit)
}

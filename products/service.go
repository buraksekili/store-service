package products

import (
	"context"
	"errors"
)

var (
	// ErrMalformedEntity indicates incorrect entity.
	ErrMalformedEntity = errors.New("malformed entity specification")
)

type ProductService interface {
	// AddProduct adds a new product. Returns new product's
	// ID and the error if exists.
	AddProduct(ctx context.Context, product Product) (string, error)

	// GetProduct fetches a product based on ID of the product.
	GetProduct(ctx context.Context, productID string) (Product, error)

	// ListProducts lists all products.
	ListProducts(ctx context.Context, offset, limit int) ([]Product, error)

	// ListVendorProducts lists all products of the given Vendor.
	ListVendorProducts(ctx context.Context, offset, limit int, vendorID string) ([]Product, error)
}

type productsService struct {
	products ProductRepository
}

func New(products ProductRepository) ProductService {
	return productsService{products}
}

func (ps productsService) AddProduct(ctx context.Context, product Product) (string, error) {
	return ps.products.Save(ctx, product)
}

func (ps productsService) GetProduct(ctx context.Context, productID string) (Product, error) {
	return ps.products.FindByID(ctx, productID)
}

func (ps productsService) ListProducts(ctx context.Context, offset, limit int) ([]Product, error) {
	return ps.products.GetProducts(ctx, offset, limit)
}

func (ps productsService) ListVendorProducts(ctx context.Context, offset, limit int, vendorID string) ([]Product, error) {
	return ps.products.FindByVendorID(ctx, offset, limit, vendorID)
}

package api

import (
	"context"
	"fmt"

	"github.com/buraksekili/store-service/users"

	"gopkg.in/mgo.v2/bson"

	"github.com/pkg/errors"

	"github.com/buraksekili/store-service/products"
	"github.com/go-kit/kit/endpoint"
)

func addProductEndpoint(svc products.ProductService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(addProductReq)
		p, err := extractProduct(req)
		if err != nil {
			return addProductRes{}, errors.Wrap(err, products.ErrMalformedEntity.Error())
		}
		pid, err := svc.AddProduct(ctx, p)
		if err != nil {
			return addProductRes{}, errors.Wrap(err, products.ErrMalformedEntity.Error())
		}
		return addProductRes{ID: pid}, nil
	}
}

func extractProduct(r addProductReq) (products.Product, error) {
	vendorID := bson.ObjectId(r.VendorID)
	if !vendorID.Valid() {
		return products.Product{}, fmt.Errorf("invalid vendor id, got: %s", vendorID)
	}
	return products.Product{
		ID:          bson.NewObjectId(),
		Name:        r.Name,
		Category:    r.Category,
		Description: r.Description,
		ImageURL:    r.ImageURL,
		Price:       r.Price,
		Stock:       r.Stock,
		Comments:    []products.Comment{},
		Vendor:      users.Vendor{ID: vendorID},
	}, nil
}

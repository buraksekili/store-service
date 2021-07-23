package api

import (
	"context"

	"github.com/buraksekili/store-service/users"

	"github.com/pkg/errors"

	"github.com/buraksekili/store-service/products"
	"github.com/go-kit/kit/endpoint"
)

func addProductEndpoint(svc products.ProductService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(addProductReq)
		p, err := extractProduct(req)
		if err != nil {
			return nil, errors.Wrap(err, products.ErrMalformedEntity.Error())
		}
		pid, err := svc.AddProduct(ctx, p)
		if err != nil {
			return nil, errors.Wrap(err, products.ErrMalformedEntity.Error())
		}
		return addProductRes{ID: pid}, nil
	}
}

func getProductsEndpoint(svc products.ProductService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(getProductsReq)
		res, err := svc.ListProducts(ctx, req.offset, req.limit)
		if err != nil {
			return listProductsRes{}, err
		}
		return listProductsRes{Total: res.Total, Offset: res.Offset, Limit: res.Limit, Products: res.Products}, nil
	}
}

func getProductEndpoint(svc products.ProductService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(getProductReq)
		p, err := svc.GetProduct(ctx, req.ProductID)
		if err != nil {
			return getProductRes{}, err
		}
		return getProductRes{Product: p}, err
	}
}

func vendorsProductEndpoint(svc products.ProductService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(listVendorProductsReq)
		prods, err := svc.ListVendorProducts(ctx, req.VendorID)
		if err != nil {
			return listVendorProductsRes{}, err
		}
		return prods, nil
	}
}

func getAllCommentsEndpoint(svc products.ProductService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(getAllCommentsReq)
		cp, err := svc.GetComments(ctx, req.offset, req.limit)
		if err != nil {
			return getAllCommentsRes{}, err
		}
		return getAllCommentsRes{
			Total:    cp.Total,
			Offset:   cp.Offset,
			Limit:    cp.Limit,
			Comments: cp.Comments,
		}, nil
	}
}

func getProductCommentsEndpoint(svc products.ProductService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(getProductCommentsReq)
		comments, err := svc.GetProduct(ctx, req.ProductID)
		if err != nil {
			return getProductCommentsRes{}, nil
		}
		return getProductCommentsRes{Comments: comments.Comments}, nil
	}
}

func extractProduct(r addProductReq) (products.Product, error) {
	return products.Product{
		Name:        r.Name,
		Category:    r.Category,
		Description: r.Description,
		ImageURL:    r.ImageURL,
		Price:       r.Price,
		Stock:       r.Stock,
		Vendor:      users.Vendor{ID: r.VendorID},
	}, nil
}

package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/buraksekili/store-service/pkg/logger"

	"github.com/buraksekili/store-service/users"
	"github.com/pkg/errors"

	urlhelper "github.com/buraksekili/store-service/pkg/url"

	"github.com/gorilla/mux"

	"github.com/buraksekili/store-service/products"
	httptransport "github.com/go-kit/kit/transport/http"
)

func MakeHTTPHandler(svc products.ProductService, logger logger.Logger) http.Handler {
	r := mux.NewRouter()

	options := []httptransport.ServerOption{
		// httptransport.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
		httptransport.ServerErrorEncoder(encodeError),
	}

	r.Methods("POST").Path("/products").Handler(httptransport.NewServer(
		addProductEndpoint(svc),
		decodeAddProductReq,
		encodeResponse,
		options...,
	))

	r.Methods("GET").Path("/products").Handler(httptransport.NewServer(
		getProductsEndpoint(svc),
		decodeGetProductsReq,
		encodeResponse,
		options...,
	))

	r.Methods("GET").Path("/products/{product_id}").Handler(httptransport.NewServer(
		getProductEndpoint(svc),
		decodeGetProductReq,
		encodeResponse,
		options...,
	))

	r.Methods("GET").Path("/products/vendors/{vendor_id}").Handler(httptransport.NewServer(
		vendorsProductEndpoint(svc),
		decodeVendorProductsReq,
		encodeResponse,
		options...,
	))

	r.Methods("GET").Path("/products/comments/").Handler(httptransport.NewServer(
		getAllCommentsEndpoint(svc),
		decodeGetAllCommentsReq,
		encodeResponse,
		options...,
	))

	r.Methods("GET").Path("/products/comments/{product_id}").Handler(httptransport.NewServer(
		getProductCommentsEndpoint(svc),
		decodeGetProductCommentsReq,
		encodeResponse,
		options...,
	))

	return r
}

func decodeAddProductReq(_ context.Context, r *http.Request) (interface{}, error) {
	var req addProductReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return addProductReq{}, err
	}
	return req, nil
}

func decodeGetProductsReq(ctx context.Context, r *http.Request) (interface{}, error) {
	offset, limit, err := parseOffsetLimitQueryParams(r)
	if err != nil {
		return getProductsReq{}, err
	}
	return getProductsReq{offset, limit}, nil
}

func decodeGetProductReq(ctx context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	pid, ok := vars["product_id"]
	if !ok {
		return getProductReq{}, users.ErrInvalidRequestPath
	}
	return getProductReq{pid}, nil
}

func decodeVendorProductsReq(ctx context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	vid, ok := vars["vendor_id"]
	if !ok {
		return listVendorProductsReq{}, users.ErrInvalidRequestPath
	}
	return listVendorProductsReq{VendorID: vid}, nil
}

func decodeGetAllCommentsReq(ctx context.Context, r *http.Request) (interface{}, error) {
	offset, limit, err := parseOffsetLimitQueryParams(r)
	if err != nil {
		return getAllCommentsReq{}, err
	}
	return getAllCommentsReq{offset: offset, limit: limit}, nil
}

func decodeGetProductCommentsReq(ctx context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	pid, ok := vars["product_id"]
	if !ok {
		return getProductCommentsReq{}, users.ErrInvalidRequestPath
	}
	return getProductCommentsReq{pid}, nil
}

func parseOffsetLimitQueryParams(r *http.Request) (offset, limit int, err error) {
	offset, err = urlhelper.ParseIntQueryParams("offset", 0, r)
	if err != nil {
		return
	}
	if offset < 0 {
		return 0, 0, errors.Wrap(fmt.Errorf("offset cannot be smaller than 0"), users.ErrInvalidRequestPath.Error())
	}

	limit, err = urlhelper.ParseIntQueryParams("limit", 10, r)
	if err != nil {
		return
	}
	if limit < 0 {
		return 0, 0, errors.Wrap(fmt.Errorf("limit cannot be smaller than 0"), users.ErrInvalidRequestPath.Error())
	}
	return
}

type errorer interface {
	error() error
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(errorer); ok && e.error() != nil {
		encodeError(ctx, e.error(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(codeFrom(err))
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}

func codeFrom(err error) int {
	switch err {
	case products.ErrMalformedEntity:
		return http.StatusBadRequest
	case users.ErrInvalidRequestPath:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}

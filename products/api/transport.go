package api

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-kit/kit/transport"

	"github.com/go-kit/kit/log"
	"github.com/gorilla/mux"

	"github.com/buraksekili/store-service/products"
	httptransport "github.com/go-kit/kit/transport/http"
)

func MakeHTTPHandler(svc products.ProductService, logger log.Logger) http.Handler {
	r := mux.NewRouter()

	options := []httptransport.ServerOption{
		httptransport.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
		httptransport.ServerErrorEncoder(encodeError),
	}

	r.Methods("POST").Path("/products").Handler(httptransport.NewServer(
		addProductEndpoint(svc),
		decodeAddProductReq,
		encodeResponse,
		options...,
	))
	return r
}

// errorer is implemented by all concrete response types that may contain
// errors. It allows us to change the HTTP response code without needing to
// trigger an endpoint (transport-level) error. For more information, read the
// big comment in endpoints.go.
type errorer interface {
	error() error
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(errorer); ok && e.error() != nil {
		// Not a Go kit transport error, but a business-logic error.
		// Provide those as HTTP errors.
		encodeError(ctx, e.error(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

func decodeAddProductReq(_ context.Context, r *http.Request) (interface{}, error) {
	var req addProductReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	return req, nil
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	if err == nil {
		panic("encodeError with nil error")
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(codeFrom(err))
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}

// TODO: implement errors
func codeFrom(err error) int {
	switch err {
	default:
		return http.StatusInternalServerError
	}
}

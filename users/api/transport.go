package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	urlhelper "github.com/buraksekili/store-service/pkg/url"

	"github.com/buraksekili/store-service/pkg/logger"

	"github.com/pkg/errors"

	"github.com/buraksekili/store-service/users"
	httptransport "github.com/go-kit/kit/transport/http"

	"github.com/gorilla/mux"
)

func MakeHTTPHandler(svc users.UserService, logger logger.Logger) http.Handler {
	r := mux.NewRouter()

	options := []httptransport.ServerOption{
		httptransport.ServerErrorEncoder(encodeError),
	}

	r.Methods("GET").Path("/users").Handler(httptransport.NewServer(
		getUsersEndpoint(svc),
		decodeGetUsersReq,
		encodeResponse,
		options...,
	))

	r.Methods("GET").Path("/users/{user_id}").Handler(httptransport.NewServer(
		getUserEndpoint(svc),
		decodeGetUserReq,
		encodeResponse,
		options...,
	))

	r.Methods("POST").Path("/users/signup").Handler(httptransport.NewServer(
		addUserEndpoint(svc),
		decodeAddUserReq,
		encodeResponse,
		options...,
	))

	r.Methods("POST").Path("/users/login").Handler(httptransport.NewServer(
		loginEndpoint(svc),
		decodeAddUserReq,
		encodeResponse,
		options...,
	))

	r.Methods("POST").Path("/vendors").Handler(httptransport.NewServer(
		addVendorEndpoint(svc),
		decodeAddVendorReq,
		encodeResponse,
		options...,
	))

	r.Methods("GET").Path("/vendors").Handler(httptransport.NewServer(
		getVendorsEndpoint(svc),
		decodeGetUsersReq,
		encodeResponse,
		options...,
	))

	return r
}

func decodeGetUsersReq(ctx context.Context, r *http.Request) (interface{}, error) {
	offset, err := urlhelper.ParseIntQueryParams("offset", 0, r)
	if err != nil {
		return nil, err
	}
	if offset < 0 {
		return nil, errors.Wrap(fmt.Errorf("offset cannot be smaller than 0"), users.ErrInvalidRequestPath.Error())
	}

	limit, err := urlhelper.ParseIntQueryParams("limit", 10, r)
	if err != nil {
		return nil, err
	}
	if limit < 0 {
		return nil, errors.Wrap(fmt.Errorf("limit cannot be smaller than 0"), users.ErrInvalidRequestPath.Error())
	}

	return getUsersReq{offset, limit}, nil
}

func decodeAddUserReq(ctx context.Context, r *http.Request) (interface{}, error) {
	if !strings.Contains(r.Header.Get("Content-Type"), "application/json") {
		return nil, users.ErrUnsupportedContentType
	}
	var user users.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		return nil, errors.Wrap(users.ErrMalformedEntity, err.Error())
	}
	return addUserReq{User: user}, nil
}

func decodeGetUserReq(ctx context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	uID, ok := vars["user_id"]
	if !ok {
		return nil, users.ErrInvalidRequestPath
	}
	return getUserReq{uID}, nil
}

func decodeAddVendorReq(ctx context.Context, r *http.Request) (interface{}, error) {
	if !strings.Contains(r.Header.Get("Content-Type"), "application/json") {
		return nil, users.ErrUnsupportedContentType
	}
	var vendor users.Vendor
	if err := json.NewDecoder(r.Body).Decode(&vendor); err != nil {
		return nil, errors.Wrap(users.ErrMalformedEntity, err.Error())
	}
	return addVendorReq{vendor}, nil
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
	if err == nil {
		panic("encodeError with nil error")
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(codeFrom(err))
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}

func codeFrom(err error) int {
	switch err {
	case users.ErrUnauthorized:
		return http.StatusUnauthorized
	case users.ErrMalformedEntity:
		return http.StatusBadRequest
	case users.ErrInvalidRequestPath:
		return http.StatusBadRequest
	case users.ErrUnsupportedContentType:
		return http.StatusUnsupportedMediaType
	default:
		return http.StatusInternalServerError
	}
}

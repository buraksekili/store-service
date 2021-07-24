package api

import (
	"context"
	"strings"

	"github.com/buraksekili/store-service/users"
	"github.com/go-kit/kit/endpoint"
)

func getUsersEndpoint(svc users.UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(getUsersReq)
		u, err := svc.GetUsers(ctx, req.offset, req.limit)
		if err != nil {
			return nil, err
		}
		return getUsersRes{
			Total:  u.Total,
			Offset: u.Offset,
			Limit:  u.Limit,
			Users:  u.Users,
		}, nil
	}
}

func addUserEndpoint(svc users.UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(addUserReq)
		uid, err := svc.AddUser(ctx, req.User)
		if err != nil {
			return nil, err
		}
		return addUserRes{uid}, nil
	}
}

func getUserEndpoint(svc users.UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(getUserReq)
		user, err := svc.GetUser(ctx, req.UserID)
		if err != nil {
			return nil, err
		}
		user.Password = ""
		return getUserRes{user.ID, user.Username, user.Email}, nil
	}
}

func loginEndpoint(svc users.UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(addUserReq)
		if valid := validate(req.User); !valid {
			return nil, users.ErrMalformedEntity
		}
		uID, err := svc.Login(ctx, req.User)
		if err != nil {
			return nil, err
		}
		return addUserRes{uID}, nil
	}
}

func addVendorEndpoint(svc users.UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(addVendorReq)
		if valid := validate(req.Vendor); !valid {
			return nil, users.ErrMalformedEntity
		}
		vid, err := svc.AddVendor(ctx, req.Vendor)
		if err != nil {
			return nil, err
		}
		return addVendorRes{vid}, nil
	}
}

func getVendorsEndpoint(svc users.UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(getUsersReq)
		u, err := svc.GetVendors(ctx, req.offset, req.limit)
		if err != nil {
			return nil, err
		}
		return getVendorsRes{
			Total:   u.Total,
			Offset:  u.Offset,
			Limit:   u.Limit,
			Vendors: u.Vendors,
		}, nil
	}
}

func validate(i interface{}) bool {
	switch v := i.(type) {
	case users.User:
		if len(strings.TrimSpace(v.Email)) == 0 ||
			len(strings.TrimSpace(v.Password)) == 0 {
			return false
		}
		return true
	case users.Vendor:
		if len(strings.TrimSpace(v.Description)) == 0 ||
			len(strings.TrimSpace(v.Name)) == 0 {
			return false
		}
		return true
	}
	return false
}

package api

import "github.com/buraksekili/store-service/users"

type getUsersReq struct {
	offset int64
	limit  int64
}

type addUserReq struct {
	User users.User
}

type getUserReq struct {
	UserID string
}

type addVendorReq struct {
	Vendor users.Vendor
}

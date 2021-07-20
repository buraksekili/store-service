package api

import "github.com/buraksekili/store-service/users"

type getUsersRes struct {
	Total  int          `json:"total"`
	Offset int          `json:"offset"`
	Limit  int          `json:"limit"`
	Users  []users.User `json:"users"`
}

type getVendorsRes struct {
	Total   int            `json:"total"`
	Offset  int            `json:"offset"`
	Limit   int            `json:"limit"`
	Vendors []users.Vendor `json:"vendors"`
}

type addUserRes struct {
	UserID string `json:"id"`
}

type getUserRes struct {
	UserID   string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type addVendorRes struct {
	VendorID string `json:"id"`
}

package users

import (
	"context"
)

type User struct {
	ID       string `json:"id" bson:"_id"`
	Username string `json:"username" bson:"username"`
	Email    string `json:"email" bson:"email"`
	Password string `json:",omitempty" bson:"password"`
}

type UserRepository interface {
	CreateUser(ctx context.Context, user User) (string, error)
	ListUsers(ctx context.Context, offset, limit int64) (UserPage, error)
	GetUser(ctx context.Context, userID string) (User, error)
	GetUserByEmail(ctx context.Context, userEmail string) (User, error)

	CreateVendor(ctx context.Context, vendor Vendor) (string, error)
	ListVendors(ctx context.Context, offset, limit int64) (VendorPage, error)
	GetVendorByName(ctx context.Context, vendorName string) (Vendor, error)
}

type Vendor struct {
	ID          string `json:"id" bson:"_id"`
	Name        string `json:"name" bson:"name"`
	Description string `json:"description" bson:"description"`
}

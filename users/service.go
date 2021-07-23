package users

import (
	"context"
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

var (
	ErrMalformedEntity        = errors.New("malformed entity specification")
	ErrUnauthorized           = errors.New("invalid credentials provided")
	ErrUnsupportedContentType = errors.New("unsupported content type")
	ErrInvalidRequestPath     = errors.New("invalid request path provided")
)

// UserPage represents a response for fetching multiple Users.
// Instead of returning all users at once, UserPage includes
// Total, Offset and Limit values to ease pagination of the response.
type UserPage struct {
	Total  int
	Offset int
	Limit  int
	Users  []User
}

// VendorPage represents a response for fetching multiple Vendors.
// Instead of returning all vendors at once, VendorPage includes
// Total, Offset and Limit values to ease pagination of the response.
type VendorPage struct {
	Total   int
	Offset  int
	Limit   int
	Vendors []Vendor
}

type UserService interface {
	// AddUser adds a new user. Returns new user's
	// ID and the error if exists.
	AddUser(ctx context.Context, user User) (string, error)

	// GetUser fetches a user based on ID of the user.
	GetUser(ctx context.Context, userID string) (User, error)

	// GetUsers fetches all users.
	GetUsers(ctx context.Context, offset, limit int) (UserPage, error)

	// Login authenticates the given user. Returns non-nil
	// error if the authentication fails.
	Login(ctx context.Context, user User) (string, error)

	// AddVendor adds a new vendor. Returns new vendor's
	// ID and the error if exists.
	AddVendor(ctx context.Context, vendor Vendor) (string, error)

	// GetVendor retrieves a vendor based on ID.
	GetVendor(ctx context.Context, vendorID string) (Vendor, error)

	// GetVendors fetches all vendors.
	GetVendors(ctx context.Context, offset, limit int) (VendorPage, error)
}

type usersService struct {
	users UserRepository
}

func New(users UserRepository) UserService {
	return usersService{users}
}

func (us usersService) AddUser(ctx context.Context, user User) (string, error) {
	u, _ := us.users.GetUserByEmail(ctx, user.Email)
	if u.Email == user.Email || u.Username == user.Email {
		return "", fmt.Errorf("user credentials already taken")
	}
	return us.users.CreateUser(ctx, user)
}

func (us usersService) GetUser(ctx context.Context, userID string) (User, error) {
	return us.users.GetUser(ctx, userID)
}

func (us usersService) GetUsers(ctx context.Context, offset, limit int) (UserPage, error) {
	return us.users.ListUsers(ctx, offset, limit)
}

func (us usersService) Login(ctx context.Context, user User) (string, error) {
	u, err := us.users.GetUserByEmail(ctx, user.Email)
	if err != nil {
		return "", errors.Wrap(err, ErrUnauthorized.Error())
	}
	if u.Password != user.Password ||
		u.Email != user.Email {
		return "", ErrUnauthorized
	}
	return u.ID, nil
}

func (us usersService) AddVendor(ctx context.Context, vendor Vendor) (string, error) {
	v, _ := us.users.GetVendorByName(ctx, vendor.Name)
	if strings.TrimSpace(v.Name) == strings.TrimSpace(vendor.Name) {
		return "", fmt.Errorf("vendor credentials already taken")
	}
	return us.users.CreateVendor(ctx, vendor)
}

func (us usersService) GetVendors(ctx context.Context, offset, limit int) (VendorPage, error) {
	return us.users.ListVendors(ctx, offset, limit)
}

func (us usersService) GetVendor(ctx context.Context, vendorID string) (Vendor, error) {
	return us.users.GetVendorByID(ctx, vendorID)
}

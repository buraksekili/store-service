package mongo

import (
	"context"
	"fmt"

	"github.com/pkg/errors"

	"github.com/buraksekili/store-service/users"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	DB      = "store"
	VENDORS = "vendors"
	USERS   = "users"
)

type MongoDBLayer struct {
	session *mgo.Session
}

func NewMongoDBLayer(url string) (users.UserRepository, error) {
	s, err := mgo.Dial(url)
	return &MongoDBLayer{s}, err
}

func (m *MongoDBLayer) CreateUser(_ context.Context, user users.User) (string, error) {
	s := m.session.Copy()
	defer s.Close()
	user.ID = bson.NewObjectId().Hex()
	return user.ID, s.DB(DB).C(USERS).Insert(user)
}

func (m *MongoDBLayer) ListUsers(_ context.Context, offset, limit int64) (users.UserPage, error) {
	s := m.session.Copy()
	defer s.Close()

	listUser := []users.User{}
	fields := bson.M{"_id": 1, "username": 1, "email": 1}
	err := s.DB(DB).C(USERS).Find(nil).Sort("email").Skip(int(offset * limit)).Limit(int(limit)).Select(fields).All(&listUser)
	if err != nil {
		return users.UserPage{}, errors.Wrap(fmt.Errorf("cannot run a query to obtain a list of users"), err.Error())
	}
	count, err := s.DB(DB).C(USERS).Find(nil).Count()
	if err != nil {
		return users.UserPage{}, errors.Wrap(fmt.Errorf("cannot run a query to obtain total count"), err.Error())
	}

	return users.UserPage{
		Total:  int64(count),
		Offset: offset,
		Limit:  limit,
		Users:  listUser,
	}, nil
}

func (m *MongoDBLayer) GetUser(_ context.Context, userID string) (users.User, error) {
	s := m.session.Copy()
	defer s.Close()

	u := users.User{}
	fields := bson.M{"_id": 1, "username": 1, "email": 1}
	if err := s.DB(DB).C(USERS).Find(bson.M{"_id": userID}).Select(fields).One(&u); err != nil {
		return u, fmt.Errorf("cannot find user %s, err: %v", userID, err)
	}
	return u, nil
}

func (m *MongoDBLayer) GetUserByEmail(_ context.Context, userEmail string) (users.User, error) {
	s := m.session.Copy()
	defer s.Close()

	u := users.User{}
	if err := s.DB(DB).C(USERS).Find(bson.M{"email": userEmail}).One(&u); err != nil {
		return u, fmt.Errorf("cannot find user %s, err: %v", userEmail, err)
	}
	return u, nil
}

func (m *MongoDBLayer) CreateVendor(_ context.Context, vendor users.Vendor) (string, error) {
	s := m.session.Copy()
	defer s.Close()
	vendor.ID = bson.NewObjectId().Hex()
	return vendor.ID, s.DB(DB).C(VENDORS).Insert(&vendor)
}

func (m *MongoDBLayer) ListVendors(_ context.Context, offset, limit int64) (users.VendorPage, error) {
	s := m.session.Copy()
	defer s.Close()

	vendors := []users.Vendor{}
	if err := s.DB(DB).C(VENDORS).Find(nil).All(&vendors); err != nil {
		return users.VendorPage{}, errors.Wrap(fmt.Errorf("cannot run a query to obtain a list of vendors"), err.Error())
	}
	count, err := s.DB(DB).C(VENDORS).Find(nil).Count()
	if err != nil {
		return users.VendorPage{}, errors.Wrap(fmt.Errorf("cannot run a query to obtain the total count"), err.Error())
	}
	return users.VendorPage{
		Total:   int64(count),
		Offset:  offset,
		Limit:   limit,
		Vendors: vendors,
	}, err
}

func (m *MongoDBLayer) GetVendorByName(ctx context.Context, vendorName string) (users.Vendor, error) {
	s := m.session.Copy()
	defer s.Close()

	v := users.Vendor{}
	if err := s.DB(DB).C(VENDORS).Find(bson.M{"name": vendorName}).One(&v); err != nil {
		return v, fmt.Errorf("cannot find vendor %s, err: %v", vendorName, err)
	}
	return v, nil
}

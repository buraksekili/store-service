package api

import (
	"context"
	"fmt"
	"time"

	"github.com/buraksekili/store-service/pkg/logger"

	"github.com/buraksekili/store-service/users"
)

type loggingMiddleware struct {
	logger logger.Logger
	svc    users.UserService
}

func LoggingMiddleware(svc users.UserService, logger logger.Logger) users.UserService {
	return &loggingMiddleware{logger, svc}
}

func (lm loggingMiddleware) AddUser(ctx context.Context, user users.User) (uid string, err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method `signup` for user %s took %s to complete", user.Email, time.Since(begin))
		if err != nil {
			lm.logger.Error(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}
		lm.logger.Info(fmt.Sprintf("%s without errors.", message))
	}(time.Now())
	return lm.svc.AddUser(ctx, user)
}

func (lm loggingMiddleware) GetUser(ctx context.Context, userID string) (user users.User, err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method `get_user` for user %s took %s to complete", user.Username, time.Since(begin))
		if err != nil {
			lm.logger.Error(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}
		lm.logger.Info(fmt.Sprintf("%s without errors.", message))
	}(time.Now())
	return lm.svc.GetUser(ctx, userID)
}

func (lm loggingMiddleware) GetUsers(ctx context.Context, offset, limit int64) (up users.UserPage, err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method `get_users` took %s to complete", time.Since(begin))
		if err != nil {
			lm.logger.Error(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}
		lm.logger.Info(fmt.Sprintf("%s without errors.", message))

	}(time.Now())
	return lm.svc.GetUsers(ctx, offset, limit)
}

func (lm loggingMiddleware) Login(ctx context.Context, user users.User) (userID string, err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method `login` for user %s took %s to complete", user.Username, time.Since(begin))
		if err != nil {
			lm.logger.Error(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}
		lm.logger.Info(fmt.Sprintf("%s without errors.", message))

	}(time.Now())
	return lm.svc.Login(ctx, user)
}

func (lm loggingMiddleware) AddVendor(ctx context.Context, vendor users.Vendor) (vid string, err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method `add_vendor` for vendor %s took %s to complete", vendor.Name, time.Since(begin))
		if err != nil {
			lm.logger.Error(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}
		lm.logger.Info(fmt.Sprintf("%s without errors.", message))
	}(time.Now())
	return lm.svc.AddVendor(ctx, vendor)
}

func (lm loggingMiddleware) GetVendors(ctx context.Context, offset, limit int64) (vp users.VendorPage, err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method `get_vendors` took %s to complete", time.Since(begin))
		if err != nil {
			lm.logger.Error(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}
		lm.logger.Info(fmt.Sprintf("%s without errors.", message))

	}(time.Now())
	return lm.svc.GetVendors(ctx, offset, limit)
}

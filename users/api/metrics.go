package api

import (
	"context"
	"time"

	"github.com/buraksekili/store-service/users"
	"github.com/go-kit/kit/metrics"
)

type metricsMiddleware struct {
	requestCount   metrics.Counter
	requestLatency metrics.Histogram
	svc            users.UserService
}

func NewMetricsMiddleware(svc users.UserService, counter metrics.Counter, latency metrics.Histogram) users.UserService {
	return &metricsMiddleware{counter, latency, svc}
}

func (ms metricsMiddleware) AddUser(ctx context.Context, user users.User) (string, error) {
	defer func(begin time.Time) {
		ms.requestCount.With("method", "signup").Add(1)
		ms.requestLatency.With("method", "signup").Observe(time.Since(begin).Seconds())
	}(time.Now())
	return ms.svc.AddUser(ctx, user)
}

func (ms metricsMiddleware) GetUser(ctx context.Context, userID string) (users.User, error) {
	defer func(begin time.Time) {
		ms.requestCount.With("method", "get_user").Add(1)
		ms.requestLatency.With("method", "get_user").Observe(time.Since(begin).Seconds())
	}(time.Now())
	return ms.svc.GetUser(ctx, userID)
}

func (ms metricsMiddleware) GetUsers(ctx context.Context, offset, limit int64) (users.UserPage, error) {
	defer func(begin time.Time) {
		ms.requestCount.With("method", "get_users").Add(1)
		ms.requestLatency.With("method", "get_users").Observe(time.Since(begin).Seconds())
	}(time.Now())
	return ms.svc.GetUsers(ctx, offset, limit)
}

func (ms metricsMiddleware) Login(ctx context.Context, user users.User) (string, error) {
	defer func(begin time.Time) {
		ms.requestCount.With("method", "login").Add(1)
		ms.requestLatency.With("method", "login").Observe(time.Since(begin).Seconds())
	}(time.Now())
	return ms.svc.Login(ctx, user)
}

func (ms metricsMiddleware) AddVendor(ctx context.Context, vendor users.Vendor) (string, error) {
	defer func(begin time.Time) {
		ms.requestCount.With("method", "add_vendor").Add(1)
		ms.requestLatency.With("method", "add_vendor").Observe(time.Since(begin).Seconds())
	}(time.Now())
	return ms.svc.AddVendor(ctx, vendor)
}

func (ms metricsMiddleware) GetVendors(ctx context.Context, offset, limit int64) (users.VendorPage, error) {
	defer func(begin time.Time) {
		ms.requestCount.With("method", "get_vendors").Add(1)
		ms.requestLatency.With("method", "get_vendors").Observe(time.Since(begin).Seconds())
	}(time.Now())
	return ms.svc.GetVendors(ctx, offset, limit)

}

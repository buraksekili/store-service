package api

import (
	"context"
	"time"

	"github.com/buraksekili/store-service/products"
	"github.com/go-kit/kit/metrics"
)

type metricsMiddleware struct {
	requestCount   metrics.Counter
	requestLatency metrics.Histogram
	svc            products.ProductService
}

func NewMetricsMiddleware(svc products.ProductService, counter metrics.Counter, latency metrics.Histogram) products.ProductService {
	return &metricsMiddleware{counter, latency, svc}
}

func (ms metricsMiddleware) AddProduct(ctx context.Context, product products.Product) (s string, err error) {
	defer func(begin time.Time) {
		ms.requestCount.With("method", "add_product").Add(1)
		ms.requestLatency.With("method", "add_product").Observe(time.Since(begin).Seconds())
	}(time.Now())
	return ms.svc.AddProduct(ctx, product)
}

func (ms metricsMiddleware) GetProduct(ctx context.Context, productID string) (p products.Product, err error) {
	defer func(begin time.Time) {
		ms.requestCount.With("method", "get_product").Add(1)
		ms.requestLatency.With("method", "get_product").Observe(time.Since(begin).Seconds())
	}(time.Now())
	return ms.svc.GetProduct(ctx, productID)
}

func (ms metricsMiddleware) ListProducts(ctx context.Context, offset, limit int) (p products.ProductPage, err error) {
	defer func(begin time.Time) {
		ms.requestCount.With("method", "list_products").Add(1)
		ms.requestLatency.With("method", "list_products").Observe(time.Since(begin).Seconds())
	}(time.Now())
	return ms.svc.ListProducts(ctx, offset, limit)
}

func (ms metricsMiddleware) ListVendorProducts(ctx context.Context, vendorID string) (p []products.Product, err error) {
	defer func(begin time.Time) {
		ms.requestCount.With("method", "vendor_products").Add(1)
		ms.requestLatency.With("method", "vendor_products").Observe(time.Since(begin).Seconds())
	}(time.Now())
	return ms.svc.ListVendorProducts(ctx, vendorID)
}

func (ms metricsMiddleware) GetComments(ctx context.Context, offset, limit int) (pc products.CommentPage, err error) {
	defer func(begin time.Time) {
		ms.requestCount.With("method", "get_comments").Add(1)
		ms.requestLatency.With("method", "get_comments").Observe(time.Since(begin).Seconds())
	}(time.Now())
	return ms.svc.GetComments(ctx, offset, limit)
}

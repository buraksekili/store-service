package api

import (
	"context"
	"fmt"
	"time"

	"github.com/buraksekili/store-service/pkg/logger"
	"github.com/buraksekili/store-service/products"
)

type loggingMiddleware struct {
	logger logger.Logger
	svc    products.ProductService
}

func LoggingMiddleware(svc products.ProductService, logger logger.Logger) products.ProductService {
	return &loggingMiddleware{logger, svc}
}

func (lm loggingMiddleware) AddProduct(ctx context.Context, product products.Product) (s string, err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method `add_product` for product %s took %s to complete", product.Name, time.Since(begin))
		if err != nil {
			lm.logger.Error(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}
		lm.logger.Info(fmt.Sprintf("%s without errors.", message))
	}(time.Now())
	return lm.svc.AddProduct(ctx, product)
}

func (lm loggingMiddleware) GetProduct(ctx context.Context, productID string) (p products.Product, err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method `get_product` for product %s took %s to complete", productID, time.Since(begin))
		if err != nil {
			lm.logger.Error(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}
		lm.logger.Info(fmt.Sprintf("%s without errors.", message))
	}(time.Now())
	return lm.svc.GetProduct(ctx, productID)
}

func (lm loggingMiddleware) ListProducts(ctx context.Context, offset, limit int) (p products.ProductPage, err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method `list_products` took %s to complete", time.Since(begin))
		if err != nil {
			lm.logger.Error(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}
		lm.logger.Info(fmt.Sprintf("%s without errors.", message))
	}(time.Now())
	return lm.svc.ListProducts(ctx, offset, limit)
}

func (lm loggingMiddleware) ListVendorProducts(ctx context.Context, vendorID string) (p []products.Product, err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method `vendor_products` for vendor %s took %s to complete", vendorID, time.Since(begin))
		if err != nil {
			lm.logger.Error(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}
		lm.logger.Info(fmt.Sprintf("%s without errors.", message))
	}(time.Now())
	return lm.svc.ListVendorProducts(ctx, vendorID)
}

func (lm loggingMiddleware) GetComments(ctx context.Context, offset, limit int) (p products.CommentPage, err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method `get_comments` took %s to complete", time.Since(begin))
		if err != nil {
			lm.logger.Error(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}
		lm.logger.Info(fmt.Sprintf("%s without errors.", message))
	}(time.Now())
	return lm.svc.GetComments(ctx, offset, limit)
}

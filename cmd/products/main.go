package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	stdprometheus "github.com/prometheus/client_golang/prometheus"

	logpkg "github.com/buraksekili/store-service/pkg/logger"

	"github.com/buraksekili/store-service/products"
	"github.com/buraksekili/store-service/products/api"
	"github.com/buraksekili/store-service/products/persistence/mongo"
)

const productServiceURL = ":8181"

func main() {
	logger := logpkg.New()

	u := fmt.Sprintf("mongodb://mongo:27017")
	productsRepo, err := mongo.NewMongoDBLayer(u)
	if err != nil {
		fmt.Println("CANNOT CONNECT TO", u)
		os.Exit(1)
	}

	var svc products.ProductService
	svc = products.New(productsRepo)
	svc = api.LoggingMiddleware(svc, logger)
	svc = api.NewMetricsMiddleware(
		svc,
		kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
			Namespace: "products",
			Subsystem: "api",
			Name:      "request_count",
			Help:      "Number of requests received.",
		}, []string{"method"}),
		kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
			Namespace: "products",
			Subsystem: "api",
			Name:      "request_latency_microseconds",
			Help:      "Total duration of requests in microseconds.",
		}, []string{"method"}),
	)

	errs := make(chan error)
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errs <- fmt.Errorf("%s", <-c)
	}()

	go func() {
		logger.Info(fmt.Sprintf("Product service started on %s", productServiceURL))
		errs <- http.ListenAndServe(productServiceURL, api.MakeHTTPHandler(svc, logger))
	}()

	logger.Info(fmt.Sprintf("Users service exits, %v", <-errs))
}

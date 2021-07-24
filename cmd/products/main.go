package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	usersMongo "github.com/buraksekili/store-service/users/persistence/mongo"

	"github.com/buraksekili/store-service/users"

	"github.com/buraksekili/store-service/config/persistence"
	"github.com/buraksekili/store-service/pkg/logger"
	"github.com/buraksekili/store-service/products"
	"github.com/buraksekili/store-service/products/api"
	"github.com/buraksekili/store-service/products/persistence/mongo"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
)

func main() {
	logger := logger.New()

	productsRepo := initPersistenceLayer(logger)
	usersRepo := initUsersRepo(logger)
	productServiceURL := getPort(logger)

	var svc products.ProductService
	svc = products.New(productsRepo, usersRepo)
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

	logger.Error(fmt.Sprintf("Product service exits, %v", <-errs))
}

func initPersistenceLayer(logger logger.Logger) products.ProductRepository {
	cp := persistence.NewMongoConfigParser()
	if err := cp.Parse(); err != nil {
		logger.Error(fmt.Sprintf("cannot extract MongoDB Config, err: %v", err))
		os.Exit(1)
	}

	addr, err := cp.Address()
	if err != nil {
		logger.Error(fmt.Sprintf("cannot construct address for MongoDB, err: %v", err))
		os.Exit(1)
	}

	productsRepo, err := mongo.NewMongoDBLayer(addr)
	if err != nil {
		logger.Error(fmt.Sprintf("cannot dial %s for the DB, err: %v", addr, err))
		os.Exit(1)
	}
	return productsRepo
}

func initUsersRepo(logger logger.Logger) users.UserRepository {
	cp := persistence.NewMongoConfigParser()
	if err := cp.Parse(); err != nil {
		logger.Error(fmt.Sprintf("cannot extract MongoDB Config, err: %v", err))
		os.Exit(1)
	}

	addr, err := cp.Address()
	if err != nil {
		logger.Error(fmt.Sprintf("cannot construct address for MongoDB, err: %v", err))
		os.Exit(1)
	}

	usersRepo, err := usersMongo.NewMongoDBLayer(addr)
	if err != nil {
		logger.Error(fmt.Sprintf("cannot dial %s for the DB, err: %v", addr, err))
		os.Exit(1)
	}
	return usersRepo
}

func getPort(logger logger.Logger) (v string) {
	if v = os.Getenv("S_PRODUCTS_PORT"); v == "" {
		logger.Error("cannot get S_PRODUCTS_PORT environment variable")
		os.Exit(1)
	}
	return fmt.Sprintf(":%s", v)
}

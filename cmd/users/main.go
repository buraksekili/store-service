package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	logpkg "github.com/buraksekili/store-service/pkg/logger"

	"github.com/buraksekili/store-service/users"
	"github.com/buraksekili/store-service/users/persistence/mongo"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	stdprometheus "github.com/prometheus/client_golang/prometheus"

	"github.com/buraksekili/store-service/users/api"
)

var (
	amqpAddress    = "amqp://guest:guest@localhost:5672"
	exchangeName   = "tests"
	userServiceURL = ":8282"
)

func main() {
	logger := logpkg.New()

	u := fmt.Sprintf("mongodb://mongo:27017")
	usersRepo, err := mongo.NewMongoDBLayer(u)
	if err != nil {
		fmt.Println("CANNOT CONNECT TO", u)
	}
	var svc users.UserService
	svc = users.New(usersRepo)
	svc = api.LoggingMiddleware(svc, logger)
	svc = api.NewMetricsMiddleware(
		svc,
		kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
			Namespace: "users",
			Subsystem: "api",
			Name:      "request_count",
			Help:      "Number of requests received.",
		}, []string{"method"}),
		kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
			Namespace: "users",
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
		logger.Info(fmt.Sprintf("Users service started on %s", userServiceURL))
		errs <- http.ListenAndServe(userServiceURL, api.MakeHTTPHandler(svc, logger))
	}()

	logger.Info(fmt.Sprintf("Users service exits, %v", <-errs))
}

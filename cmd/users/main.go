package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/buraksekili/store-service/pkg/hasher"

	amqppkg "github.com/buraksekili/store-service/amqp"
	amqputil "github.com/buraksekili/store-service/config/amqp"
	"github.com/buraksekili/store-service/config/persistence"
	"github.com/buraksekili/store-service/pkg/logger"
	"github.com/buraksekili/store-service/users"
	"github.com/buraksekili/store-service/users/api"
	"github.com/buraksekili/store-service/users/persistence/mongo"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"github.com/streadway/amqp"
)

func main() {
	logger := logger.New()

	usersRepo := initPersistenceLayer(logger)
	publisher := getPublisher(logger)
	userServiceURL := getPort(logger)

	h := hasher.New(14)

	var svc users.UserService
	svc = users.New(usersRepo, *publisher, h)
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

	logger.Error(fmt.Sprintf("Users service exits, %v", <-errs))
}

func initPersistenceLayer(logger logger.Logger) users.UserRepository {
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

	usersRepo, err := mongo.NewMongoDBLayer(addr)
	if err != nil {
		logger.Error(fmt.Sprintf("cannot dial %s for the DB, err: %v", addr, err))
		os.Exit(1)
	}
	return usersRepo
}

func getPublisher(log logger.Logger) *amqppkg.AMQPPublisher {
	ac := amqputil.ExtractAMQPConfigs()
	conn, err := amqp.Dial(ac.Addr)
	if err != nil {
		log.Error(fmt.Sprintf("cannot connect to AMQP addr: %s, err: %s", ac.Addr, err.Error()))
		os.Exit(1)
	}

	publisher, err := amqppkg.NewAMQPPublisher(conn, ac.Exchange)
	if err != nil {
		log.Error(fmt.Sprintf("cannot get new AMQP publisher, err: %s", err.Error()))
		os.Exit(1)
	}
	return publisher
}

func getPort(logger logger.Logger) (v string) {
	if v = os.Getenv("S_USERS_PORT"); v == "" {
		logger.Error("cannot get S_USERS_PORT environment variable")
		os.Exit(1)
	}
	return fmt.Sprintf(":%s", v)
}

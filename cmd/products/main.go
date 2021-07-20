package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-kit/kit/log"

	"github.com/buraksekili/store-service/products"
	"github.com/buraksekili/store-service/products/api"
	"github.com/buraksekili/store-service/products/persistence/mongo"
)

const productServiceURL = ":8181"

func main() {

	logger := log.NewLogfmtLogger(os.Stdout)
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)
	logger = log.With(logger, "caller", log.DefaultCaller)

	u := fmt.Sprintf("mongodb://mongo:27017")
	productsRepo, err := mongo.NewMongoDBLayer(u)
	if err != nil {
		fmt.Println("CANNOT CONNECT TO", u)
	}
	svc := products.New(productsRepo)
	h := api.MakeHTTPHandler(svc, log.With(logger, "component", "HTTP"))

	errs := make(chan error)
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errs <- fmt.Errorf("%s", <-c)
	}()

	go func() {
		logger.Log("transport", "HTTP", "addr", productServiceURL)
		errs <- http.ListenAndServe(productServiceURL, h)
	}()

	logger.Log("exit", <-errs)
}

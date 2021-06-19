package main

import (
	"fmt"
	"log"
	"os"

	"github.com/buraksekili/store-service/config"

	"github.com/buraksekili/store-service/src/productservice/rest"

	"github.com/buraksekili/store-service/db/mongo"
)

var productServiceURL = ":8181"

func main() {

	if v := os.Getenv("PRODUCT_SERVICE_URL"); v != "" {
		productServiceURL = v
	}

	c := config.ReadDBConfig()
	MONGOURL := fmt.Sprintf("mongodb://%s:27017", c.DBHost)
	if v := os.Getenv("MONGODB_URL"); v != "" {
		MONGOURL = v
	}

	log.Println("Connecting to database", MONGOURL)
	h, err := mongo.NewMongoDBLayer(MONGOURL)
	if err != nil {
		log.Fatalf("cannot connect to the database %s, err: %v", MONGOURL, err)
	}
	log.Printf("Connected to %s!\n", MONGOURL)

	err = rest.ServerREST(productServiceURL, h)
	if err != nil {
		panic(err)
	}
}

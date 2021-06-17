package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/buraksekili/store-service/conf"
	"github.com/buraksekili/store-service/src/productservice/rest"

	"github.com/buraksekili/store-service/db/mongo"
)

var configFile = flag.String("f", "./config.json", "path for config file")

func main() {

	flag.Parse()

	// extract configuration
	config, _ := conf.ReadConfig(*configFile)

	log.Println("Connecting to database")
	h, err := mongo.NewMongoDBLayer(fmt.Sprintf("mongodb://%s:%s@%s", config.DBUser, config.DBPass, config.DBConnection))
	if err != nil {
		panic(err)
	}
	log.Printf("Connected to %s!\n", config.DBType)

	err = rest.ServerREST(config.AddrREST, h)
	if err != nil {
		panic(err)
	}
}

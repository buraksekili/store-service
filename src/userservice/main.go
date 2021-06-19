package main

import (
	"fmt"
	"log"
	"os"

	amqp2 "github.com/buraksekili/store-service/amqp"

	"github.com/streadway/amqp"

	"github.com/buraksekili/store-service/config"
	"github.com/buraksekili/store-service/db/mongo"
	"github.com/buraksekili/store-service/src/userservice/rest"
)

var (
	amqpAddress    = "amqp://guest:guest@localhost:5672"
	exchangeName   = "tests"
	userServiceURL = ":8282"
)

func main() {

	if u := os.Getenv("AMQP_URL"); u != "" {
		amqpAddress = u
	}
	if en := os.Getenv("AMQP_EXCHANGE_NAME"); en != "" {
		exchangeName = en
	}
	if v := os.Getenv("USER_SERVICE_URL"); v != "" {
		userServiceURL = v
	}

	// Connect to AMQP for RabbitMQ, default address for RabbitMQ is 'amqp://guest:guest@localhost:5672'.
	conn, err := amqp.Dial(amqpAddress)
	if err != nil {
		log.Fatalf("cannot connect to AMQP addr: %s, err: %s", amqpAddress, err.Error())
	}

	publisher, err := amqp2.NewAMQPPublisher(conn, exchangeName)
	if err != nil {
		log.Fatalf("cannot get new AMQP publisher, err: %s", err.Error())
	}

	c := config.ReadDBConfig()
	// MONGOURL := fmt.Sprintf("%s://%s:%s@%s:%s/%s", c.DBType, c.DBUser, c.DBPass, c.DBHost, c.DBPort, c.DBName)
	MONGOURL := fmt.Sprintf("mongodb://%s:27017", c.DBHost)
	if v := os.Getenv("MONGODB_URL"); v != "" {
		MONGOURL = v
	}

	log.Println("Connecting to database", MONGOURL)
	h, err := mongo.NewMongoDBLayer(MONGOURL)
	if err != nil {
		panic(err)
	}
	log.Printf("Connected to %s!\n", c.DBType)

	err = rest.ServerREST(userServiceURL, h, *publisher)
	if err != nil {
		panic(err)
	}
}

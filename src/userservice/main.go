package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	amqp2 "github.com/buraksekili/store-service/amqp"

	"github.com/streadway/amqp"

	"github.com/buraksekili/store-service/conf"
	"github.com/buraksekili/store-service/db/mongo"
	"github.com/buraksekili/store-service/src/userservice/rest"
)

const (
	ADD_USER_EVENT = "add_user"
)

var (
	configFile   = flag.String("f", "./config.json", "path for config file")
	amqpAddress  = "amqp://guest:guest@localhost:5672"
	exchangeName = "tests"
)

func main() {
	flag.Parse()

	if u := os.Getenv("AMQP_URL"); u != "" {
		amqpAddress = u
	}
	if en := os.Getenv("AMQP_EXCHANGE_NAME"); en != "" {
		exchangeName = en
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

	listener, err := amqp2.NewAMQPConsumer(conn, "listener_queue", exchangeName)
	if err != nil {
		log.Fatalf("cannot get new AMQP publisher, err: %s", err.Error())
	}

	config, _ := conf.ReadConfig(*configFile)
	log.Println("Connecting to database")
	h, err := mongo.NewMongoDBLayer(fmt.Sprintf("mongodb://%s:%s@%s", config.DBUser, config.DBPass, config.DBConnection))
	if err != nil {
		panic(err)
	}
	log.Printf("Connected to %s!\n", config.DBType)

	go func() {
		if err := listener.Listen([]string{ADD_USER_EVENT}); err != nil {
			fmt.Printf("cannot listen, err: %v", err)
		}
	}()

	err = rest.ServerREST(":8282", h, *publisher)
	if err != nil {
		panic(err)
	}
}

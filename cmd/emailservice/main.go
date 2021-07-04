package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/buraksekili/store-service/src/emailservice/emailer"
	"github.com/buraksekili/store-service/config/email"
	"github.com/buraksekili/store-service/db"
	amqp2 "github.com/buraksekili/store-service/amqp"
	"github.com/streadway/amqp"
)

const (
	ADD_USER_EVENT = "add_user"
)

var (
	amqpAddress  = "amqp://guest:guest@localhost:5672"
	exchangeName = "tests"
	queueName    = "listener_queue"
)

func main() {
	flag.Parse()

	sc := email.ExtractSMTPConfig()
	e := emailer.NewEmailer(*sc)

	if u := os.Getenv("AMQP_URL"); u != "" {
		amqpAddress = u
	}
	if en := os.Getenv("AMQP_EXCHANGE_NAME"); en != "" {
		exchangeName = en
	}
	if qn := os.Getenv("AMQP_QUEUE_NAME"); qn != "" {
		queueName = qn
	}

	// Connect to AMQP for RabbitMQ, default address for RabbitMQ is 'amqp://guest:guest@localhost:5672'.
	conn, err := amqp.Dial(amqpAddress)
	if err != nil {
		log.Fatalf("cannot connect to AMQP addr: %s, err: %s", amqpAddress, err.Error())
	}

	listener, err := amqp2.NewAMQPConsumer(conn, queueName, exchangeName)
	if err != nil {
		log.Fatalf("cannot get new AMQP publisher, err: %s", err.Error())
	}

	msgChan, err := listener.Listen([]string{ADD_USER_EVENT})
	if err != nil {
		log.Fatalf("cannot listen, err: %v", err)
	}

	log.Println("Emailservice listening...")

	for {
		for msg := range msgChan {
			u := &db.User{}
			fmt.Println("received a message", string(msg))
			if err = json.Unmarshal(msg, u); err != nil {
				fmt.Println("\ncannot unmarshal the message", string(msg))
				continue
			}

			e.Config.To = []string{u.Email}
			fmt.Println([]string{u.Email})
			if err = e.SendEmail(); err != nil {
				fmt.Println("cannot send email, err: ", err)
				continue
			}
			log.Println("Successfully send an email")
		}
	}
}

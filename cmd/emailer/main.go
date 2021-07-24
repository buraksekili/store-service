package main

import (
	"encoding/json"
	"fmt"
	"os"

	rabbitmq "github.com/buraksekili/store-service/amqp"
	amqputil "github.com/buraksekili/store-service/config/amqp"
	"github.com/buraksekili/store-service/emailer"
	"github.com/buraksekili/store-service/pkg/logger"
	"github.com/buraksekili/store-service/users"
	"github.com/streadway/amqp"
)

const (
	addUserEvent = "add_user"
)

func main() {
	log := logger.New()

	e := NewEmailService(log)
	listener := initAMQP(log)

	msgChan, err := listener.Listen([]string{addUserEvent}, log)
	if err != nil {
		log.Error(fmt.Sprintf("cannot listen messages via AMQP, err: %v", err))
		os.Exit(1)
	}

	log.Info("Email service waiting for messages.")
	for {
		for msg := range msgChan {
			u := &users.User{}
			if err = json.Unmarshal(msg, u); err != nil {
				continue
			}

			content := emailer.Email{
				To:      u.Email,
				Subject: fmt.Sprintf("Welcome to the Store %s", u.Username),
				Message: fmt.Sprintf("Hello %s!\n"+
					"Your account successfully created.\n", u.Username),
			}
			if err = e.SendEmail(content); err != nil {
				log.Error(fmt.Sprintf("cannot send an email to %v, err: %v", u.Email, err))
				continue
			}
			log.Info(fmt.Sprintf("Successfully send an email to %v, err: %v", u.Email, err))
		}
	}
}

func initAMQP(log logger.Logger) *rabbitmq.AMQPConsumer {
	ac := amqputil.ExtractAMQPConfigs()
	conn, err := amqp.Dial(ac.Addr)
	if err != nil {
		log.Error(fmt.Sprintf("cannot connect to AMQP addr: %s, err: %s", ac.Addr, err.Error()))
		os.Exit(1)
	}
	listener, err := rabbitmq.NewAMQPConsumer(conn, ac.Queue, ac.Exchange)
	if err != nil {
		log.Error(fmt.Sprintf("cannot get new AMQP Publisher, err: %v", err))
		os.Exit(1)
	}
	return listener
}

func NewEmailService(logger logger.Logger) *emailer.Emailer {
	sc, err := emailer.ExtractSMTPConfig()
	if err != nil {
		logger.Error(fmt.Sprintf("cannot extract SMTP config, err: %v", err))
		os.Exit(1)
	}
	return emailer.New(sc)
}

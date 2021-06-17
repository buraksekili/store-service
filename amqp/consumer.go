package amqp

import (
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

type AMQPConsumer struct {
	conn      *amqp.Connection
	queueName string
	exchange  string
}

func NewAMQPConsumer(conn *amqp.Connection, queueName, exchange string) (*AMQPConsumer, error) {
	consumer := &AMQPConsumer{conn, queueName, exchange}
	c, err := consumer.conn.Channel()
	if err != nil {
		return nil, err
	}
	defer c.Close()

	if err = c.ExchangeDeclare(consumer.exchange, "topic", true, false, false, false, nil); err != nil {
		return nil, err
	}
	if _, err = c.QueueDeclare(consumer.queueName, true, false, false, false, nil); err != nil {
		return nil, err
	}

	return consumer, nil
}

func (c *AMQPConsumer) Listen(topics []string) error {
	ch, err := c.conn.Channel()
	if err != nil {
		return err
	}

	for _, t := range topics {
		if err := ch.QueueBind(c.queueName, t, c.exchange, false, nil); err != nil {
			return err
		}
	}

	msgs, err := ch.Consume(c.queueName, "", false, false, false, false, nil)
	if err != nil {
		return err
	}
	forever := make(chan bool)
	go func() {
		log.Println("\t===\twaiting to receive a message\t===")
		for msg := range msgs {
			log.Printf("Received a message: %s\n", msg.Body)
			if err = msg.Ack(false); err != nil {
				fmt.Printf("cannot acknowledge msg=%#v, err: %v", msg, err)
			}
		}
	}()

	<-forever
	return nil
}

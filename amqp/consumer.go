package amqp

import (
	"fmt"

	"github.com/buraksekili/store-service/pkg/logger"
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

func (c *AMQPConsumer) Listen(topics []string, logger logger.Logger) (<-chan []byte, error) {
	ch, err := c.conn.Channel()
	if err != nil {
		return nil, err
	}

	for _, t := range topics {
		if err := ch.QueueBind(c.queueName, t, c.exchange, false, nil); err != nil {
			return nil, err
		}
	}

	msgs, err := ch.Consume(c.queueName, "", false, false, false, false, nil)
	if err != nil {
		return nil, err
	}

	forever := make(chan []byte)
	go func() {
		logger.Info("RabbitMQ: waiting to receive a message.")
		for msg := range msgs {
			forever <- msg.Body
			if err = msg.Ack(false); err != nil {
				logger.Error(fmt.Sprintf("cannot acknowledge msg:%#v, err: %v", msg, err))
			}
		}
	}()

	return forever, nil
}

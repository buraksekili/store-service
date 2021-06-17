package amqp

import (
	"encoding/json"

	"github.com/streadway/amqp"
)

type AMQPPublisher struct {
	conn     *amqp.Connection
	exchange string
}

func NewAMQPPublisher(conn *amqp.Connection, exchange string) (*AMQPPublisher, error) {
	publisher := &AMQPPublisher{conn, exchange}

	c, err := publisher.conn.Channel()
	if err != nil {
		return nil, err
	}
	defer c.Close()
	if err = c.ExchangeDeclare(exchange, "topic", true, false, false, false, nil); err != nil {
		return nil, err
	}

	return publisher, nil
}

func (p *AMQPPublisher) Publish(e Event) error {
	c, err := p.conn.Channel()
	if err != nil {
		return err
	}
	defer c.Close()

	jb, err := json.Marshal(e)
	if err != nil {
		return err
	}

	return c.Publish(
		p.exchange,
		e.Name(),
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        jb,
		},
	)
}

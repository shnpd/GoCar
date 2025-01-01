package amqpclt

import (
	"context"
	carpb "coolcar/car/api/gen/v1"
	"encoding/json"
	"fmt"

	"github.com/streadway/amqp"
)

type Publisher struct {
	ch       *amqp.Channel
	exchange string
}

// NewPublisher creates a new publisher.
func NewPublisher(conn *amqp.Connection, exchange string) (*Publisher, error) {
	// channel is a virtual connection inside a connection
	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("cannot allocate channel: %v", err)
	}

	err = ch.ExchangeDeclare(
		exchange, // name
		"fanout", // kind
		true,     // durable
		false,    // autoDelete
		false,    // internal
		false,    // noWait
		nil,      // args
	)
	if err != nil {
		return nil, fmt.Errorf("cannot declare exchange: %v", err)
	}
	return &Publisher{
		ch:       ch,
		exchange: exchange,
	}, nil
}

// Publish publishes a message.
func (p *Publisher) Publish(c context.Context, car *carpb.CarEntity) error {
	b, err := json.Marshal(car)
	if err != nil {
		return fmt.Errorf("cannot marshal car: %v", err)
	}

	return p.ch.Publish(
		p.exchange,
		"",    // routing key
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			Body: b,
		}, // message
	)
}

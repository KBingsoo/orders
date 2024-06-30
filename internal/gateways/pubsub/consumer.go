package pubsub

import (
	card "github.com/KBingsoo/cards/pkg/models/event"
	"github.com/KBingsoo/orders/pkg/models/event"
	"github.com/seosoojin/go-rabbit/rabbit"
	"github.com/streadway/amqp"
)

type cardConsumer struct {
	internalConsumer rabbit.Consumer[card.Event]
}

func NewCardConsumer(conn *amqp.Connection) (*cardConsumer, error) {
	internalConsumer, err := rabbit.NewConsumer(conn, decode[card.Event])
	if err != nil {
		return nil, err
	}

	return &cardConsumer{
		internalConsumer: internalConsumer,
	}, nil
}

func (c *cardConsumer) Consume(fn func(card.Event) error) error {
	return c.internalConsumer.Consume(fn)
}

type consumer struct {
	internalConsumer rabbit.Consumer[event.Event]
}

func NewConsumer(conn *amqp.Connection) (*consumer, error) {
	internalConsumer, err := rabbit.NewConsumer(conn, decode[event.Event])
	if err != nil {
		return nil, err
	}

	return &consumer{
		internalConsumer: internalConsumer,
	}, nil
}

func (c *consumer) Consume(fn func(event.Event) error) error {
	return c.internalConsumer.Consume(fn)
}

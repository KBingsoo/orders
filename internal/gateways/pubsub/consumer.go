package pubsub

import (
	"github.com/KBingsoo/cards/pkg/models/event"
	"github.com/seosoojin/go-rabbit/rabbit"
	"github.com/streadway/amqp"
)

type consumer struct {
	internalConsumer rabbit.Consumer[event.Event]
}

func NewCardConsumer(conn *amqp.Connection) (*consumer, error) {
	internalConsumer, err := rabbit.NewConsumer(conn, decode)
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

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

func (c *consumer) Consume() error {
	return c.internalConsumer.Consume(c.handler)
}

func (c *consumer) handler(entry event.Event) error {
	switch entry.Type {
	case event.Update:
	default:
		return nil
	}

	return nil
}

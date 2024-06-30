package pubsub

import (
	"encoding/json"
	"fmt"

	card "github.com/KBingsoo/cards/pkg/models/event"
	"github.com/KBingsoo/orders/pkg/models/event"
	"github.com/seosoojin/go-rabbit/rabbit"
	"github.com/seosoojin/go-rabbit/rabbit/message"
	"github.com/streadway/amqp"
)

type cardProducer struct {
	internalProducer rabbit.Producer
}

func NewCardProducer(conn *amqp.Connection) (*cardProducer, error) {
	internalProducer, err := rabbit.NewProducer(conn)
	if err != nil {
		return nil, err
	}

	return &cardProducer{
		internalProducer: internalProducer,
	}, nil
}

func (p *cardProducer) Emit(event card.Event) error {

	b, err := json.Marshal(event)
	if err != nil {
		return err
	}

	key := fmt.Sprintf("card.%s.%s.%s", event.Card.ID, event.Type, event.Time.Format("2006-01-02"))

	msg := message.Message{
		Key: key,
		Value: amqp.Publishing{
			ContentType: "application/json",
			Body:        b,
		},
	}

	return p.internalProducer.Emit(msg)
}

type producer struct {
	internalProducer rabbit.Producer
}

func NewProducer(conn *amqp.Connection) (*producer, error) {
	internalProducer, err := rabbit.NewProducer(conn)
	if err != nil {
		return nil, err
	}

	return &producer{
		internalProducer: internalProducer,
	}, nil
}

func (p *producer) Emit(event event.Event) error {
	b, err := json.Marshal(event)
	if err != nil {
		return err
	}

	key := fmt.Sprintf("order.%s.%s", event.Order.ID, event.Type)

	msg := message.Message{
		Key: key,
		Value: amqp.Publishing{
			ContentType: "application/json",
			Body:        b,
		},
	}

	return p.internalProducer.Emit(msg)
}

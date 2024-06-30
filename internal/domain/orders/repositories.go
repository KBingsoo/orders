package orders

import (
	card "github.com/KBingsoo/cards/pkg/models/event"
	"github.com/KBingsoo/entities/pkg/models"
	"github.com/KBingsoo/orders/pkg/models/event"
	"github.com/literalog/go-wise/wise"
)

type Repository wise.MongoRepository[models.Order]
type UserRepo wise.MongoRepository[models.User]

type CardProducer interface {
	Emit(event card.Event) error
}

type CardConsumer interface {
	Consume(fn func(card.Event) error) error
}

type Producer interface {
	Emit(event event.Event) error
}

type Consumer interface {
	Consume(fn func(event.Event) error) error
}

package orders

import (
	"github.com/KBingsoo/cards/pkg/models/event"
	"github.com/KBingsoo/entities/pkg/models"
	"github.com/literalog/go-wise/wise"
)

type Repository wise.MongoRepository[models.Order]

type CardProducer interface {
	Emit(event event.Event) error
}

type CardConsumer interface {
	Consume(fn func(event.Event) error) error
}

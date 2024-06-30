package pubsub

import (
	"encoding/json"

	card "github.com/KBingsoo/cards/pkg/models/event"
	"github.com/KBingsoo/orders/pkg/models/event"
)

type dataTypes interface {
	card.Event | event.Event
}

func decode[T dataTypes](data []byte) (T, error) {
	e := new(T)
	if err := json.Unmarshal(data, e); err != nil {
		return *e, err
	}

	return *e, nil
}

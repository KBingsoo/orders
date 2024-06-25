package pubsub

import (
	"encoding/json"

	"github.com/KBingsoo/cards/pkg/models/event"
)

func decode(data []byte) (event.Event, error) {
	e := new(event.Event)
	if err := json.Unmarshal(data, e); err != nil {
		return event.Event{}, err
	}

	return *e, nil
}

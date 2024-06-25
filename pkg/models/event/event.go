package event

import (
	"context"
	"time"

	"github.com/KBingsoo/entities/pkg/models"
)

type EventType string

const (
	Create EventType = "order_create"
	Update EventType = "order_update"
	Delete EventType = "order_delete"

	Succeed EventType = "order_succeed"
	Error   EventType = "order_error"
)

type Event struct {
	Type    EventType       `json:"type"`
	Time    time.Time       `json:"time"`
	Order   models.Order    `json:"order"`
	Context context.Context `json:"context"`
}

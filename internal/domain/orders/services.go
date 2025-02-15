package orders

import (
	"context"
	"time"

	cards "github.com/KBingsoo/cards/pkg/models/event"
	"github.com/KBingsoo/entities/pkg/models"
	"github.com/google/uuid"
)

type Manager interface {
	GetByID(ctx context.Context, orderID string) (models.Order, error)
	Create(ctx context.Context, order *models.Order) error
	ConfirmShipment(ctx context.Context, orderID string) error
	Fulfill(ctx context.Context, orderID string) error
}

type itemsMap struct {
	internalMal map[string]bool
	order       models.Order
	succeed     int
}

type manager struct {
	repository   Repository
	cardProducer CardProducer
	cardConsumer CardConsumer
	orderMap     map[string]itemsMap
}

func NewManager(repository Repository, cardProducer CardProducer, cardConsumer CardConsumer) *manager {
	return &manager{
		repository:   repository,
		cardProducer: cardProducer,
		cardConsumer: cardConsumer,
		orderMap:     make(map[string]itemsMap),
	}
}

func (m *manager) GetByID(ctx context.Context, orderID string) (models.Order, error) {
	return m.repository.Find(ctx, orderID)
}

func (m *manager) Create(ctx context.Context, order *models.Order) error {
	if order.ID == "" {
		order.ID = uuid.NewString()
	}

	return m.repository.Upsert(ctx, order.ID, *order)
}

func (m *manager) ConfirmShipment(ctx context.Context, orderID string) error {
	order, err := m.repository.Find(ctx, orderID)
	if err != nil {
		return err
	}

	order.Status = models.COMPLETED
	return m.repository.Upsert(ctx, orderID, order)
}

func (m *manager) Fulfill(ctx context.Context, orderID string) error {
	order, err := m.repository.Find(ctx, orderID)
	if err != nil {
		return err
	}

	m.orderMap[orderID] = itemsMap{
		internalMal: make(map[string]bool),
		order:       order,
	}

	for _, card := range order.Cards {
		event := cards.Event{
			Type: cards.OrderFulfill,
			Time: time.Now(),
			Card: models.Card{
				ID: card,
			},
			OrderID: orderID,
			Context: ctx,
		}

		err = m.cardProducer.Emit(event)
		if err != nil {
			return err
		}

		m.orderMap[orderID].internalMal[card] = false
	}

	return nil

}

func (m *manager) revert(ctx context.Context, orderID string) error {
	return m.cardProducer.Emit(cards.Event{
		Type:    cards.OrderRevert,
		Time:    time.Now(),
		OrderID: orderID,
		Context: ctx,
	})
}

func (m *manager) cardHandler(entry cards.Event) error {
	switch entry.Type {
	case cards.Succeed:
		order, ok := m.orderMap[entry.OrderID]
		if !ok {
			m.revert(entry.Context, entry.OrderID)
			return nil
		}
		order.succeed++
		order.internalMal[entry.Card.ID] = true

		if order.succeed == len(order.order.Cards) {
			order.order.Status = models.SHIPPING
			err := m.repository.Upsert(entry.Context, entry.OrderID, order.order)
			if err != nil {
				m.revert(entry.Context, entry.OrderID)
				return err
			}
			delete(m.orderMap, entry.OrderID)
		}

	case cards.Error:
		err := m.revert(entry.Context, entry.OrderID)
		if err != nil {
			return err
		}
		delete(m.orderMap, entry.OrderID)

	}

	return nil
}

func (m *manager) Consume() error {
	return m.cardConsumer.Consume(m.cardHandler)
}

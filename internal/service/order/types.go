package order

import (
	responseOrder "backend/internal/controllers/http/response/order"
	"backend/internal/entity/order"
	"context"
	"sync"

	"github.com/google/uuid"
)

type (
	Repo interface {
		GetOrderByID(ctx context.Context, orderID uuid.UUID) (order.Order, error)
		GetOrders(ctx context.Context, limit, offset int, filters map[string]any) ([]order.Order, error)
		CreateOrder(ctx context.Context, body order.Order) error
		UpdateOrder(ctx context.Context, body order.Order) error

		GetEventByID(ctx context.Context, eventID uuid.UUID) (order.OrderEvent, error)
		GetEventsByOrder(ctx context.Context, orderID uuid.UUID) ([]order.OrderEvent, error)
		EventExists(ctx context.Context, eventID uuid.UUID) (bool, error)
		CreateEvent(ctx context.Context, body order.OrderEvent) error
	}

	Service struct {
		repo         Repo
		clients      map[uuid.UUID][]chan responseOrder.OrderEvent
		clientMux    sync.Mutex
		eventBuffers map[uuid.UUID]*EventBuffer
	}

	EventBuffer struct {
		mu     sync.Mutex
		events []responseOrder.OrderEvent
	}
)

func NewService(repo Repo) *Service {
	return &Service{
		repo:         repo,
		clients:      make(map[uuid.UUID][]chan responseOrder.OrderEvent),
		eventBuffers: make(map[uuid.UUID]*EventBuffer),
	}
}

func NewEventBuffer() *EventBuffer {
	return &EventBuffer{
		events: make([]responseOrder.OrderEvent, 0),
	}
}

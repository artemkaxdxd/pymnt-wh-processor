package order

import (
	"backend/config"
	"backend/internal/entity/order"
	"time"

	"github.com/google/uuid"
)

type (
	Order struct {
		OrderID   uuid.UUID          `json:"order_id"`
		UserID    uuid.UUID          `json:"user_id"`
		Status    config.OrderStatus `json:"status"`
		IsFinal   bool               `json:"is_final"`
		CreatedAt time.Time          `json:"created_at"`
		UpdatedAt time.Time          `json:"updated_at"`
	}

	OrderEvent struct {
		EventID     uuid.UUID          `json:"-"`
		OrderID     uuid.UUID          `json:"order_id"`
		UserID      uuid.UUID          `json:"user_id"`
		OrderStatus config.OrderStatus `json:"order_status"`
		IsFinal     bool               `json:"is_final"`
		CreatedAt   time.Time          `json:"created_at"`
		UpdatedAt   time.Time          `json:"updated_at"`
	}
)

func OrderToResponse(o order.Order) Order {
	return Order{
		OrderID:   o.ID,
		UserID:    o.UserID,
		Status:    o.Status,
		IsFinal:   o.IsFinal,
		CreatedAt: o.CreatedAt,
		UpdatedAt: o.UpdatedAt,
	}
}

func OrdersToResponse(orders []order.Order) []Order {
	resp := make([]Order, len(orders))
	for i, v := range orders {
		resp[i] = OrderToResponse(v)
	}
	return resp
}

func OrderToEventResponse(o order.Order) OrderEvent {
	return OrderEvent{
		OrderID:     o.ID,
		UserID:      o.UserID,
		OrderStatus: o.Status,
		IsFinal:     o.IsFinal,
		CreatedAt:   o.CreatedAt,
		UpdatedAt:   o.UpdatedAt,
	}
}

func EventToResponse(o order.OrderEvent) OrderEvent {
	return OrderEvent{
		EventID:     o.ID,
		OrderID:     o.OrderID,
		UserID:      o.UserID,
		OrderStatus: o.OrderStatus,
		IsFinal:     config.IsFinalOrderStatus(o.OrderStatus),
		CreatedAt:   o.CreatedAt,
		UpdatedAt:   o.UpdatedAt,
	}
}

func EventsToResponse(events []order.OrderEvent) []OrderEvent {
	resp := make([]OrderEvent, len(events))
	for i, v := range events {
		resp[i] = EventToResponse(v)
	}
	return resp
}

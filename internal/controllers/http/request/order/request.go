package order

import (
	"backend/config"
	"backend/internal/entity/order"
	"time"

	"github.com/google/uuid"
)

type OrderEvent struct {
	EventID     uuid.UUID          `json:"event_id"`
	OrderID     uuid.UUID          `json:"order_id"`
	UserID      uuid.UUID          `json:"user_id"`
	OrderStatus config.OrderStatus `json:"order_status"`
	UpdatedAt   time.Time          `json:"updated_at"`
	CreatedAt   time.Time          `json:"created_at"`
}

func (o OrderEvent) ToEntity() order.OrderEvent {
	return order.OrderEvent{
		ID:          o.EventID,
		OrderID:     o.OrderID,
		UserID:      o.UserID,
		OrderStatus: o.OrderStatus,
		UpdatedAt:   o.UpdatedAt,
		CreatedAt:   o.CreatedAt,
	}
}

package order

import (
	"backend/config"
	"time"

	"github.com/google/uuid"
)

type (
	Order struct {
		ID        uuid.UUID
		UserID    uuid.UUID
		Status    config.OrderStatus
		IsFinal   bool
		UpdatedAt time.Time
		CreatedAt time.Time
	}

	OrderEvent struct {
		ID          uuid.UUID
		OrderID     uuid.UUID
		UserID      uuid.UUID
		OrderStatus config.OrderStatus
		UpdatedAt   time.Time
		CreatedAt   time.Time
	}
)

func (Order) TableName() string {
	return "orders"
}

func (OrderEvent) TableName() string {
	return "order_events"
}

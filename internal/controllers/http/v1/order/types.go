package order

import (
	"backend/config"
	requestOrder "backend/internal/controllers/http/request/order"
	responseOrder "backend/internal/controllers/http/response/order"
	"backend/pkg/logger"
	"context"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type (
	Service interface {
		GetOrders(ctx context.Context, limit, offset int, filters map[string]any) ([]responseOrder.Order, config.ServiceCode, error)
		ProcessOrderEvent(ctx context.Context, body requestOrder.OrderEvent) (config.ServiceCode, error)
		Subscribe(orderID uuid.UUID) <-chan responseOrder.OrderEvent
	}

	handler struct {
		l   logger.Logger
		svc Service
	}
)

func InitHandler(
	g *gin.Engine,
	l logger.Logger,
	svc Service,
) {
	handler := handler{
		l:   l,
		svc: svc,
	}

	g.POST("/webhooks/payments/orders", handler.processOrderPayment)

	orders := g.Group("orders")
	{
		orders.GET("", handler.getOrders)
		orders.GET("/:order_id/events", handler.subscribeToOrders)
	}
}

package order

import (
	"backend/internal/entity/order"
	"backend/pkg/db"
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

type Repo struct {
	db.Database
}

func NewRepo(db db.Database) Repo {
	return Repo{Database: db}
}

func (r Repo) GetOrderByID(ctx context.Context, orderID uuid.UUID) (order order.Order, err error) {
	err = r.Database.Instance().WithContext(ctx).Raw(`
		SELECT * FROM orders WHERE id = ?`, orderID.String()).Scan(&order).Error
	return
}

func (r Repo) GetOrders(ctx context.Context, limit, offset int, filters map[string]any) (orders []order.Order, err error) {
	var (
		query strings.Builder
		args  []any
	)

	query.WriteString(`
		SELECT * FROM orders
		WHERE 1 = 1`)

	if userId, ok := filters["user_id"]; ok {
		query.WriteString(" AND user_id = ?")
		args = append(args, userId)
	}

	if statuses, ok := filters["status"]; ok {
		query.WriteString(" AND status IN (?)")
		args = append(args, statuses)
	}

	if isFinal, ok := filters["is_final"]; ok {
		query.WriteString(" AND is_final = ?")
		args = append(args, isFinal)
	}

	if sortBy, ok := filters["sort_by"]; ok {
		query.WriteString(fmt.Sprintf(" ORDER BY %s", sortBy))
	}

	if sortOrder, ok := filters["sort_order"]; ok {
		query.WriteString(fmt.Sprintf(" %s", sortOrder))
	}

	query.WriteString(" LIMIT ? OFFSET ?")
	args = append(args, limit, offset)

	err = r.Database.Instance().WithContext(ctx).Raw(query.String(), args...).Scan(&orders).Error
	return
}

func (r Repo) CreateOrder(ctx context.Context, body order.Order) error {
	return r.Database.Instance().WithContext(ctx).Create(&body).Error
}

func (r Repo) UpdateOrder(ctx context.Context, body order.Order) error {
	return r.Database.Instance().WithContext(ctx).Where("id = ?", body.ID).Updates(&body).Error
}

func (r Repo) GetEventByID(ctx context.Context, eventID uuid.UUID) (event order.OrderEvent, err error) {
	err = r.Database.Instance().WithContext(ctx).Raw(`
		SELECT * FROM order_events WHERE id = ?`, eventID.String()).Scan(&event).Error
	return
}

func (r Repo) GetEventsByOrder(ctx context.Context, orderID uuid.UUID) (events []order.OrderEvent, err error) {
	err = r.Database.Instance().WithContext(ctx).Raw(`
		SELECT oe.*
		FROM order_events oe
		JOIN orders o ON oe.order_id = o.id
		WHERE oe.order_id = ?
		ORDER BY oe.updated_at ASC`, orderID).Scan(&events).Error
	return
}

func (r Repo) EventExists(ctx context.Context, eventID uuid.UUID) (exists bool, err error) {
	err = r.Database.Instance().WithContext(ctx).Raw(`
		SELECT EXISTS (
			SELECT 1 FROM order_events WHERE id = ?)`, eventID).Scan(&exists).Error
	return
}

func (r Repo) CreateEvent(ctx context.Context, body order.OrderEvent) error {
	return r.Database.Instance().WithContext(ctx).Create(&body).Error
}

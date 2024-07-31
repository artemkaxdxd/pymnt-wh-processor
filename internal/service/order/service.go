package order

import (
	"backend/config"
	requestOrder "backend/internal/controllers/http/request/order"
	responseOrder "backend/internal/controllers/http/response/order"
	entityOrder "backend/internal/entity/order"
	"context"
	"slices"
	"time"

	"github.com/google/uuid"
)

func (s *Service) GetOrders(ctx context.Context, limit, offset int, filters map[string]any) ([]responseOrder.Order, config.ServiceCode, error) {
	orders, err := s.repo.GetOrders(ctx, limit, offset, filters)
	if err != nil {
		return nil, config.CodeDatabaseError, err
	}
	return responseOrder.OrdersToResponse(orders), config.CodeOK, nil
}

func (s *Service) ProcessOrderEvent(ctx context.Context, body requestOrder.OrderEvent) (config.ServiceCode, error) {
	exists, err := s.repo.EventExists(ctx, body.EventID)
	if err != nil {
		return config.CodeDatabaseError, err
	}
	if exists {
		return config.CodeDuplicateEvent, config.ErrDuplicateEvent
	}

	order, err := s.repo.GetOrderByID(ctx, body.OrderID)
	if err != nil {
		return config.CodeDatabaseError, err
	}
	if order.IsFinal {
		return config.CodeOrderEnded, err
	}

	isFinal := config.IsFinalOrderStatus(body.OrderStatus)

	if order.ID != uuid.Nil {
		order.Status = body.OrderStatus
		order.UpdatedAt = body.UpdatedAt
		order.IsFinal = isFinal
		err = s.repo.UpdateOrder(ctx, order)
	} else {
		order = entityOrder.Order{
			ID:        body.OrderID,
			UserID:    body.UserID,
			Status:    body.OrderStatus,
			IsFinal:   isFinal,
			CreatedAt: body.CreatedAt,
			UpdatedAt: body.UpdatedAt,
		}
		err = s.repo.CreateOrder(ctx, order)
	}
	if err != nil {
		return config.CodeDatabaseError, err
	}

	event := body.ToEntity()
	if err = s.repo.CreateEvent(ctx, event); err != nil {
		return config.CodeDatabaseError, err
	}

	s.bufferEvent(body.OrderID, responseOrder.EventToResponse(event))
	s.checkAndSendEvents(body.OrderID)

	return config.CodeOK, nil
}

func (s *Service) bufferEvent(orderID uuid.UUID, event responseOrder.OrderEvent) {
	s.clientMux.Lock()
	defer s.clientMux.Unlock()

	if s.eventBuffers[orderID] == nil {
		s.eventBuffers[orderID] = NewEventBuffer()
	}
	s.eventBuffers[orderID].AddEvent(event)
}

func (s *Service) checkAndSendEvents(orderID uuid.UUID) {
	s.clientMux.Lock()
	defer s.clientMux.Unlock()

	buffer := s.eventBuffers[orderID]
	if buffer == nil {
		return
	}

	pastEvents, err := s.repo.GetEventsByOrder(context.Background(), orderID)
	if err != nil {
		return
	}

	// Combine past events and buffered events for sequence validation
	allEvents := append(responseOrder.EventsToResponse(pastEvents), buffer.GetEvents()...)

	for buffer.HasNextEvent() {
		nextEvent := buffer.GetNextEvent()

		if !s.isValidSequence(allEvents, nextEvent) {
			buffer.AddEvent(nextEvent) // Re-buffer the event if the sequence is incorrect
			return
		}

		for _, ch := range s.clients[orderID] {
			ch <- nextEvent
		}
		allEvents = append(allEvents, nextEvent)

		if nextEvent.IsFinal {
			buffer.ClearEvents()
			for _, ch := range s.clients[orderID] {
				close(ch)
			}
			delete(s.clients, orderID)
		}
	}
}

func (s *Service) Subscribe(orderID uuid.UUID) <-chan responseOrder.OrderEvent {
	s.clientMux.Lock()
	defer s.clientMux.Unlock()

	ch := make(chan responseOrder.OrderEvent, 10)
	s.clients[orderID] = append(s.clients[orderID], ch)

	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()

		var chanClosed bool
		pastEvents, _ := s.repo.GetEventsByOrder(ctx, orderID)
		for _, event := range pastEvents {
			isFinal := config.IsFinalOrderStatus(event.OrderStatus)
			ch <- responseOrder.EventToResponse(event)
			if isFinal {
				chanClosed = true
				s.closeCh(ch, orderID)
				return
			}
		}

		<-ctx.Done()
		if !chanClosed {
			s.closeCh(ch, orderID)
		}
	}()

	return ch
}

func (s *Service) closeCh(ch chan responseOrder.OrderEvent, orderID uuid.UUID) {
	close(ch)
	s.clientMux.Lock()
	defer s.clientMux.Unlock()
	for i, clientCh := range s.clients[orderID] {
		if clientCh == ch {
			s.clients[orderID] = append(s.clients[orderID][:i], s.clients[orderID][i+1:]...)
			break
		}
	}
}

func (*Service) isValidSequence(events []responseOrder.OrderEvent, newEvent responseOrder.OrderEvent) bool {
	if len(events) == 0 {
		return newEvent.OrderStatus == config.OrderStatusCreated
	}

	// Needs optimization (maybe use trie)
	lastEvent := events[len(events)-1]
	validTransitions := map[config.OrderStatus][]config.OrderStatus{
		config.OrderStatusCreated:     {config.OrderStatusPending, config.OrderStatusChangedMind, config.OrderStatusFailed},
		config.OrderStatusPending:     {config.OrderStatusConfirmed, config.OrderStatusChangedMind, config.OrderStatusFailed},
		config.OrderStatusConfirmed:   {config.OrderStatusChinazes, config.OrderStatusChangedMind, config.OrderStatusFailed},
		config.OrderStatusChinazes:    {config.OrderStatusMoneyBack},
		config.OrderStatusChangedMind: {},
		config.OrderStatusFailed:      {},
		config.OrderStatusMoneyBack:   {},
	}

	return slices.Contains(validTransitions[lastEvent.OrderStatus], newEvent.OrderStatus)
}

package order

import (
	responseOrder "backend/internal/controllers/http/response/order"
)

func (b *EventBuffer) AddEvent(event responseOrder.OrderEvent) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.events = append(b.events, event)
}

func (b *EventBuffer) HasNextEvent() bool {
	b.mu.Lock()
	defer b.mu.Unlock()
	return len(b.events) > 0
}

func (b *EventBuffer) GetNextEvent() responseOrder.OrderEvent {
	b.mu.Lock()
	defer b.mu.Unlock()
	if len(b.events) > 0 {
		event := b.events[0]
		b.events = b.events[1:]
		return event
	}
	return responseOrder.OrderEvent{}
}

func (b *EventBuffer) GetEvents() []responseOrder.OrderEvent {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.events
}

func (b *EventBuffer) ClearEvents() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.events = nil
}

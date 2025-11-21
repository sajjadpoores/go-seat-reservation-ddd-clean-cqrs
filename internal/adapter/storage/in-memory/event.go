package storage

import (
	"context"
	"fmt"
	"sync"

	"github.com/sajjadpoores/go-seat-reservation-ddd-clean-cqrs/internal/domain/event"
)

type InMemoryRepository struct {
	mu     sync.RWMutex
	events map[event.EventID]*event.Event
}

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		events: make(map[event.EventID]*event.Event),
	}
}

func (r *InMemoryRepository) FindByID(ctx context.Context, id event.EventID) (*event.Event, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	e, exists := r.events[id]

	if !exists {
		return nil, event.ErrEventNotFound
	}

	return deepCopy(e), nil
}

func (r *InMemoryRepository) Save(ctx context.Context, event *event.Event) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.events[event.Id] = deepCopy(event)
	return nil
}

func deepCopy(e *event.Event) *event.Event {
	newE := *e
	newE.Holds = make(map[event.HoldID]*event.Hold)
	for k, v := range e.Holds {
		val := *v
		newE.Holds[k] = &val
	}
	return &newE
}

func (r *InMemoryRepository) PrintAllEvents(ctx context.Context) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, e := range r.events {
		fmt.Printf("Event ID: %s, Title: %s, Capacity: %d, Sold: %d, Holds: %d\n", e.Id, e.Title, e.Capacity, e.Sold, len(e.Holds))
		for hold := range e.Holds {
			fmt.Printf("  Hold ID: %s, Quantity: %d, Customer ID: %s, Expires At: %s\n", hold, e.Holds[hold].Quantity, e.Holds[hold].CustomerID, e.Holds[hold].ExpiresAt)
		}
	}
}

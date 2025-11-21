package event

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrEventCapacityExceeded = errors.New("event capacity exceeded")
	ErrHoldNotFound          = errors.New("hold not found")
)

type EventID uuid.UUID
type HoldID uuid.UUID
type CustomerID uuid.UUID

type Hold struct {
	Id         HoldID
	quantity   int
	customerID CustomerID
	expiresAt  time.Time
}

func (h *Hold) isExpired(currentTime time.Time) bool {
	if h.expiresAt.Before(currentTime) {
		return true
	}
	return false
}

type Event struct {
	Id       EventID
	title    string
	capacity int
	sold     int
	holds    map[HoldID]*Hold
}

func (e *Event) pruneExpiredHolds(currentTime time.Time) {
	for HoldId, hold := range e.holds {
		if hold.isExpired(currentTime) {
			delete(e.holds, HoldId)
		}
	}
}

func (e *Event) calculateActiveHolds(currentTime time.Time) int {
	total := e.capacity
	e.pruneExpiredHolds(currentTime)
	for _, hold := range e.holds {
		total -= hold.quantity
	}
	return total
}

func createNewHold(e *Event, customerID CustomerID, quantity int, holdDuration time.Duration, currentTime time.Time) (*Hold, error) {
	holdID := HoldID(uuid.New())
	hold := &Hold{
		Id:         holdID,
		quantity:   quantity,
		customerID: customerID,
		expiresAt:  currentTime.Add(holdDuration),
	}
	return hold, nil
}

func (e *Event) requestHold(customerID CustomerID, quantity int, holdDuration time.Duration, currentTime time.Time) (*Hold, error) {
	e.pruneExpiredHolds(currentTime)
	availables := e.calculateActiveHolds(currentTime)
	if quantity+e.sold > availables {
		return nil, ErrEventCapacityExceeded
	}

	hold, err := createNewHold(e, customerID, quantity, holdDuration, currentTime)
	if err != nil {
		return nil, err
	}

	e.holds[hold.Id] = hold
	return hold, nil
}

func (e *Event) confirmHold(holdID HoldID, currentTime time.Time) error {
	e.pruneExpiredHolds(currentTime)
	hold, exists := e.holds[holdID]

	if !exists {
		return ErrHoldNotFound
	}

	e.sold += hold.quantity
	// TODO: EMIT AN EVENT SO THAT LATER A NEW ORDER IS CREATED
	delete(e.holds, holdID)
	return nil
}

func (e *Event) cancelHold(holdID HoldID, currentTime time.Time) error {
	_, exists := e.holds[holdID]

	if !exists {
		return ErrHoldNotFound
	}

	delete(e.holds, holdID)
	return nil
}

func CreateNewEvent(title string, capacity int) *Event {
	eventID := EventID(uuid.New())
	return &Event{
		Id:       eventID,
		title:    title,
		capacity: capacity,
		sold:     0,
		holds:    make(map[HoldID]*Hold),
	}
}

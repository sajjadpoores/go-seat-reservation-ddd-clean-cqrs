package event

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrEventCapacityExceeded = errors.New("event capacity exceeded")
	ErrHoldNotFound          = errors.New("hold not found")
	ErrEventNotFound         = errors.New("event not found")
)

type EventID uuid.UUID
type HoldID uuid.UUID
type CustomerID uuid.UUID

type Hold struct {
	Id         HoldID
	Quantity   int
	CustomerID CustomerID
	ExpiresAt  time.Time
}

func (h *Hold) isExpired(currentTime time.Time) bool {
	if h.ExpiresAt.Before(currentTime) {
		return true
	}
	return false
}

type Event struct {
	Id       EventID
	Title    string
	Capacity int
	Sold     int
	Holds    map[HoldID]*Hold
}

func (e *Event) pruneExpiredHolds(currentTime time.Time) {
	for HoldId, hold := range e.Holds {
		if hold.isExpired(currentTime) {
			delete(e.Holds, HoldId)
		}
	}
}

func (e *Event) calculateActiveHolds(currentTime time.Time) int {
	total := e.Capacity
	e.pruneExpiredHolds(currentTime)
	for _, hold := range e.Holds {
		total -= hold.Quantity
	}
	return total
}

func createNewHold(e *Event, customerID CustomerID, quantity int, holdDuration time.Duration, currentTime time.Time) (*Hold, error) {
	holdID := HoldID(uuid.New())
	hold := &Hold{
		Id:         holdID,
		Quantity:   quantity,
		CustomerID: customerID,
		ExpiresAt:  currentTime.Add(holdDuration),
	}
	return hold, nil
}

func (e *Event) RequestHold(customerID CustomerID, quantity int, holdDuration time.Duration, currentTime time.Time) (*Hold, error) {
	e.pruneExpiredHolds(currentTime)
	availables := e.calculateActiveHolds(currentTime)
	if quantity+e.Sold > availables {
		return nil, ErrEventCapacityExceeded
	}

	hold, err := createNewHold(e, customerID, quantity, holdDuration, currentTime)
	if err != nil {
		return nil, err
	}

	e.Holds[hold.Id] = hold
	return hold, nil
}

func (e *Event) ConfirmHold(holdID HoldID, currentTime time.Time) error {
	e.pruneExpiredHolds(currentTime)
	hold, exists := e.Holds[holdID]

	if !exists {
		return ErrHoldNotFound
	}

	e.Sold += hold.Quantity
	// TODO: EMIT AN EVENT SO THAT LATER A NEW ORDER IS CREATED
	delete(e.Holds, holdID)
	return nil
}

func (e *Event) cancelHold(holdID HoldID, currentTime time.Time) error {
	_, exists := e.Holds[holdID]

	if !exists {
		return ErrHoldNotFound
	}

	delete(e.Holds, holdID)
	return nil
}

func CreateNewEvent(title string, capacity int) *Event {
	eventID := EventID(uuid.New())
	return &Event{
		Id:       eventID,
		Title:    title,
		Capacity: capacity,
		Sold:     0,
		Holds:    make(map[HoldID]*Hold),
	}
}

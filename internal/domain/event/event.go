package event

import (
	"time"

	"github.com/google/uuid"
)

type EventID uuid.UUID
type HoldID uuid.UUID
type CustomerID uuid.UUID

type Hold struct {
	id HoldID
	quantity int
	customerID CustomerID
	expiresAt time.Time
}

type Event struct {
	id EventID
	title string
	capacity int
	sold int
	holds map[HoldID]*Hold
}
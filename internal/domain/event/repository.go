package event

import "context"

type repository interface {
	FindByID(ctx context.Context, id EventID) (*Event, error)
	Save(ctx context.Context, event *Event) error
}

package forecast

import "context"

type Reader interface {
	Read(ctx context.Context, slot *Slot) ([]*Weather, error)
}

type Repository interface {
	Save(ctx context.Context, weather *Weather) error
}

package calendar

import (
	"context"
	"github.com/3n0ugh/kalenderium/pkg/calendar/repository"
)

type Service interface {
	CreateEvent(ctx context.Context, event repository.Event) (uint64, error)
	ListEvent(ctx context.Context, userId uint64) ([]*repository.Event, error)
	DeleteEvent(ctx context.Context, eventId uint64) error
	ServiceStatus(ctx context.Context) error
}

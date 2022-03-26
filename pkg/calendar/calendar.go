package calendar

import (
	"context"
	"fmt"
	"github.com/3n0ugh/kalenderium/internal/validator"
	"github.com/3n0ugh/kalenderium/pkg/calendar/repository"
	"github.com/go-kit/log"
	"github.com/pkg/errors"
	"net/http"
	"os"
	"time"
)

type calendarService struct {
	calendarRepository repository.CalendarRepository
}

func NewService(calendarRepository repository.CalendarRepository) Service {
	return &calendarService{
		calendarRepository: calendarRepository,
	}
}

// CreateEvent -> Add the given event to calendar database and returns the eventId
func (c *calendarService) CreateEvent(ctx context.Context, event repository.Event) (uint64, error) {
	v := validator.New()
	if repository.ValidateEvent(v, event); !v.Valid() {
		logger.Log(fmt.Sprintf("validation error: %v", v.Errors))
		return 0, errors.New(fmt.Sprintf("%v", v.Errors))
	}

	err := c.calendarRepository.CreateEvent(ctx, &event)
	if err != nil {
		logger.Log("failed to create event", time.Now())
		return 0, errors.Wrap(err, "failed to create event")
	}

	return event.Id, nil
}

// ListEvent -> Get events from database according to userId and return events
func (c *calendarService) ListEvent(ctx context.Context, userId uint64) ([]repository.Event, error) {
	events, err := c.calendarRepository.ListEvent(ctx, userId)
	if err != nil {
		logger.Log("failed to get events", time.Now())
		return nil, errors.Wrap(err, "failed to get events")
	}
	return events, nil
}

// DeleteEvent -> Delete event from database according to eventId
func (c *calendarService) DeleteEvent(ctx context.Context, eventId uint64, userId uint64) error {
	err := c.calendarRepository.DeleteEvent(ctx, eventId, userId)
	if err != nil {
		logger.Log("failed to delete event", time.Now())
		return errors.Wrap(err, "failed to delete event")
	}
	return nil
}

// ServiceStatus -> A health-check mechanism
func (c *calendarService) ServiceStatus(ctx context.Context) (int, error) {
	if err := c.calendarRepository.ServiceStatus(ctx); err != nil {
		logger.Log("calendar service status error", time.Now())
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}

var logger log.Logger

func init() {
	logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)
}

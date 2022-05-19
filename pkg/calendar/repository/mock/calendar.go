package mock

import (
	"context"
	"fmt"
	"github.com/3n0ugh/kalenderium/pkg/calendar/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

var Event = &repository.Event{
	Id:      primitive.NewObjectID(),
	UserId:  22,
	Name:    "Spring Time",
	Details: "Spring adds new life and new beauty to all that is.",
	Start:   time.Time{}.Add(time.Second),
	End:     time.Time{}.Add(time.Second),
	Color:   "#00FF00",
}

var Event2 = &repository.Event{
	Id:      primitive.NewObjectID(),
	UserId:  22,
	Name:    "Summer Time",
	Details: "Oh, the summer night, has a smile of light, and she sits on a sapphire throne",
	Start:   time.Time{}.Add(time.Second),
	End:     time.Time{}.Add(time.Second),
	Color:   "#00FFFF",
}

type CalendarRepository interface {
	CreateEvent(ctx context.Context, event *repository.Event) error
	ListEvent(ctx context.Context, userId uint64) ([]repository.Event, error)
	DeleteEvent(ctx context.Context, eventId string, userId uint64) error
	ServiceStatus(ctx context.Context) error
}

type calendarRepository struct {
}

func NewCalendarRepository() CalendarRepository {
	return &calendarRepository{}
}

func (c *calendarRepository) CreateEvent(_ context.Context, event *repository.Event) error {
	event.Id, _ = primitive.ObjectIDFromHex("6285f86bb502f9d335124b04")
	fmt.Println(primitive.NewObjectID())
	return nil
}

func (c *calendarRepository) ListEvent(_ context.Context, userId uint64) ([]repository.Event, error) {
	if userId == Event.UserId {
		var eventList = []repository.Event{
			*Event,
			*Event2,
		}
		return eventList, nil
	}
	return nil, repository.ErrRecordNotFound
}

func (c *calendarRepository) DeleteEvent(_ context.Context, eventId string, userId uint64) error {
	if userId == Event.UserId && eventId == Event.Id.String() {
		return nil
	}
	return repository.ErrRecordNotFound
}

func (c *calendarRepository) ServiceStatus(_ context.Context) error {
	return nil
}

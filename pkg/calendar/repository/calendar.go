package repository

import (
	"context"
	"github.com/3n0ugh/kalenderium/internal/validator"
	db "github.com/3n0ugh/kalenderium/pkg/calendar/database"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"strings"
	"time"
)

var ErrRecordNotFound = errors.New("record not found")

type Event struct {
	Id      primitive.ObjectID `json:"id"`
	UserId  uint64             `json:"user_id"`
	Name    string             `json:"name"`
	Details string             `json:"details,omitempty"`
	Start   time.Time          `json:"start"`
	End     time.Time          `json:"end"`
	Color   string             `json:"color"`
}

type CalendarRepository interface {
	CreateEvent(ctx context.Context, event *Event) error
	ListEvent(ctx context.Context, userId uint64) ([]Event, error)
	DeleteEvent(ctx context.Context, eventId string, userId uint64) error
	ServiceStatus(ctx context.Context) error
}

type calendarRepository struct {
	collection *mongo.Collection
	db         *mongo.Client
}

func NewCalendarRepository(conn db.Connection) CalendarRepository {
	coll := conn.DB().Database("kalenderium").Collection("calendar")
	return &calendarRepository{
		collection: coll,
		db:         conn.DB(),
	}
}

func ValidateEvent(v *validator.Validator, event Event) {
	v.Check(event.Name != "", "name", "must be provided")
	v.Check(len(event.Name) <= 80, "name", "must not be more than 80 bytes long")

	v.Check(event.Color != "", "color", "must be provided")
	v.Check(strings.HasPrefix(event.Color, "#"), "color", "must be start with #")
	v.Check(len(event.Color) == 7, "color", "must be 7 bytes long")

	v.Check(len(event.Details) <= 1100, "details", "must not be more than 1100 bytes long")

	v.Check(!event.Start.IsZero(), "start", "must be provided")
	v.Check(!event.End.IsZero(), "end", "must be provided")
}

// CreateEvent -> Adds event to the events database with given userId
func (c *calendarRepository) CreateEvent(ctx context.Context, event *Event) error {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	e := bson.D{
		{"user_id", event.UserId},
		{"name", event.Name},
		{"details", event.Details},
		{"start", event.Start},
		{"end", event.End},
		{"color", event.Color},
	}

	i, err := c.collection.InsertOne(ctx, e)
	if err != nil {
		return errors.Wrap(err, "failed to insert event")
	}

	event.Id = i.InsertedID.(primitive.ObjectID)
	return nil
}

// ListEvent -> Gets events from database according to given userId
func (c *calendarRepository) ListEvent(ctx context.Context, userId uint64) ([]Event, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	cursor, err := c.collection.Find(ctx, bson.D{{"user_id", userId}})
	if err != nil {
		return nil, errors.Wrap(err, "failed to get events from database")
	}
	if cursor.Err() != nil {
		return nil, errors.Wrap(cursor.Err(), "failed to get events from database")
	}

	var events []Event
	var results []bson.M
	if err = cursor.All(ctx, &results); err != nil {
		return nil, errors.Wrap(err, "failed to get events from database")
	}
	for _, result := range results {
		event := Event{
			Id:      result["_id"].(primitive.ObjectID),
			UserId:  uint64(result["user_id"].(int64)),
			Name:    result["name"].(string),
			Details: result["details"].(string),
			Start:   result["start"].(primitive.DateTime).Time(),
			End:     result["end"].(primitive.DateTime).Time(),
			Color:   result["color"].(string),
		}
		events = append(events, event)
	}

	return events, nil
}

// DeleteEvent -> Deletes event according to given userId and eventId
func (c *calendarRepository) DeleteEvent(ctx context.Context, eventId string, userId uint64) error {
	id, err := primitive.ObjectIDFromHex(eventId)
	if err != nil {
		return errors.Wrap(err, "failed to convert eventId to ObjectID")
	}
	return c.collection.FindOneAndDelete(ctx, bson.M{"user_id": userId, "_id": id}).Err()
}

// ServiceStatus -> A health check mechanism
func (c *calendarRepository) ServiceStatus(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	return c.db.Ping(ctx, readpref.Primary())
}

package repository

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/3n0ugh/kalenderium/internal/validator"
	db "github.com/3n0ugh/kalenderium/pkg/calendar/database"
	"github.com/pkg/errors"
	"time"
)

var ErrRecordNotFound = errors.New("record not found")

type Event struct {
	EventId   uint64    `json:"event_id"`
	UserId    uint64    `json:"user_id"`
	Title     string    `json:"title"`
	Body      string    `json:"body,omitempty"`
	AttendAt  time.Time `json:"attend_at"`
	CreatedAt time.Time `json:"created_at"`
}

type CalendarRepository interface {
	CreateEvent(ctx context.Context, event *Event) error
	ListEvent(ctx context.Context, userId uint64) ([]*Event, error)
	DeleteEvent(ctx context.Context, eventId uint64, userId uint64) error
	ServiceStatus(ctx context.Context) error
}

type calendarRepository struct {
	db *sql.DB
}

func NewCalendarRepository(conn db.Connection) CalendarRepository {
	return &calendarRepository{db: conn.DB()}
}

func ValidateEvent(v *validator.Validator, event Event) {
	v.Check(event.Title != "", "title", "must be provided")
	v.Check(len(event.Title) <= 80, "title", "must not be more than 80 bytes long")

	v.Check(len(event.Body) <= 1100, "body", "must not be more than 1100 bytes long")

	v.Check(time.Now().Before(event.AttendAt), "attend_at", "out of date")
}

// CreateEvent -> Adds event to the events database with given userId
func (c *calendarRepository) CreateEvent(ctx context.Context, event *Event) error {
	query := `INSERT INTO events (user_id, title, body, attend_at)
			VALUES ($1, $2, $3, $4)
			RETURNING event_id`

	args := []interface{}{event.UserId, event.Title, event.Body, event.AttendAt}

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	return c.db.QueryRowContext(ctx, query, args...).Scan(
		&event.EventId)
}

// ListEvent -> Gets events from database according to given userId
func (c *calendarRepository) ListEvent(ctx context.Context, userId uint64) ([]*Event, error) {
	query := `SELECT event_id, user_id, title, body, attend_at, created_at FROM events
			WHERE user_id = $1`

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	rows, err := c.db.QueryContext(ctx, query, userId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrRecordNotFound
		}
		return nil, err
	}

	defer func() {
		err = rows.Close()
		if err != nil {
			panic(err)
		}
	}()

	var events = make([]*Event, 1)

	for rows.Next() {
		var event Event
		err = rows.Scan(
			&event.EventId,
			&event.UserId,
			&event.Title,
			&event.Body,
			&event.AttendAt,
			&event.CreatedAt,
		)

		if err != nil {
			return nil, err
		}

		events = append(events, &event)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return events, nil
}

// DeleteEvent -> Deletes event according to given userId and eventId
func (c *calendarRepository) DeleteEvent(ctx context.Context, eventId uint64, userId uint64) error {
	query := `DELETE FROM events 
			WHERE event_id = $1 AND user_id = $2`

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	res, err := c.db.ExecContext(ctx, query, eventId, userId)
	if err != nil {
		return err
	}

	effectedRow, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if effectedRow == 0 {
		return ErrRecordNotFound
	}

	if effectedRow > 1 {
		return errors.New(fmt.Sprintf("expected to affect 1 row, affected %d", effectedRow))
	}

	return nil
}

// ServiceStatus -> A health check mechanism
func (c *calendarRepository) ServiceStatus(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	return c.db.PingContext(ctx)
}

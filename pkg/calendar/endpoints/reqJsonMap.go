package endpoints

import "github.com/3n0ugh/kalenderium/pkg/calendar/repository"

// CreateEventRequest -> CreateEvent endpoint's  input structures
type CreateEventRequest struct {
	Event repository.Event `json:"event"`
}

// CreateEventResponse -> CreateEvent endpoint's output structure
type CreateEventResponse struct {
	EventId uint64 `json:"event_id,omitempty"`
	Err     string `json:"err,omitempty"`
}

// ListEventRequest -> ListEvent endpoint's  input structures
type ListEventRequest struct {
	UserId uint64 `json:"user_id"`
}

// ListEventResponse -> ListEvent endpoint's output structure
type ListEventResponse struct {
	Events []repository.Event `json:"events,omitempty"`
	Err    string             `json:"err,omitempty"`
}

// DeleteEventRequest -> DeleteEvent endpoint's  input structures
type DeleteEventRequest struct {
	EventId uint64 `json:"event_id"`
	UserId  uint64 `json:"user_id"`
}

// DeleteEventResponse -> DeleteEvent endpoint's output structure
type DeleteEventResponse struct {
	Err string `json:"err,omitempty"`
}

// ServiceStatusRequest -> ServiceStatus endpoint's  input structures
type ServiceStatusRequest struct{}

// ServiceStatusResponse -> ServiceStatus endpoint's output structure
type ServiceStatusResponse struct {
	Code int    `json:"code,omitempty"`
	Err  string `json:"err,omitempty"`
}

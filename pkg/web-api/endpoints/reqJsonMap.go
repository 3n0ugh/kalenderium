package endpoints

import (
	"github.com/3n0ugh/kalenderium/internal/token"
	repo "github.com/3n0ugh/kalenderium/pkg/account/repository"
	"github.com/3n0ugh/kalenderium/pkg/calendar/repository"
)

// AddEventRequest -> CreateEvent endpoint's  input structures
type AddEventRequest struct {
	Event repository.Event `json:"event"`
}

// AddEventResponse -> CreateEvent endpoint's output structure
type AddEventResponse struct {
	EventId uint64 `json:"event_id,omitempty"`
	Err     string `json:"err,omitempty"`
}

// ListEventRequest -> CreateEvent endpoint's  input structures
type ListEventRequest struct {
	UserId uint64 `json:"user_id"`
}

// ListEventResponse -> CreateEvent endpoint's output structure
type ListEventResponse struct {
	Events []repository.Event `json:"events,omitempty"`
	Err    string             `json:"err,omitempty"`
}

// DeleteEventRequest -> CreateEvent endpoint's  input structures
type DeleteEventRequest struct {
	EventId uint64 `json:"eventId,omitempty"`
	UserId  uint64 `json:"userId"`
}

// DeleteEventResponse -> CreateEvent endpoint's output structure
type DeleteEventResponse struct {
	Err string `json:"err,omitempty"`
}

// SignUpRequest -> CreateEvent endpoint's  input structures
type SignUpRequest struct {
	User repo.User `json:"user"`
}

// SignUpResponse -> CreateEvent endpoint's output structure
type SignUpResponse struct {
	UserId uint64      `json:"userId,omitempty"`
	Token  token.Token `json:"token,omitempty"`
	Err    string      `json:"err,omitempty"`
}

// LoginRequest -> CreateEvent endpoint's  input structures
type LoginRequest struct {
	User repo.User `json:"user"`
}

// LoginResponse -> CreateEvent endpoint's output structure
type LoginResponse struct {
	UserId uint64      `json:"userId,omitempty"`
	Token  token.Token `json:"token,omitempty"`
	Err    string      `json:"err,omitempty"`
}

// LogoutRequest -> CreateEvent endpoint's  input structures
type LogoutRequest struct {
	Token token.Token `json:"token"`
}

// LogoutResponse -> CreateEvent endpoint's output structure
type LogoutResponse struct {
	Err string `json:"err,omitempty"`
}

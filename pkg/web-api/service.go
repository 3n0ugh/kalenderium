package web_api

import (
	"context"
	"github.com/3n0ugh/kalenderium/internal/token"
	repo "github.com/3n0ugh/kalenderium/pkg/account/repository"
	"github.com/3n0ugh/kalenderium/pkg/calendar/repository"
)

type Service interface {
	AddEvent(ctx context.Context, event repository.Event) (string, error)
	ListEvent(ctx context.Context, userId uint64) ([]repository.Event, error)
	DeleteEvent(ctx context.Context, eventId string, userId uint64) error

	SignUp(ctx context.Context, user repo.User) (uint64, token.Token, error)
	Login(ctx context.Context, user repo.User) (uint64, token.Token, error)
	Logout(ctx context.Context, token token.Token) error
}

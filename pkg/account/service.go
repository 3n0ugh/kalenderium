package account

import (
	"context"
	"github.com/3n0ugh/kalenderium/internal/token"
	"github.com/3n0ugh/kalenderium/pkg/account/repository"
)

type Service interface {
	IsAuth(ctx context.Context, token token.Token) (token.Token, error)
	SignUp(ctx context.Context, user repository.User) (uint64, token.Token, error)
	Login(ctx context.Context, user repository.User) (uint64, token.Token, error)
	Logout(ctx context.Context, token token.Token) error
	ServiceStatus(ctx context.Context) (int, error)
}

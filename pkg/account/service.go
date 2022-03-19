package account

import (
	"context"
	"github.com/3n0ugh/kalenderium/pkg/account/repository"
)

type Service interface {
	IsAuth(ctx context.Context, token string) error
	SignUp(ctx context.Context, user repository.User) (string, error)
	Login(ctx context.Context, user repository.User) (string, error)
	Logout(ctx context.Context, token string) error
}

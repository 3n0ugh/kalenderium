package mock

import (
	"context"
	"github.com/3n0ugh/kalenderium/pkg/account/repository"
	"github.com/pkg/errors"
)

var (
	ErrDuplicateEmail = errors.New("duplicate email")
	ErrRecordNotFound = errors.New("record not found")
)

type AccountRepository interface {
	CreateUser(ctx context.Context, user *repository.User) error
	GetUser(ctx context.Context, email string) (*repository.User, error)
	ServiceStatus(ctx context.Context) error
}

type accountRepository struct {
}

var User = &repository.User{
	UserID:       22,
	Email:        "test@22.com",
	Password:     "test_1234!",
	PasswordHash: []byte("$2a$12$Wck7TSDfSYhn0GxOchYEJe3xX5w3MItslZPFUHRGbTYIQxjgHPHCe"),
}

func NewAccountRepository() AccountRepository {
	return &accountRepository{}
}

func (a *accountRepository) CreateUser(_ context.Context, user *repository.User) error {
	if user.Email == User.Email {
		return ErrDuplicateEmail
	}
	user.UserID = uint64(23)
	return nil
}

func (a *accountRepository) GetUser(_ context.Context, email string) (*repository.User, error) {
	if email != User.Email {
		return nil, ErrRecordNotFound
	}
	return User, nil
}

func (a *accountRepository) ServiceStatus(_ context.Context) error {
	return nil
}

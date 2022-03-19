package account

import (
	"context"
	"fmt"
	"github.com/3n0ugh/kalenderium/internal/validator"
	"github.com/3n0ugh/kalenderium/pkg/account/repository"
	"github.com/3n0ugh/kalenderium/pkg/account/store"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"strconv"
)

type accountService struct {
	accountRepository repository.AccountRepository
	serializableStore store.SerializableStore
}

func NewService(accountRepository repository.AccountRepository, customRedisStore store.SerializableStore) Service {
	return &accountService{
		accountRepository: accountRepository,
		serializableStore: customRedisStore,
	}
}

func (a *accountService) IsAuth(ctx context.Context, token string) error {
	return nil
}

// SignUp creates new user and session token and return session token
func (a *accountService) SignUp(ctx context.Context, user repository.User) (string, error) {
	// Hash the user's plain-text password
	err := user.Set(user.Password)
	if err != nil {
		return "", errors.Wrap(err, "failed to hash password")
	}

	v := validator.New()
	repository.ValidateUser(v, &user)
	if !v.Valid() {
		return "", errors.New(fmt.Sprintf("failed user data validation: %v", v.Errors))
	}

	// Add new user to account database
	err = a.accountRepository.CreateUser(ctx, user)
	if err != nil {
		return "", errors.Wrap(err, "failed to create new user")
	}

	// New session token
	token := uuid.New().String()

	// Add session token to Redis
	err = a.serializableStore.Set(ctx, strconv.FormatUint(user.UserID, 10), token)
	if err != nil {
		return "", errors.Wrap(err, "failed to set session token to Redis")
	}

	return token, nil
}

func (a *accountService) Login(ctx context.Context, user repository.User) (string, error) {
	return "", nil
}
func (a *accountService) Logout(ctx context.Context, token string) error {
	return nil
}
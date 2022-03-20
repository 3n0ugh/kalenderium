package account

import (
	"context"
	"fmt"
	"github.com/3n0ugh/kalenderium/internal/validator"
	"github.com/3n0ugh/kalenderium/pkg/account/repository"
	"github.com/3n0ugh/kalenderium/pkg/account/store"
	"github.com/google/uuid"
	"github.com/pkg/errors"
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

// IsAuth checks the redis for the given token is existed or not
func (a *accountService) IsAuth(ctx context.Context, token string) error {
	// Check token is valid uuid
	if _, err := uuid.Parse(token); err != nil {
		return errors.Wrap(err, "invalid uuid")
	}

	// Get session info from redis
	_, err := a.serializableStore.Get(ctx, token)
	if err != nil {
		return errors.Wrap(err, "session is not available")
	}

	return nil
}

// SignUp creates a new user and session token and returns session token
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
	err = a.accountRepository.CreateUser(ctx, &user)
	if err != nil {
		return "", errors.Wrap(err, "failed to create new user")
	}

	// New session token
	token := uuid.New().String()

	// Add session token to Redis
	err = a.serializableStore.Set(ctx, user.UserID, token)
	if err != nil {
		return "", errors.Wrap(err, "failed to set session token to Redis")
	}

	return token, nil
}

// Login checks are given user exist in the database, if exist return session token
func (a *accountService) Login(ctx context.Context, user repository.User) (string, error) {
	err := user.Set(user.Password)
	if err != nil {
		return "", errors.Wrap(err, "failed to hash password")
	}

	v := validator.New()
	repository.ValidateUser(v, &user)
	if !v.Valid() {
		return "", errors.New(fmt.Sprintf("failed user data validation: %v", v.Errors))
	}

	// Get user from account database
	usr, err := a.accountRepository.GetUser(ctx, user.Email)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			return "", errors.Wrap(err, "user not found")
		}
		return "", err
	}

	// Compare password hashes
	err = usr.Matches(user.Password)
	if err != nil {
		return "", errors.Wrap(err, "wrong password")
	}

	// New session token
	token := uuid.New().String()

	// Add session token to Redis
	err = a.serializableStore.Set(ctx, usr.UserID, token)
	if err != nil {
		return "", errors.Wrap(err, "failed to set session token to Redis")
	}

	return token, nil
}

// Logout removes session token from redis
func (a *accountService) Logout(ctx context.Context, token string) error {
	// Check token is valid uuid
	if _, err := uuid.Parse(token); err != nil {
		return errors.Wrap(err, "invalid uuid")
	}

	// Delete token from redis
	if err := a.serializableStore.Delete(ctx, token); err != nil {
		return errors.Wrap(err, "failed to delete session token from redis")
	}
	return nil
}

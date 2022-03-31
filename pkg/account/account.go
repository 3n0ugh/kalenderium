package account

import (
	"context"
	"crypto/sha256"
	"fmt"
	"github.com/3n0ugh/kalenderium/internal/token"
	"github.com/3n0ugh/kalenderium/internal/validator"
	"github.com/3n0ugh/kalenderium/pkg/account/repository"
	"github.com/3n0ugh/kalenderium/pkg/account/store"
	"github.com/go-kit/log"
	"github.com/pkg/errors"
	"net/http"
	"os"
	"time"
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
func (a *accountService) IsAuth(ctx context.Context, sessionToken token.Token) (token.Token, error) {
	// Check token is valid uuid
	v := validator.New()
	token.ValidateTokenPlaintext(v, sessionToken.PlainText)
	if !v.Valid() {
		logger.Log("failed to validate token")
		return token.Token{}, errors.New(fmt.Sprintf("failed to validate token: %v", v.Errors))
	}

	// Get session info from redis
	tkn, err := a.serializableStore.Get(ctx, sessionToken.PlainText)
	if err != nil {
		logger.Log("session is not available")
		return token.Token{}, errors.New("session is not available")
	}

	return tkn, nil
}

// SignUp creates a new user and session token and returns session token
func (a *accountService) SignUp(ctx context.Context, user repository.User) (uint64, token.Token, error) {
	// Hash the user's plain-text password
	err := user.Set(user.Password)
	if err != nil {
		logger.Log("failed to hash password")
		return 0, token.Token{}, errors.New("failed to hash password")
	}

	// Validate the user
	v := validator.New()
	repository.ValidateUser(v, &user)
	if !v.Valid() {
		logger.Log("failed user data validation: %v", v.Errors)
		return 0, token.Token{}, errors.New(fmt.Sprintf("failed user data validation: %v", v.Errors))
	}

	// Add new user to account database
	err = a.accountRepository.CreateUser(ctx, &user)
	if err != nil {
		logger.Log("failed to create new user")
		return 0, token.Token{}, errors.New("failed to create new user")
	}

	// New session token
	sessionToken, err := token.GenerateToken(user.UserID, time.Minute*60, token.ScopeAuthentication)
	if err != nil {
		logger.Log("failed to generate token")
		return 0, token.Token{}, errors.New("failed to generate token")
	}

	// Add session token to Redis
	err = a.serializableStore.Set(ctx, sessionToken)
	if err != nil {
		logger.Log("failed to set session token to redis")
		return 0, token.Token{}, errors.New("failed to set session token to redis")
	}

	return user.UserID, *sessionToken, nil
}

// Login checks are given user exist in the database, if exist return session token
func (a *accountService) Login(ctx context.Context, user repository.User) (uint64, token.Token, error) {
	err := user.Set(user.Password)
	if err != nil {
		logger.Log("failed to hash password")
		return 0, token.Token{}, errors.New("failed to hash password")
	}

	// Validate the user
	v := validator.New()
	repository.ValidateUser(v, &user)
	if !v.Valid() {
		logger.Log("failed user data validation")
		return 0, token.Token{}, errors.New(" failed user data validation")

	}

	// Get user from account database
	usr, err := a.accountRepository.GetUser(ctx, user.Email)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			logger.Log("user not found")
			return 0, token.Token{}, errors.New("user not found")
		}
		return 0, token.Token{}, err
	}

	// Compare password hashes
	match, err := usr.Matches(user.Password)
	if err != nil {
		logger.Log("failed to  password and hash")
		return 0, token.Token{}, errors.New("failed to compare password and hash")
	}
	if !match {
		logger.Log("wrong password")
		return 0, token.Token{}, errors.New("wrong password")
	}

	// New session token
	sessionToken, err := token.GenerateToken(usr.UserID, time.Minute*60, token.ScopeAuthentication)
	if err != nil {
		logger.Log("failed to generate token")
		return 0, token.Token{}, errors.New("failed to generate token")
	}

	// Add session token to Redis
	err = a.serializableStore.Set(ctx, sessionToken)
	if err != nil {
		logger.Log("failed to set session token to redis")
		return 0, token.Token{}, errors.New("failed to set session token to redis")
	}

	return usr.UserID, *sessionToken, nil
}

// Logout removes session token from redis
func (a *accountService) Logout(ctx context.Context, sessionToken token.Token) error {
	// Check token is valid uuid
	v := validator.New()
	token.ValidateTokenPlaintext(v, sessionToken.PlainText)
	if !v.Valid() {
		logger.Log("failed to validate token")
		return errors.New("failed to validate token")
	}

	hash := sha256.Sum256([]byte(sessionToken.PlainText))
	sessionToken.Hash = hash[:]

	// Delete token from redis
	if err := a.serializableStore.Delete(ctx, string(sessionToken.Hash)); err != nil {
		logger.Log("failed to delete session token from redis")
		return errors.New("failed to delete session token from redis")
	}
	return nil
}

func (a *accountService) ServiceStatus(ctx context.Context) (int, error) {
	logger.Log("Checking the Service health...")
	err := a.accountRepository.ServiceStatus(ctx)
	if err != nil {
		return http.StatusInternalServerError, errors.New("")
	}
	return http.StatusOK, nil
}

var logger log.Logger

func init() {
	logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)
}

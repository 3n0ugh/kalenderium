package mock

import (
	"context"
	"github.com/3n0ugh/kalenderium/internal/token"
	"github.com/pkg/errors"
	"time"
)

type SerializableStore interface {
	Get(ctx context.Context, sessionTokenHash string) (token.Token, error)
	Set(ctx context.Context, sessionToken *token.Token) error
	Delete(ctx context.Context, token string) error
}

type redisStore struct {
}

var Token = &token.Token{
	PlainText: "GAJPGPWLD6KISED2QS34A6ERWU",
	Hash:      []byte("75c1ea94931bfefe90b8684ca17220a3bf8813aef074c91226bc8fafacd55fb3"),
	UserID:    22,
	Expiry:    time.Time{},
	Scope:     token.ScopeAuthentication,
}

func CustomRedisStore(_ context.Context) SerializableStore {
	return &redisStore{}
}

func (r redisStore) Get(_ context.Context, sessionTokenHash string) (token.Token, error) {
	if sessionTokenHash != Token.PlainText {
		return token.Token{}, errors.New("session not found")
	}
	return *Token, nil
}

func (r redisStore) Set(_ context.Context, sessionToken *token.Token) error {
	sessionToken.PlainText = Token.PlainText
	sessionToken.Hash = Token.Hash
	sessionToken.Expiry = Token.Expiry
	return nil
}

func (r redisStore) Delete(_ context.Context, _ string) error {
	return nil
}

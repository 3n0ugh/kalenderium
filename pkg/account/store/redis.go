package store

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"github.com/3n0ugh/kalenderium/internal/config"
	"github.com/3n0ugh/kalenderium/internal/token"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
)

type SerializableStore interface {
	Get(ctx context.Context, sessionTokenHash string) (token.Token, error)
	Set(ctx context.Context, sessionToken *token.Token) error
	Delete(ctx context.Context, token string) error
}

type redisStore struct {
	client *redis.Client
}

// CustomRedisStore established new Redis connection
func CustomRedisStore(ctx context.Context, cfg config.AccountServiceConfigurations) SerializableStore {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisUrl,
		Password: cfg.RedisPass,
		DB:       0,
	})

	_, err := client.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Failed to ping Redis: %v", err)
	}

	return &redisStore{
		client: client,
	}
}

// Delete removes token from Redis
func (r redisStore) Delete(ctx context.Context, token string) error {
	_, err := r.client.Del(ctx, token).Result()
	if err != nil {
		return errors.New("problem")
	}
	return nil
}

// Get session token from Redis
func (r redisStore) Get(ctx context.Context, sessionToken string) (token.Token, error) {
	hash := sha256.Sum256([]byte(sessionToken))
	tHash := hash[:]
	userSession, err := r.client.Get(ctx, string(tHash)).Result()
	if err != nil {
		return token.Token{}, errors.New("session not found")
	}

	var session token.Token
	err = json.Unmarshal([]byte(userSession), &session)
	if err != nil {
		return token.Token{}, errors.New("failed to unmarshal session")
	}
	return session, nil
}

// Set creates session token for 1 hour
func (r redisStore) Set(ctx context.Context, sessionToken *token.Token) error {
	session, err := json.Marshal(sessionToken)
	if err != nil {
		return errors.New("failed to marshal session")
	}

	err = r.client.Set(ctx, string(sessionToken.Hash), session, time.Minute*60).Err()
	if err != nil {
		return errors.New("failed to save session to redis")
	}
	return nil
}

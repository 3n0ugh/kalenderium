package store

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
)

type UserSession struct {
	CreatedAt time.Time `json:"created_at"`
	UserID    uint64    `json:"user_id"`
}

type SerializableStore interface {
	Get(ctx context.Context, token string) (UserSession, error)
	Set(ctx context.Context, id uint64, token string) error
	Delete(ctx context.Context, token string) error
}

type redisStore struct {
	client *redis.Client
}

// CustomRedisStore established new Redis connection
func CustomRedisStore(ctx context.Context) SerializableStore {
	client := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_URL"),
		Password: "",
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
		return errors.Wrap(err, "problem")
	}
	return nil
}

// Get session token from Redis
func (r redisStore) Get(ctx context.Context, token string) (UserSession, error) {
	userSession, err := r.client.Get(ctx, token).Result()
	if err != nil {
		return UserSession{}, errors.Wrap(err, "session not found")
	}

	var session UserSession
	err = json.Unmarshal([]byte(userSession), &session)
	if err != nil {
		return UserSession{}, errors.Wrap(err, "failed to unmarshal session")
	}
	return session, nil
}

// Set creates session token for 1 hour
func (r redisStore) Set(ctx context.Context, id uint64, token string) error {
	userSession := UserSession{
		CreatedAt: time.Now(),
		UserID:    id,
	}

	session, err := json.Marshal(userSession)
	if err != nil {
		return errors.Wrap(err, "failed to marshal session")
	}

	err = r.client.Set(ctx, token, string(session), time.Minute*60).Err()
	if err != nil {
		fmt.Println(session)
		return errors.Wrap(err, "failed to save session to redis")
	}
	return nil
}

package store

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
)

type UserSession struct {
	CreatedAt time.Time `json:"created_at"`
	UserID    uint64    `json:"user_id"`
	Email     string    `json:"email"`
}

type SerializableStore interface {
	Get(ctx context.Context, id string) (UserSession, error)
	Set(ctx context.Context, id string, session string) error
	Delete(ctx context.Context, id string) error
}

type redisStore struct {
	client *redis.Client
}

// CustomRedisStore established new Redis connection
func CustomRedisStore(ctx context.Context) SerializableStore {
	client := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_URL"),
		Password: os.Getenv("REDIS_PASSWORD"),
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
func (r redisStore) Delete(ctx context.Context, id string) error {
	_, err := r.client.Del(ctx, id).Result()
	if err != nil {
		return errors.Wrap(err, "problem")
	}
	return nil
}

// Get session token from Redis
func (r redisStore) Get(ctx context.Context, id string) (UserSession, error) {
	userSession, err := r.Get(ctx, id)
	if err != nil {
		return UserSession{}, err
	}
	return userSession, nil
}

// Set creates session token for 1 hour
func (r redisStore) Set(ctx context.Context, id string, session string) error {
	err := r.client.Set(ctx, id, session, time.Minute*60).Err()
	if err != nil {
		return errors.Wrap(err, "failed to save session to redis")
	}

	return nil
}

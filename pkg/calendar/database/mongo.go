package database

import (
	"context"
	"github.com/3n0ugh/kalenderium/internal/config"
	_ "github.com/lib/pq"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"time"
)

type Connection interface {
	Close()
	DB() *mongo.Client
}

type conn struct {
	database *mongo.Client
}

func NewConnection(cfg config.CalendarServiceConfigurations) (Connection, error) {
	// Create an empty connection pool
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(cfg.URI))
	if err != nil {
		return nil, err
	}

	// Create a context with a 5-second timeout deadline
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// If the connection couldn't be established successfully
	// within the 5-second deadline, then this will return an error
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return nil, err
	}

	return &conn{database: client}, nil
}

func (c *conn) Close() {
	c.Close()
}

func (c *conn) DB() *mongo.Client {
	return c.database
}

package database

import (
	"context"
	"database/sql"
	"github.com/3n0ugh/kalenderium/internal/config"
	"time"
)

type Connection interface {
	Close()
	DB() *sql.DB
}

type conn struct {
	database *sql.DB
}

func NewConnection(cfg config.AccountServiceConfigurations) (Connection, error) {
	// Create an empty connection pool
	db, err := sql.Open("mysql", cfg.DSN)
	if err != nil {
		return nil, err
	}

	duration, err := time.ParseDuration(cfg.MaxIdleTime)
	if err != nil {
		return nil, err
	}

	// Database connection configs
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxIdleTime(duration)

	// Create a context with a 5-second timeout deadline
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// If the connection couldn't be established successfully
	// within the 5-second deadline, then this will return an error
	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return &conn{database: db}, nil
}

func (c *conn) Close() {
	c.Close()
}

func (c *conn) DB() *sql.DB {
	return c.database
}

package database

import (
	"fmt"
	"os"
	"strconv"
)

type Config interface {
	Dsn() string
	DbName() string
	DbMaxOpenConns() int
	DbMaxIdleConns() int
	DbMaxIdleTime() string
}

type config struct {
	dbUser       string
	dbPass       string
	dbHost       string
	dbPort       int
	dbName       string
	dsn          string
	maxOpenConns int
	maxIdleConns int
	maxIdleTime  string
	sslMode      string
}

func NewConfig() Config {
	var cfg config
	cfg.dbUser = os.Getenv("DATABASE_USER")
	cfg.dbPass = os.Getenv("DATABASE_PASS")
	cfg.dbHost = os.Getenv("DATABASE_HOST")
	cfg.dbName = os.Getenv("DATABASE_NAME")
	cfg.sslMode = os.Getenv("DATABASE_SSL_MODE")
	cfg.maxOpenConns, _ = strconv.Atoi(os.Getenv("DATABASE_MAX_OPEN_CONNS"))
	cfg.maxIdleConns, _ = strconv.Atoi(os.Getenv("DATABASE_MAX_IDLE_CONNS"))
	cfg.maxIdleTime = os.Getenv("DATABASE_MAX_IDLE_TIME")
	cfg.dbPort, _ = strconv.Atoi(os.Getenv("DATABASE_PORT"))
	cfg.dsn = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s",
		cfg.dbHost, cfg.dbUser, cfg.dbPass, cfg.dbName, cfg.dbPort, cfg.sslMode)
	return &cfg
}

func (c *config) Dsn() string {
	return c.dsn
}

func (c *config) DbName() string {
	return c.dbName
}

func (c *config) DbMaxOpenConns() int {
	return c.maxOpenConns
}

func (c *config) DbMaxIdleConns() int {
	return c.maxIdleConns
}

func (c *config) DbMaxIdleTime() string {
	return c.maxIdleTime
}

package database

import (
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
	cfg.dbUser = os.Getenv("ACCOUNT_DATABASE_USER")
	cfg.dbPass = os.Getenv("ACCOUNT_DATABASE_PASS")
	cfg.dbHost = os.Getenv("ACCOUNT_DATABASE_HOST")
	cfg.dbName = os.Getenv("ACCOUNT_DATABASE_NAME")
	cfg.sslMode = os.Getenv("ACCOUNT_DATABASE_SSL_MODE")
	cfg.maxOpenConns, _ = strconv.Atoi(os.Getenv("ACCOUNT_DATABASE_MAX_OPEN_CONNS"))
	cfg.maxIdleConns, _ = strconv.Atoi(os.Getenv("ACCOUNT_DATABASE_MAX_IDLE_CONNS"))
	cfg.maxIdleTime = os.Getenv("ACCOUNT_DATABASE_MAX_IDLE_TIME")
	cfg.dbPort, _ = strconv.Atoi(os.Getenv("ACCOUNT_DATABASE_PORT"))
	cfg.dsn = os.Getenv("ACCOUNT_DB_DSN")
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

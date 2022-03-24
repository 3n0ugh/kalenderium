package config

import (
	"fmt"
	"os"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

const (
	defaultEnv = "dev"
)

type Config struct {
	Calendar CalendarServiceConfigurations `mapstructure:"calendar_service"`
	Account  AccountServiceConfigurations  `mapstructure:"account_service"`
}

type CalendarServiceConfigurations struct {
	DBUser       string `mapstructure:"db_user"`
	DBPass       string `mapstructure:"db_pass"`
	DBHost       string `mapstructure:"db_host"`
	DBPort       int    `mapstructure:"db_port"`
	DBName       string `mapstructure:"db_name"`
	DSN          string `mapstructure:"db_dsn"`
	MaxOpenConns int    `mapstructure:"db_max_open_conns"`
	MaxIdleConns int    `mapstructure:"db_max_idle_conns"`
	MaxIdleTime  string `mapstructure:"db_max_idle_time"`
	SSLMode      string `mapstructure:"db_ssl_mode"`
}

type AccountServiceConfigurations struct {
	DBUser       string `mapstructure:"db_user"`
	DBPass       string `mapstructure:"db_pass"`
	DBHost       string `mapstructure:"db_host"`
	DBPort       int    `mapstructure:"db_port"`
	DBName       string `mapstructure:"db_name"`
	DSN          string `mapstructure:"db_dsn"`
	MaxOpenConns int    `mapstructure:"db_max_open_conns"`
	MaxIdleConns int    `mapstructure:"db_max_idle_conns"`
	MaxIdleTime  string `mapstructure:"db_max_idle_time"`
	SSLMode      string `mapstructure:"db_ssl_mode"`
	RedisUrl     string `mapstructure:"redis_url"`
	RedisPass    string `mapstructure:"redis_pass"`
}

var cfgReader *configReader

type (
	configReader struct {
		configFile string
		v          *viper.Viper
	}
)

// GetValueByKey gets value by key from the  config
func GetValueByKey(key string) (string, error) {
	newConfigReader()

	var err error
	if err = cfgReader.v.ReadInConfig(); err != nil {
		return "", errors.Wrap(err, "failed to read config file")
	}
	return cfgReader.v.GetString(key), nil
}

// GetConfigByKey gets config to given struct by key
func GetConfigByKey(key string, config interface{}) error {
	newConfigReader()

	var err error
	if err = cfgReader.v.ReadInConfig(); err != nil {
		return errors.Wrap(err, "failed to read config file")
	}

	if err = cfgReader.v.UnmarshalKey(key, config); err != nil {
		return errors.Wrap(err, "failed to unmarshal config")
	}
	return nil
}

func GetConfig(config *Config) error {
	newConfigReader()

	var err error
	if err = cfgReader.v.ReadInConfig(); err != nil {
		return errors.Wrap(err, "failed to read config file")
	}

	if err = cfgReader.v.Unmarshal(config); err != nil {
		return errors.Wrap(err, "failed to unmarshal config")
	}
	return nil
}

// getEnvironment gets environment if fail return fallback
func getEnvironment(env, fallback string) string {
	e := os.Getenv(env)
	if e == "" {
		return fallback
	}
	return e
}

// newConfigReader creates new viper config reader.
func newConfigReader() {
	env := getEnvironment("APP_ENVIRONMENT", defaultEnv)
	configFile := fmt.Sprintf("api.%s.yaml", env)

	v := viper.GetViper()
	v.SetConfigType("yaml")
	v.SetConfigFile(configFile)
	v.AddConfigPath(".")

	cfgReader = &configReader{
		configFile: configFile,
		v:          v,
	}
}
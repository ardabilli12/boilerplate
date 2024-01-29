package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Environment    string `mapstructure:"ENVIRONMENT"`
	ServiceName    string `mapstructure:"SERVICE_NAME"`
	ServiceVersion string `mapstructure:"SERVICE_VERSION"`

	MySQLConfig    MySQLConfig    `mapstructure:"DATABASE"`
	RedisConfig    RedisConfig    `mapstructure:"REDIS"`
	RabbitMQConfig RabbitMQConfig `mapstructure:"RABBITMQ"`
}

type Secret struct {
	MySQLSecret    MySQLSecret    `mapstructure:"DATABASE"`
	RedisSecret    RedisSecret    `mapstructure:"REDIS"`
	RabbitMQSecret RabbitMQSecret `mapstructure:"RABBITMQ"`
	SecuritySecret SecuritySecret `mapstructure:"SECURITY"`

	NewRelicLicenseKey string `mapstructure:"NEW_RELIC_LICENSE_KEY"`
}

type MySQLConfig struct {
	Dialect      string `mapstructure:"DIALECT"`
	Host         string `mapstructure:"HOST"`
	Port         string `mapstructure:"PORT"`
	MaxIdleConns int    `mapstructure:"MAX_IDLE_CONNS"`
	MaxOpenConns int    `mapstructure:"MAX_OPEN_CONNS"`
	MaxIdleTime  int    `mapstructure:"MAX_IDLE_TIME"`
	MaxLifeTime  int    `mapstructure:"MAX_LIFE_TIME"`
}

type MySQLSecret struct {
	Database string `mapstructure:"DB_NAME"`
	Username string `mapstructure:"USERNAME"`
	Password string `mapstructure:"PASSWORD"`
}

type RedisConfig struct {
	Host    string `mapstructure:"HOST"`
	Port    string `mapstructure:"PORT"`
	CacheDB int    `mapstructure:"CACHE_DB"`
}

type RedisSecret struct {
	Password string `mapstructure:"PASSWORD"`
}

type RabbitMQConfig struct {
	Host string `mapstructure:"HOST"`
	Port string `mapstructure:"PORT"`
}

type RabbitMQSecret struct {
	Username string `mapstructure:"USERNAME"`
	Password string `mapstructure:"PASSWORD"`
}

type SecuritySecret struct {
	CardSecretKey string `mapstructure:"CARD_SECRET_KEY"`
	CardIV        string `mapstructure:"CARD_IV"`
}

func LoadConfig(configPath, secretPath string) (*Config, *Secret, error) {
	// Load Config
	viper.SetConfigFile(configPath)
	if err := viper.ReadInConfig(); err != nil {
		return nil, nil, fmt.Errorf("error reading config file: %w", err)
	}
	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, nil, fmt.Errorf("unable to decode into struct: %w", err)
	}

	viper.Reset()

	// Load Secret
	viper.SetConfigFile(secretPath)
	var secret Secret
	if err := viper.ReadInConfig(); err != nil {
		return nil, nil, fmt.Errorf("error reading secret file: %w", err)
	}
	if err := viper.Unmarshal(&secret); err != nil {
		return nil, nil, fmt.Errorf("unable to decode into struct: %w", err)
	}

	return &config, &secret, nil
}

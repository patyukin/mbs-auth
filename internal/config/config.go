package config

import (
	"fmt"

	configLoader "github.com/patyukin/mbs-pkg/pkg/config"
)

type Config struct {
	MinLogLevel string `validate:"required,oneof=debug info warn error" yaml:"min_log_level"`
	JwtSecret   string `validate:"required"                             yaml:"jwt_secret"`
	HttpServer  struct {
		Port int `validate:"required,numeric" yaml:"port"`
	} `yaml:"http_server" validate:"required"`
	GRPCServer struct {
		Port              int `validate:"required,numeric" yaml:"port"`
		MaxConnectionIdle int `yaml:"max_connection_idle"`
		Timeout           int `yaml:"timeout"`
		MaxConnectionAge  int `yaml:"max_connection_age"`
	} `yaml:"grpc_server" validate:"required"`
	PostgreSQLDSN string `validate:"required" yaml:"postgresql_dsn"`
	RedisDSN      string `validate:"required" yaml:"redis_dsn"`
	RabbitMQUrl   string `validate:"required" yaml:"rabbitmq_url"`
	Kafka         struct {
		Brokers       []string `validate:"required" yaml:"brokers"`
		ConsumerGroup string   `validate:"required" yaml:"consumer_group"`
		Topics        []string `validate:"required" yaml:"topics"`
	} `yaml:"kafka" validate:"required"`
	TelegramBotName string `validate:"required" yaml:"telegram_bot_name"`
	TracerHost      string `validate:"required" yaml:"tracer_host"`
}

func LoadConfig() (*Config, error) {
	var config Config
	err := configLoader.LoadConfig(&config)
	if err != nil {
		return nil, fmt.Errorf("error loading config: %w", err)
	}

	return &config, nil
}

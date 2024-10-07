package config

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

type Config struct {
	MinLogLevel string `yaml:"min_log_level" validate:"required,oneof=debug info warn error"`
	JwtSecret   string `yaml:"jwt_secret" validate:"required"`
	HttpServer  struct {
		Port int `yaml:"port" validate:"required,numeric"`
	} `yaml:"http_server" validate:"required"`
	SwaggerServer struct {
		Port int `yaml:"port" validate:"required,numeric"`
	} `yaml:"swagger_server" validate:"required"`
	GRPCServer struct {
		Port int `yaml:"port" validate:"required,numeric"`
	} `yaml:"grpc_server" validate:"required"`
	PostgreSQL struct {
		Host     string `yaml:"host" validate:"required"`
		Port     int    `yaml:"port" validate:"required,numeric"`
		User     string `yaml:"user" validate:"required"`
		Password string `yaml:"password" validate:"required"`
		Name     string `yaml:"name" validate:"required"`
	} `yaml:"postgresql"`
	Redis struct {
		Host string `yaml:"host" validate:"required"`
		Port int    `yaml:"port" validate:"required,numeric"`
	} `yaml:"redis"`
	RabbitMQ struct {
		URL        string `yaml:"url" validate:"required"`
		QueueName  string `yaml:"queue_name" validate:"required"`
		Durable    bool   `yaml:"durable" validate:"required,bool"`
		AutoDelete bool   `yaml:"auto_delete" validate:"required,bool"`
		Exclusive  bool   `yaml:"exclusive" validate:"required,bool"`
		NoWait     bool   `yaml:"no_wait" validate:"required,bool"`
	} `yaml:"rabbitmq"`
	TelegramToken string `yaml:"telegram_token" validate:"required"`
}

func LoadConfig() (*Config, error) {
	yamlConfigFilePath := os.Getenv("YAML_CONFIG_FILE_PATH")
	if yamlConfigFilePath == "" {
		return nil, fmt.Errorf("yaml config file path is not set")
	}

	f, err := os.Open(yamlConfigFilePath)
	if err != nil {
		return nil, fmt.Errorf("unable to open config file: %w", err)
	}

	defer func(f *os.File) {
		if err = f.Close(); err != nil {
			log.Printf("unable to close config file: %v", err)
		}
	}(f)

	var config Config
	decoder := yaml.NewDecoder(f)
	if err = decoder.Decode(&config); err != nil {
		return nil, fmt.Errorf("unable to decode config file: %w", err)
	}

	validate := validator.New()
	if err = validate.Struct(&config); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return &config, nil
}

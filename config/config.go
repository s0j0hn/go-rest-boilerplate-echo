package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

// Config holds all application configuration
type Config struct {
	App       AppConfig      `yaml:"app" validate:"required"`
	Server    ServerConfig   `yaml:"server" validate:"required"`
	Database  DatabaseConfig `yaml:"database" validate:"required"`
	RabbitMQ  RabbitMQConfig `yaml:"rabbitmq" validate:"required"`
	WebSocket WebSocketConfig `yaml:"websocket"`
}

// AppConfig holds application-level configuration
type AppConfig struct {
	Environment string `yaml:"environment" validate:"required,oneof=development staging production"`
	LogLevel    string `yaml:"log_level" validate:"required,oneof=debug info warn error"`
}

// ServerConfig holds HTTP server configuration
type ServerConfig struct {
	Address      string        `yaml:"address" validate:"required"`
	Port         int           `yaml:"port" validate:"required,min=1024,max=65535"`
	ReadTimeout  time.Duration `yaml:"read_timeout" validate:"required"`
	WriteTimeout time.Duration `yaml:"write_timeout" validate:"required"`
	IdleTimeout  time.Duration `yaml:"idle_timeout" validate:"required"`
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Host     string `yaml:"host" validate:"required"`
	Port     int    `yaml:"port" validate:"required,min=1,max=65535"`
	User     string `yaml:"user" validate:"required"`
	Password string `yaml:"password" validate:"required"`
	Name     string `yaml:"name" validate:"required"`
	SSLMode  string `yaml:"ssl_mode" validate:"required,oneof=disable require"`
}

// RabbitMQConfig holds RabbitMQ configuration
type RabbitMQConfig struct {
	Host        string `yaml:"host" validate:"required"`
	Port        int    `yaml:"port" validate:"required,min=1,max=65535"`
	User        string `yaml:"user" validate:"required"`
	Password    string `yaml:"password" validate:"required"`
	PushQueue   string `yaml:"push_queue" validate:"required"`
	ListenQueue string `yaml:"listen_queue" validate:"required"`
}

// WebSocketConfig holds WebSocket configuration
type WebSocketConfig struct {
	Address string `yaml:"address"`
}

var globalConfig *Config

// Load loads and validates the configuration from config.yaml
func Load() (*Config, error) {
	if globalConfig != nil {
		return globalConfig, nil
	}

	vp := viper.New()
	vp.SetConfigName("config")
	vp.SetConfigType("yaml")
	vp.AddConfigPath(".")
	vp.AddConfigPath("./config")

	// Set defaults
	vp.SetDefault("app.environment", "development")
	vp.SetDefault("app.log_level", "info")
	vp.SetDefault("server.read_timeout", "30s")
	vp.SetDefault("server.write_timeout", "30s")
	vp.SetDefault("server.idle_timeout", "120s")
	vp.SetDefault("database.port", 5432)
	vp.SetDefault("database.ssl_mode", "disable")
	vp.SetDefault("rabbitmq.port", 5672)

	if err := vp.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := vp.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// Validate configuration
	validate := validator.New()
	if err := validate.Struct(&cfg); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	globalConfig = &cfg
	return globalConfig, nil
}

// Get returns the global configuration
func Get() *Config {
	if globalConfig == nil {
		panic("configuration not loaded, call Load() first")
	}
	return globalConfig
}

// IsProd returns true if running in production environment
func (c *Config) IsProd() bool {
	return c.App.Environment == "production"
}

// GetAddress returns the full server address
func (c *Config) GetAddress() string {
	if c.Server.Address != "" {
		return c.Server.Address
	}
	return fmt.Sprintf(":%d", c.Server.Port)
}

// GetPort returns the server port as string
func (c *Config) GetPort() string {
	if c.Server.Address != "" && strings.Contains(c.Server.Address, ":") {
		return strings.Split(c.Server.Address, ":")[1]
	}
	return fmt.Sprintf("%d", c.Server.Port)
}

// GetDatabaseDSN returns the PostgreSQL connection string
func (c *Config) GetDatabaseDSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Database.Host,
		c.Database.Port,
		c.Database.User,
		c.Database.Password,
		c.Database.Name,
		c.Database.SSLMode,
	)
}

// GetRabbitMQDSN returns the RabbitMQ connection string
func (c *Config) GetRabbitMQDSN() string {
	return fmt.Sprintf(
		"amqp://%s:%s@%s:%d/",
		c.RabbitMQ.User,
		c.RabbitMQ.Password,
		c.RabbitMQ.Host,
		c.RabbitMQ.Port,
	)
}

// Backward compatibility functions (deprecated)
func IsProd() bool {
	return Get().IsProd()
}

func GetPort() string {
	return Get().GetPort()
}

func GetAddress() string {
	return Get().GetAddress()
}

func GetWebSocketAddress() string {
	return Get().WebSocket.Address
}

func GetDatabaseAccess() string {
	return Get().GetDatabaseDSN()
}

func GetRabbitMQAccess() string {
	return Get().GetRabbitMQDSN()
}

func GetAMQPPushQueue() string {
	return Get().RabbitMQ.PushQueue
}

func GetAMQPQListenQueue() string {
	return Get().RabbitMQ.ListenQueue
}

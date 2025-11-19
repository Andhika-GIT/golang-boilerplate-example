// config.go placeholder
// internal/config/config.go
package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// Config holds all application configuration
type Config struct {
	App      AppConfig
	Server   ServerConfig
	Database DatabaseConfig
	Redis    RedisConfig
	RabbitMQ RabbitMQConfig
}

// AppConfig holds application-level configuration
type AppConfig struct {
	Name        string
	Version     string
	Environment string // development, staging, production
	Debug       bool
}

// ServerConfig holds HTTP server configuration
type ServerConfig struct {
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Type            string // postgres, mysql, mongodb
	Host            string
	Port            int
	User            string
	Password        string
	Name            string
	SSLMode         string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

// RedisConfig holds Redis configuration
type RedisConfig struct {
	Enabled  bool
	Host     string
	Port     int
	Password string
	DB       int
	TTL      time.Duration
}

// RabbitMQConfig holds RabbitMQ configuration
type RabbitMQConfig struct {
	Enabled  bool
	Host     string
	Port     int
	User     string
	Password string
	VHost    string
}

// Load loads configuration from environment variables
func Load() *Config {
	// Load .env file if exists (ignore error in production)
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️  No .env file found, using system environment variables")
	}

	cfg := &Config{
		App:      loadAppConfig(),
		Server:   loadServerConfig(),
		Database: loadDatabaseConfig(),
		Redis:    loadRedisConfig(),
		RabbitMQ: loadRabbitMQConfig(),
	}

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		log.Fatalf("❌ Configuration validation failed: %v", err)
	}

	log.Printf("✅ Configuration loaded for environment: %s", cfg.App.Environment)
	return cfg
}

// loadAppConfig loads application configuration
func loadAppConfig() AppConfig {
	return AppConfig{
		Name:        getEnv("APP_NAME", "user-service"),
		Version:     getEnv("APP_VERSION", "1.0.0"),
		Environment: getEnv("APP_ENV", "development"),
		Debug:       getEnvAsBool("APP_DEBUG", true),
	}
}

// loadServerConfig loads server configuration
func loadServerConfig() ServerConfig {
	return ServerConfig{
		Port:         getEnv("SERVER_PORT", "8080"),
		ReadTimeout:  getEnvAsDuration("SERVER_READ_TIMEOUT", 10*time.Second),
		WriteTimeout: getEnvAsDuration("SERVER_WRITE_TIMEOUT", 10*time.Second),
		IdleTimeout:  getEnvAsDuration("SERVER_IDLE_TIMEOUT", 60*time.Second),
	}
}

// loadDatabaseConfig loads database configuration
func loadDatabaseConfig() DatabaseConfig {
	return DatabaseConfig{
		Type:            getEnv("DB_TYPE", "postgres"),
		Host:            getEnv("DB_HOST", "localhost"),
		Port:            getEnvAsInt("DB_PORT", 5432),
		User:            getEnv("DB_USER", "postgres"),
		Password:        getEnv("DB_PASSWORD", "password"),
		Name:            getEnv("DB_NAME", "userdb"),
		SSLMode:         getEnv("DB_SSL_MODE", "disable"),
		MaxOpenConns:    getEnvAsInt("DB_MAX_OPEN_CONNS", 25),
		MaxIdleConns:    getEnvAsInt("DB_MAX_IDLE_CONNS", 5),
		ConnMaxLifetime: getEnvAsDuration("DB_CONN_MAX_LIFETIME", 5*time.Minute),
	}
}

// loadRedisConfig loads Redis configuration
func loadRedisConfig() RedisConfig {
	return RedisConfig{
		Enabled:  getEnvAsBool("REDIS_ENABLED", false),
		Host:     getEnv("REDIS_HOST", "localhost"),
		Port:     getEnvAsInt("REDIS_PORT", 6379),
		Password: getEnv("REDIS_PASSWORD", ""),
		DB:       getEnvAsInt("REDIS_DB", 0),
		TTL:      getEnvAsDuration("REDIS_TTL", 5*time.Minute),
	}
}

// loadRabbitMQConfig loads RabbitMQ configuration
func loadRabbitMQConfig() RabbitMQConfig {
	return RabbitMQConfig{
		Enabled:  getEnvAsBool("RABBITMQ_ENABLED", false),
		Host:     getEnv("RABBITMQ_HOST", "localhost"),
		Port:     getEnvAsInt("RABBITMQ_PORT", 5672),
		User:     getEnv("RABBITMQ_USER", "guest"),
		Password: getEnv("RABBITMQ_PASSWORD", "guest"),
		VHost:    getEnv("RABBITMQ_VHOST", "/"),
	}
}

// Validate validates the configuration
func (c *Config) Validate() error {
	// Validate App
	if c.App.Name == "" {
		return fmt.Errorf("APP_NAME is required")
	}

	// Validate Server
	if c.Server.Port == "" {
		return fmt.Errorf("SERVER_PORT is required")
	}

	// Validate Database
	if c.Database.Type == "" {
		return fmt.Errorf("DB_TYPE is required")
	}
	if c.Database.Host == "" {
		return fmt.Errorf("DB_HOST is required")
	}
	if c.Database.User == "" {
		return fmt.Errorf("DB_USER is required")
	}
	if c.Database.Name == "" {
		return fmt.Errorf("DB_NAME is required")
	}

	return nil
}

// GetDatabaseDSN returns database connection string
func (c *Config) GetDatabaseDSN() string {
	switch c.Database.Type {
	case "postgres":
		return fmt.Sprintf(
			"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
			c.Database.Host,
			c.Database.Port,
			c.Database.User,
			c.Database.Password,
			c.Database.Name,
			c.Database.SSLMode,
		)
	case "mysql":
		return fmt.Sprintf(
			"%s:%s@tcp(%s:%d)/%s?parseTime=true",
			c.Database.User,
			c.Database.Password,
			c.Database.Host,
			c.Database.Port,
			c.Database.Name,
		)
	case "mongodb":
		return fmt.Sprintf(
			"mongodb://%s:%s@%s:%d/%s",
			c.Database.User,
			c.Database.Password,
			c.Database.Host,
			c.Database.Port,
			c.Database.Name,
		)
	default:
		return ""
	}
}

// GetRabbitMQURL returns RabbitMQ connection URL
func (c *Config) GetRabbitMQURL() string {
	return fmt.Sprintf(
		"amqp://%s:%s@%s:%d%s",
		c.RabbitMQ.User,
		c.RabbitMQ.Password,
		c.RabbitMQ.Host,
		c.RabbitMQ.Port,
		c.RabbitMQ.VHost,
	)
}

// GetRedisAddr returns Redis connection address
func (c *Config) GetRedisAddr() string {
	return fmt.Sprintf("%s:%d", c.Redis.Host, c.Redis.Port)
}

// IsProduction returns true if running in production
func (c *Config) IsProduction() bool {
	return c.App.Environment == "production"
}

// IsDevelopment returns true if running in development
func (c *Config) IsDevelopment() bool {
	return c.App.Environment == "development"
}

// IsStaging returns true if running in staging
func (c *Config) IsStaging() bool {
	return c.App.Environment == "staging"
}

// Print prints configuration (with sensitive data masked)
func (c *Config) Print() {
	log.Println("📋 Configuration:")
	log.Printf("  App Name:        %s", c.App.Name)
	log.Printf("  Version:         %s", c.App.Version)
	log.Printf("  Environment:     %s", c.App.Environment)
	log.Printf("  Debug:           %v", c.App.Debug)
	log.Printf("  Server Port:     %s", c.Server.Port)
	log.Printf("  Database Type:   %s", c.Database.Type)
	log.Printf("  Database Host:   %s:%d", c.Database.Host, c.Database.Port)
	log.Printf("  Database Name:   %s", c.Database.Name)
	log.Printf("  Redis Enabled:   %v", c.Redis.Enabled)
	log.Printf("  RabbitMQ Enabled:%v", c.RabbitMQ.Enabled)
}

// ============================================
// Helper Functions
// ============================================

// getEnv gets environment variable with default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvAsInt gets environment variable as integer
func getEnvAsInt(key string, defaultValue int) int {
	valueStr := os.Getenv(key)
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}

// getEnvAsBool gets environment variable as boolean
func getEnvAsBool(key string, defaultValue bool) bool {
	valueStr := os.Getenv(key)
	if value, err := strconv.ParseBool(valueStr); err == nil {
		return value
	}
	return defaultValue
}

// getEnvAsDuration gets environment variable as duration
func getEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
	valueStr := os.Getenv(key)
	if value, err := time.ParseDuration(valueStr); err == nil {
		return value
	}
	return defaultValue
}

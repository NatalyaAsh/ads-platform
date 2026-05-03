package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Server        ServerConfig
	Database      DatabaseConfig
	MongoDB       MongoDBConfig
	Elasticsearch ElasticsearchConfig
	RabbitMQ      RabbitMQConfig
	Log           LogConfig
	Service       ServiceConfig
}

type ServerConfig struct {
	Port           string
	ReadTimeout    time.Duration
	WriteTimeout   time.Duration
	MaxConnections int
}

type DatabaseConfig struct {
	Host            string
	Port            string
	User            string
	Password        string
	DBName          string
	SSLMode         string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

type MongoDBConfig struct {
	URI      string
	Database string
}

type ElasticsearchConfig struct {
	Addresses []string
	Username  string
	Password  string
}

type RabbitMQConfig struct {
	URL string
}

type LogConfig struct {
	Level      string
	Format     string
	OutputPath string
}

type ServiceConfig struct {
	Name        string
	Environment string
	Version     string
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	cfg := &Config{}

	cfg.Server = ServerConfig{
		Port:           getEnv("GRPC_PORT", "50052"),
		ReadTimeout:    getEnvDuration("GRPC_READ_TIMEOUT", 10*time.Second),
		WriteTimeout:   getEnvDuration("GRPC_WRITE_TIMEOUT", 10*time.Second),
		MaxConnections: getEnvInt("GRPC_MAX_CONNECTIONS", 100),
	}

	cfg.Database = DatabaseConfig{
		Host:            getEnv("DB_HOST", "localhost"),
		Port:            getEnv("DB_PORT", "5432"),
		User:            getEnv("DB_USER", "postgres"),
		Password:        getEnv("DB_PASSWORD", "postgres"),
		DBName:          getEnv("DB_NAME", "ad_db"),
		SSLMode:         getEnv("DB_SSL_MODE", "disable"),
		MaxOpenConns:    getEnvInt("DB_MAX_OPEN_CONNS", 25),
		MaxIdleConns:    getEnvInt("DB_MAX_IDLE_CONNS", 25),
		ConnMaxLifetime: getEnvDuration("DB_CONN_MAX_LIFETIME", 5*time.Minute),
	}

	cfg.MongoDB = MongoDBConfig{
		URI:      getEnv("MONGODB_URI", "mongodb://localhost:27017"),
		Database: getEnv("MONGODB_DATABASE", "ad_media"),
	}

	cfg.Elasticsearch = ElasticsearchConfig{
		Addresses: []string{getEnv("ELASTICSEARCH_URL", "http://localhost:9200")},
		Username:  getEnv("ELASTICSEARCH_USERNAME", ""),
		Password:  getEnv("ELASTICSEARCH_PASSWORD", ""),
	}

	cfg.RabbitMQ = RabbitMQConfig{
		URL: getEnv("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/"),
	}

	cfg.Log = LogConfig{
		Level:      getEnv("LOG_LEVEL", "info"),
		Format:     getEnv("LOG_FORMAT", "json"),
		OutputPath: getEnv("LOG_OUTPUT_PATH", "stdout"),
	}

	cfg.Service = ServiceConfig{
		Name:        getEnv("SERVICE_NAME", "ad-search-service"),
		Environment: getEnv("ENV", "dev"),
		Version:     getEnv("SERVICE_VERSION", "1.0.0"),
	}

	return cfg, nil
}

func (c *Config) GetDSN() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		c.Database.User,
		c.Database.Password,
		c.Database.Host,
		c.Database.Port,
		c.Database.DBName,
		c.Database.SSLMode,
	)
}

// Вспомогательные функции
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvDuration(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}

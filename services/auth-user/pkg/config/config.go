package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

// Config - главная структура конфигурации
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
	//Redis    RedisConfig // опционально, для refresh токенов
	Log LogConfig
	//SMTP     SMTPConfig // для email подтверждения
	Service ServiceConfig
}

// ServerConfig - настройки gRPC сервера
type ServerConfig struct {
	Port           string        // gRPC порт (50051)
	ReadTimeout    time.Duration // таймаут чтения
	WriteTimeout   time.Duration // таймаут записи
	MaxConnections int           // макс. соединений
}

// DatabaseConfig - настройки PostgreSQL
type DatabaseConfig struct {
	Host            string
	Port            string
	User            string
	Password        string
	DBName          string
	SSLMode         string
	MaxOpenConns    int           // макс. открытых соединений
	MaxIdleConns    int           // макс. idle соединений
	ConnMaxLifetime time.Duration // время жизни соединения
}

// JWTConfig - настройки JWT токенов
type JWTConfig struct {
	AccessSecret  string        // секрет для access токена
	RefreshSecret string        // секрет для refresh токена
	AccessTTL     time.Duration // время жизни access токена (24h)
	RefreshTTL    time.Duration // время жизни refresh токена (720h = 30 дней)
	Issuer        string        // издатель токена (например, "auth-service")
}

// // RedisConfig - для хранения refresh токенов (опционально)
// type RedisConfig struct {
// 	Host     string
// 	Port     string
// 	Password string
// 	DB       int
// 	PoolSize int
// }

// LogConfig - настройки логирования
type LogConfig struct {
	Level      string // debug, info, warn, error
	Format     string // json или text
	OutputPath string // stdout или файл
}

// // SMTPConfig - для отправки email
// type SMTPConfig struct {
// 	Host      string
// 	Port      int
// 	Username  string
// 	Password  string
// 	FromEmail string
// 	FromName  string
// }

// ServiceConfig - настройки самого сервиса
type ServiceConfig struct {
	Name        string // "auth-user-service"
	Environment string // dev, stage, prod
	Version     string // версия сервиса
	EnableTLS   bool   // включить TLS для gRPC
	CertFile    string // путь к сертификату
	KeyFile     string // путь к ключу
}

// Load загружает конфигурацию из .env файла и переменных окружения
func Load() (*Config, error) {
	// Пытаемся загрузить .env файл (игнорируем ошибку если файла нет)
	_ = godotenv.Load()

	cfg := &Config{}

	// Загружаем серверную конфигурацию
	cfg.Server = ServerConfig{
		Port:           getEnv("GRPC_PORT", "50051"),
		ReadTimeout:    getEnvDuration("GRPC_READ_TIMEOUT", 10*time.Second),
		WriteTimeout:   getEnvDuration("GRPC_WRITE_TIMEOUT", 10*time.Second),
		MaxConnections: getEnvInt("GRPC_MAX_CONNECTIONS", 100),
	}

	// Загружаем конфигурацию БД
	cfg.Database = DatabaseConfig{
		Host:            getEnv("DB_HOST", "localhost"),
		Port:            getEnv("DB_PORT", "5432"),
		User:            getEnv("DB_USER", "postgres"),
		Password:        getEnv("DB_PASSWORD", "postgres"),
		DBName:          getEnv("DB_NAME", "auth_user_db"),
		SSLMode:         getEnv("DB_SSL_MODE", "disable"),
		MaxOpenConns:    getEnvInt("DB_MAX_OPEN_CONNS", 25),
		MaxIdleConns:    getEnvInt("DB_MAX_IDLE_CONNS", 25),
		ConnMaxLifetime: getEnvDuration("DB_CONN_MAX_LIFETIME", 5*time.Minute),
	}

	// Загружаем JWT конфигурацию
	cfg.JWT = JWTConfig{
		AccessSecret:  getEnv("JWT_ACCESS_SECRET", "default-access-secret-change-me"),
		RefreshSecret: getEnv("JWT_REFRESH_SECRET", "default-refresh-secret-change-me"),
		AccessTTL:     getEnvDuration("JWT_ACCESS_TTL", 24*time.Hour),
		RefreshTTL:    getEnvDuration("JWT_REFRESH_TTL", 720*time.Hour), // 30 дней
		Issuer:        getEnv("JWT_ISSUER", "auth-service"),
	}

	// Загружаем конфигурацию логов
	cfg.Log = LogConfig{
		Level:      getEnv("LOG_LEVEL", "info"),
		Format:     getEnv("LOG_FORMAT", "json"),
		OutputPath: getEnv("LOG_OUTPUT_PATH", "stdout"),
	}

	// Загружаем конфигурацию сервиса
	cfg.Service = ServiceConfig{
		Name:        getEnv("SERVICE_NAME", "auth-user-service"),
		Environment: getEnv("ENVIRONMENT", "dev"),
		Version:     getEnv("SERVICE_VERSION", "1.0.0"),
		EnableTLS:   getEnvBool("ENABLE_TLS", false),
		CertFile:    getEnv("TLS_CERT_FILE", ""),
		KeyFile:     getEnv("TLS_KEY_FILE", ""),
	}

	// Валидируем обязательные поля
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return cfg, nil
}

// Validate проверяет обязательные поля конфигурации
func (c *Config) Validate() error {
	// Проверяем JWT секреты в production
	if c.Service.Environment == "prod" {
		if c.JWT.AccessSecret == "default-access-secret-change-me" {
			return fmt.Errorf("JWT_ACCESS_SECRET must be changed in production")
		}
		if c.JWT.RefreshSecret == "default-refresh-secret-change-me" {
			return fmt.Errorf("JWT_REFRESH_SECRET must be changed in production")
		}
	}

	// Проверяем порт
	if c.Server.Port == "" {
		return fmt.Errorf("GRPC_PORT is required")
	}

	// Проверяем подключение к БД
	if c.Database.Host == "" {
		return fmt.Errorf("DB_HOST is required")
	}
	if c.Database.User == "" {
		return fmt.Errorf("DB_USER is required")
	}
	if c.Database.DBName == "" {
		return fmt.Errorf("DB_NAME is required")
	}

	return nil
}

// Вспомогательные функции для чтения переменных окружения

// getEnv возвращает значение переменной окружения или defaultValue
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvInt возвращает int значение переменной окружения или defaultValue
func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// getEnvBool возвращает bool значение переменной окружения или defaultValue
func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		value = strings.ToLower(value)
		if value == "true" || value == "1" || value == "yes" {
			return true
		}
		if value == "false" || value == "0" || value == "no" {
			return false
		}
	}
	return defaultValue
}

// getEnvDuration возвращает time.Duration значение переменной окружения или defaultValue
func getEnvDuration(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
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

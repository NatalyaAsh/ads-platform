package repository

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/NatalyaAsh/ads-platform/services/auth-user/pkg/config"
)

// NewPostgresDB создает подключение к PostgreSQL с настройками пула
func NewPostgresDB(cfg *config.DatabaseConfig) (*gorm.DB, error) {
	// Формируем DSN
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=UTC",
		cfg.Host,
		cfg.User,
		cfg.Password,
		cfg.DBName,
		cfg.Port,
		cfg.SSLMode,
	)

	// Настройки логгера GORM
	gormLogger := logger.Default.LogMode(logger.Info)
	if cfg.SSLMode == "disable" {
		// В dev режиме можно включить подробное логирование
		gormLogger = logger.Default.LogMode(logger.Info)
	}

	// Открываем соединение
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Получаем原生 sql.DB для настройки пула
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get sql.DB: %w", err)
	}

	// Настройка пула соединений
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(cfg.ConnMaxLifetime)

	log.Printf("Database connected: %s:%s/%s", cfg.Host, cfg.Port, cfg.DBName)
	log.Printf("Connection pool: max_open=%d, max_idle=%d, max_lifetime=%v",
		cfg.MaxOpenConns, cfg.MaxIdleConns, cfg.ConnMaxLifetime)

	return db, nil
}

// CloseDB закрывает соединение с БД
func CloseDB(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

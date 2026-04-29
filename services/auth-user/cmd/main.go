package main

import (
	"log"

	"github.com/NatalyaAsh/ads-platform/services/auth-user/internal/repository"
	"github.com/NatalyaAsh/ads-platform/services/auth-user/pkg/config"
	"github.com/NatalyaAsh/ads-platform/services/auth-user/pkg/db"
)

func main() {
	// 1. Загружаем конфигурацию
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	log.Printf("Starting %s v%s", cfg.Service.Name, cfg.Service.Version)
	log.Printf("Environment: %s", cfg.Service.Environment)

	// 2. Накатываем миграции
	dsn := cfg.GetDSN()
	migrationsPath := "file://migrations"

	log.Println("Running database migrations...")
	if err := db.RunMigrations(dsn, migrationsPath); err != nil {
		log.Fatalf("Migrations failed: %v", err)
	}

	// 3. Подключаемся к БД через GORM
	log.Println("Connecting to database...")
	gormDB, err := repository.NewPostgresDB(&cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer func() {
		if err := repository.CloseDB(gormDB); err != nil {
			log.Printf("Error closing database: %v", err)
		}
	}()

	log.Println("Database connected successfully")

	// 4. Здесь позже будут:
	// - Создание репозиториев
	// - Создание сервисного слоя
	// - Создание gRPC хендлеров
	// - Запуск gRPC сервера

	log.Println("Auth-user service initialized successfully")

	// Временная заглушка, чтобы сервис не завершался
	select {}
}

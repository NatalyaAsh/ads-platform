package main

import (
	"log"

	"github.com/NatalyaAsh/ads-platform/services/auth-user/internal/repository"
	"github.com/NatalyaAsh/ads-platform/services/auth-user/internal/server"
	"github.com/NatalyaAsh/ads-platform/services/auth-user/internal/service"
	"github.com/NatalyaAsh/ads-platform/services/auth-user/pkg/config"
	"github.com/NatalyaAsh/ads-platform/services/auth-user/pkg/db"
	"github.com/NatalyaAsh/ads-platform/services/auth-user/pkg/jwt"
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

	// 4. Создаём репозитории
	userRepo := repository.NewUserRepository(gormDB)
	refreshRepo := repository.NewRefreshTokenRepository(gormDB)

	// 5. Создаём JWT менеджер
	jwtManager := jwt.NewJWTManager(&cfg.JWT)

	// 6. Создаём сервис
	authService := service.NewAuthService(userRepo, refreshRepo, jwtManager)

	log.Println("Auth-user service initialized successfully")

	// 7. Запускаем gRPC сервер
	log.Printf("Starting gRPC server on port %s", cfg.Server.Port)
	if err := server.RunGRPCServer(cfg.Server.Port, authService); err != nil {
		log.Fatalf("Failed to start gRPC server: %v", err)
	}

}

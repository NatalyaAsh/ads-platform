package main

import (
	"context"
	"log"

	"github.com/NatalyaAsh/ads-platform/services/ad-search/internal/repository"
	"github.com/NatalyaAsh/ads-platform/services/ad-search/pkg/config"
	"github.com/NatalyaAsh/ads-platform/services/ad-search/pkg/db"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	log.Printf("Starting %s v%s", cfg.Service.Name, cfg.Service.Version)

	// Миграции
	dsn := cfg.GetDSN()
	if err := db.RunMigrations(dsn, "file://migrations"); err != nil {
		log.Fatalf("Migrations failed: %v", err)
	}

	// Подключение к БД
	gormDB, err := repository.NewPostgresDB(&cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer repository.CloseDB(gormDB)

	log.Println("Ad-search service initialized successfully")

	// Подключение к MongoDB
	mongoClient, err := repository.NewMongoClient(cfg.MongoDB.URI)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer mongoClient.Disconnect(context.Background())

	mongoDB := mongoClient.Database(cfg.MongoDB.Database)
	mediaRepo := repository.NewMediaRepository(mongoDB)
	log.Printf("Media repository initialized")
	_ = mediaRepo // временно игнорируем

	// Пока заглушка
	select {}
}

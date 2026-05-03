package main

import (
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

	// Пока заглушка
	select {}
}

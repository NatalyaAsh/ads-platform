package db

import (
	"fmt"
	"log"

	migrate "github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// RunMigrations накатывает все миграции из папки migrations
// dsn - строка подключения к БД в формате PostgreSQL
// migrationsPath - путь к папке с миграциями (например, "file://migrations")
func RunMigrations(dsn string, migrationsPath string) error {
	// Создаем экземпляр migrate
	m, err := migrate.New(
		migrationsPath,
		dsn,
	)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}
	defer m.Close()

	// Накатываем миграции
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	if err == migrate.ErrNoChange {
		log.Println("No new migrations to apply")
	} else {
		log.Println("Migrations applied successfully")
	}

	return nil
}

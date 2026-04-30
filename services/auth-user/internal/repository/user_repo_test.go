package repository

import (
	"testing"

	"github.com/NatalyaAsh/ads-platform/services/auth-user/internal/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}

	if err := db.AutoMigrate(&model.User{}, &model.RefreshToken{}); err != nil {
		t.Fatalf("Failed to migrate: %v", err)
	}

	return db
}

func cleanupDB(db *gorm.DB) {
	// SQLite не поддерживает TRUNCATE
	db.Exec("DELETE FROM refresh_tokens")
	db.Exec("DELETE FROM users")
}

func TestUserRepository_Create(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupDB(db)

	repo := NewUserRepository(db)

	tests := []struct {
		name    string
		user    *model.User
		wantErr bool
	}{
		{
			name: "valid user",
			user: &model.User{
				Email:        "test@example.com",
				PasswordHash: "hashedpass",
				Role:         model.RoleUser,
			},
			wantErr: false,
		},
		{
			name: "duplicate email",
			user: &model.User{
				Email:        "test@example.com",
				PasswordHash: "hashedpass2",
				Role:         model.RoleUser,
			},
			wantErr: true,
		},
		{
			name: "invalid email",
			user: &model.User{
				Email:        "invalid",
				PasswordHash: "hashedpass",
				Role:         model.RoleUser,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.Create(tt.user)
			if (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUserRepository_FindByEmail(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupDB(db)

	repo := NewUserRepository(db)

	// Создаем тестового пользователя
	user := &model.User{
		Email:        "findme@example.com",
		PasswordHash: "hashed",
		Role:         model.RoleUser,
	}
	if err := repo.Create(user); err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	// Тест: находим существующего
	found, err := repo.FindByEmail("findme@example.com")
	if err != nil {
		t.Errorf("FindByEmail() error = %v", err)
	}
	if found.Email != user.Email {
		t.Errorf("FindByEmail() got = %v, want %v", found.Email, user.Email)
	}

	// Тест: ищем несуществующего
	_, err = repo.FindByEmail("notexists@example.com")
	if err != model.ErrUserNotFound {
		t.Errorf("FindByEmail() error = %v, want %v", err, model.ErrUserNotFound)
	}
}

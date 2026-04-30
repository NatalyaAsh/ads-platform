package repository

import (
	"testing"
	"time"

	"github.com/NatalyaAsh/ads-platform/services/auth-user/internal/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Запуск go test ./internal/repository/... -v

func setupRefreshTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}

	if err := db.AutoMigrate(&model.User{}, &model.RefreshToken{}); err != nil {
		t.Fatalf("Failed to migrate: %v", err)
	}

	return db
}

func cleanupRefreshDB(db *gorm.DB) {
	db.Exec("DELETE FROM refresh_tokens")
	db.Exec("DELETE FROM users")
}

// createTestUser создает тестового пользователя
func createTestUser(t *testing.T, db *gorm.DB) *model.User {
	user := &model.User{
		Email:        "testuser@example.com",
		PasswordHash: "hashedpassword",
		Role:         model.RoleUser,
		IsBlocked:    false,
	}
	if err := db.Create(user).Error; err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}
	return user
}

func TestRefreshTokenRepository_Create(t *testing.T) {
	db := setupRefreshTestDB(t)
	defer cleanupRefreshDB(db)

	user := createTestUser(t, db)
	repo := NewRefreshTokenRepository(db)

	token := &model.RefreshToken{
		UserID:    user.ID,
		Token:     "test-token-123",
		ExpiresAt: time.Now().Add(24 * time.Hour),
		Revoked:   false,
	}

	err := repo.Create(token)
	if err != nil {
		t.Errorf("Create() error = %v", err)
	}

	// Проверяем, что токен сохранился
	var saved model.RefreshToken
	db.Where("token = ?", "test-token-123").First(&saved)
	if saved.ID == 0 {
		t.Error("Token was not saved")
	}
}

func TestRefreshTokenRepository_FindByToken(t *testing.T) {
	db := setupRefreshTestDB(t)
	defer cleanupRefreshDB(db)

	user := createTestUser(t, db)
	repo := NewRefreshTokenRepository(db)

	// Создаем токен
	token := &model.RefreshToken{
		UserID:    user.ID,
		Token:     "findme-token",
		ExpiresAt: time.Now().Add(24 * time.Hour),
		Revoked:   false,
	}
	repo.Create(token)

	// Находим по значению
	found, err := repo.FindByToken("findme-token")
	if err != nil {
		t.Errorf("FindByToken() error = %v", err)
	}
	if found.Token != "findme-token" {
		t.Errorf("FindByToken() got = %v, want %v", found.Token, "findme-token")
	}

	// Ищем несуществующий
	_, err = repo.FindByToken("nonexistent")
	if err != model.ErrTokenNotFound {
		t.Errorf("FindByToken() error = %v, want %v", err, model.ErrTokenNotFound)
	}
}

func TestRefreshTokenRepository_Revoke(t *testing.T) {
	db := setupRefreshTestDB(t)
	defer cleanupRefreshDB(db)

	user := createTestUser(t, db)
	repo := NewRefreshTokenRepository(db)

	// Создаем токен
	token := &model.RefreshToken{
		UserID:    user.ID,
		Token:     "revoke-me",
		ExpiresAt: time.Now().Add(24 * time.Hour),
		Revoked:   false,
	}
	repo.Create(token)

	// Находим токен
	found, _ := repo.FindByToken("revoke-me")

	// Отзываем
	err := repo.Revoke(found.ID)
	if err != nil {
		t.Errorf("Revoke() error = %v", err)
	}

	// Проверяем, что отозван
	revoked, _ := repo.FindByToken("revoke-me")
	if !revoked.Revoked {
		t.Error("Token should be revoked")
	}
}

func TestRefreshTokenRepository_RevokeByTokenValue(t *testing.T) {
	db := setupRefreshTestDB(t)
	defer cleanupRefreshDB(db)

	user := createTestUser(t, db)
	repo := NewRefreshTokenRepository(db)

	// Создаем токен
	token := &model.RefreshToken{
		UserID:    user.ID,
		Token:     "revoke-by-value",
		ExpiresAt: time.Now().Add(24 * time.Hour),
		Revoked:   false,
	}
	repo.Create(token)

	// Отзываем по значению
	err := repo.RevokeByTokenValue("revoke-by-value")
	if err != nil {
		t.Errorf("RevokeByTokenValue() error = %v", err)
	}

	// Проверяем
	revoked, _ := repo.FindByToken("revoke-by-value")
	if !revoked.Revoked {
		t.Error("Token should be revoked")
	}
}

func TestRefreshTokenRepository_RevokeAllByUserID(t *testing.T) {
	db := setupRefreshTestDB(t)
	defer cleanupRefreshDB(db)

	user := createTestUser(t, db)
	repo := NewRefreshTokenRepository(db)

	// Создаем несколько токенов
	tokens := []string{"token1", "token2", "token3"}
	for _, tokenStr := range tokens {
		token := &model.RefreshToken{
			UserID:    user.ID,
			Token:     tokenStr,
			ExpiresAt: time.Now().Add(24 * time.Hour),
			Revoked:   false,
		}
		repo.Create(token)
	}

	// Отзываем все токены пользователя
	err := repo.RevokeAllByUserID(user.ID)
	if err != nil {
		t.Errorf("RevokeAllByUserID() error = %v", err)
	}

	// Проверяем, что все отозваны
	activeTokens, _ := repo.FindActiveByUserID(user.ID)
	if len(activeTokens) != 0 {
		t.Errorf("All tokens should be revoked, but %d active", len(activeTokens))
	}
}

func TestRefreshTokenRepository_FindActiveByUserID(t *testing.T) {
	db := setupRefreshTestDB(t)
	defer cleanupRefreshDB(db)

	user := createTestUser(t, db)
	repo := NewRefreshTokenRepository(db)

	// Создаем активный токен
	activeToken := &model.RefreshToken{
		UserID:    user.ID,
		Token:     "active-token",
		ExpiresAt: time.Now().Add(24 * time.Hour),
		Revoked:   false,
	}
	repo.Create(activeToken)

	// Создаем просроченный токен
	expiredToken := &model.RefreshToken{
		UserID:    user.ID,
		Token:     "expired-token",
		ExpiresAt: time.Now().Add(-1 * time.Hour),
		Revoked:   false,
	}
	repo.Create(expiredToken)

	// Создаем отозванный токен
	revokedToken := &model.RefreshToken{
		UserID:    user.ID,
		Token:     "revoked-token",
		ExpiresAt: time.Now().Add(24 * time.Hour),
		Revoked:   true,
	}
	repo.Create(revokedToken)

	// Находим активные токены
	active, err := repo.FindActiveByUserID(user.ID)
	if err != nil {
		t.Errorf("FindActiveByUserID() error = %v", err)
	}

	// Должен быть только один активный
	if len(active) != 1 {
		t.Errorf("Expected 1 active token, got %d", len(active))
	}
	if active[0].Token != "active-token" {
		t.Errorf("Expected 'active-token', got '%s'", active[0].Token)
	}
}

func TestRefreshTokenRepository_DeleteExpired(t *testing.T) {
	db := setupRefreshTestDB(t)
	defer cleanupRefreshDB(db)

	user := createTestUser(t, db)
	repo := NewRefreshTokenRepository(db)

	// Создаем просроченный токен
	expiredToken := &model.RefreshToken{
		UserID:    user.ID,
		Token:     "expired",
		ExpiresAt: time.Now().Add(-1 * time.Hour),
		Revoked:   false,
	}
	repo.Create(expiredToken)

	// Создаем активный токен
	activeToken := &model.RefreshToken{
		UserID:    user.ID,
		Token:     "active",
		ExpiresAt: time.Now().Add(24 * time.Hour),
		Revoked:   false,
	}
	repo.Create(activeToken)

	// Удаляем просроченные
	err := repo.DeleteExpired()
	if err != nil {
		t.Errorf("DeleteExpired() error = %v", err)
	}

	// Проверяем, что просроченный удален
	_, err = repo.FindByToken("expired")
	if err != model.ErrTokenNotFound {
		t.Error("Expired token should be deleted")
	}

	// Проверяем, что активный остался
	_, err = repo.FindByToken("active")
	if err != nil {
		t.Error("Active token should remain")
	}
}

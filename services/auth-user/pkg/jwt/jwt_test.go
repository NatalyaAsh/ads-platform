package jwt

import (
	"testing"
	"time"

	"github.com/NatalyaAsh/ads-platform/services/auth-user/pkg/config"
)

func setupTestJWTManager() *JWTManager {
	cfg := &config.JWTConfig{
		AccessSecret:  "test-access-secret",
		RefreshSecret: "test-refresh-secret",
		AccessTTL:     1 * time.Hour,
		RefreshTTL:    24 * time.Hour,
		Issuer:        "test-issuer",
	}
	return NewJWTManager(cfg)
}

func TestJWTManager_GenerateAccessToken(t *testing.T) {
	manager := setupTestJWTManager()

	token, err := manager.GenerateAccessToken(1, "user")
	if err != nil {
		t.Errorf("GenerateAccessToken() error = %v", err)
	}
	if token == "" {
		t.Error("GenerateAccessToken() returned empty token")
	}
}

func TestJWTManager_GenerateRefreshToken(t *testing.T) {
	manager := setupTestJWTManager()

	token, err := manager.GenerateRefreshToken(1)
	if err != nil {
		t.Errorf("GenerateRefreshToken() error = %v", err)
	}
	if token == "" {
		t.Error("GenerateRefreshToken() returned empty token")
	}
}

func TestJWTManager_ValidateAccessToken(t *testing.T) {
	manager := setupTestJWTManager()

	// Генерируем токен
	accessToken, _ := manager.GenerateAccessToken(1, "admin")

	// Валидируем
	claims, err := manager.ValidateAccessToken(accessToken)
	if err != nil {
		t.Errorf("ValidateAccessToken() error = %v", err)
	}
	if claims.UserID != 1 {
		t.Errorf("ValidateAccessToken() userID = %v, want %v", claims.UserID, 1)
	}
	if claims.Role != "admin" {
		t.Errorf("ValidateAccessToken() role = %v, want %v", claims.Role, "admin")
	}
}

func TestJWTManager_ValidateRefreshToken(t *testing.T) {
	manager := setupTestJWTManager()

	// Генерируем refresh токен
	refreshToken, _ := manager.GenerateRefreshToken(1)

	// Валидируем
	claims, err := manager.ValidateRefreshToken(refreshToken)
	if err != nil {
		t.Errorf("ValidateRefreshToken() error = %v", err)
	}
	if claims.UserID != 1 {
		t.Errorf("ValidateRefreshToken() userID = %v, want %v", claims.UserID, 1)
	}
}

func TestJWTManager_InvalidToken(t *testing.T) {
	manager := setupTestJWTManager()

	// Неверный токен
	_, err := manager.ValidateAccessToken("invalid-token")
	if err == nil {
		t.Error("ValidateAccessToken() should return error for invalid token")
	}
}

func TestJWTManager_GetUserIDFromAccessToken(t *testing.T) {
	manager := setupTestJWTManager()

	accessToken, _ := manager.GenerateAccessToken(42, "user")

	userID, err := manager.GetUserIDFromAccessToken(accessToken)
	if err != nil {
		t.Errorf("GetUserIDFromAccessToken() error = %v", err)
	}
	if userID != 42 {
		t.Errorf("GetUserIDFromAccessToken() = %v, want %v", userID, 42)
	}
}

func TestJWTManager_RefreshTokens(t *testing.T) {
	manager := setupTestJWTManager()

	// Создаём refresh токен
	oldRefreshToken, _ := manager.GenerateRefreshToken(1)

	// Обновляем пару токенов
	newAccess, newRefresh, err := manager.RefreshTokens(oldRefreshToken)
	if err != nil {
		t.Errorf("RefreshTokens() error = %v", err)
	}
	if newAccess == "" {
		t.Error("RefreshTokens() returned empty access token")
	}
	if newRefresh == "" {
		t.Error("RefreshTokens() returned empty refresh token")
	}
}

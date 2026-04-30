// internal/service/auth_service_test.go
package service

import (
	"testing"

	"github.com/NatalyaAsh/ads-platform/services/auth-user/internal/model"
	"github.com/NatalyaAsh/ads-platform/services/auth-user/internal/repository"
	"github.com/NatalyaAsh/ads-platform/services/auth-user/pkg/config"
	"github.com/NatalyaAsh/ads-platform/services/auth-user/pkg/jwt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// setupTestAuthService создает новый сервис и БД ДЛЯ КАЖДОГО ТЕСТА
func setupTestAuthService(t *testing.T) (*AuthService, *gorm.DB) {
	// Открываем новую in-memory БД для каждого теста
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open db: %v", err)
	}

	// Создаем схему
	err = db.AutoMigrate(&model.User{}, &model.RefreshToken{})
	if err != nil {
		t.Fatalf("Failed to migrate: %v", err)
	}

	userRepo := repository.NewUserRepository(db)
	refreshRepo := repository.NewRefreshTokenRepository(db)

	// Создаем тестовый JWT конфиг
	jwtConfig := &config.JWTConfig{
		AccessSecret:  "test-secret",
		RefreshSecret: "test-refresh-secret",
		AccessTTL:     24 * 60 * 60,  // 24 часа в секундах
		RefreshTTL:    720 * 60 * 60, // 720 часов
		Issuer:        "test",
	}
	jwtManager := jwt.NewJWTManager(jwtConfig)

	authService := NewAuthService(userRepo, refreshRepo, jwtManager)

	return authService, db
}

// Тест регистрации
func TestAuthService_Register(t *testing.T) {
	service, _ := setupTestAuthService(t)

	req := &RegisterRequest{
		Email:    "test@example.com",
		Password: "password123",
	}

	resp, err := service.Register(req)
	if err != nil {
		t.Fatalf("Register() error = %v", err)
	}
	if resp.UserID == 0 {
		t.Error("Register() returned zero user ID")
	}
	if resp.AccessToken == "" {
		t.Error("Register() returned empty access token")
	}
	if resp.RefreshToken == "" {
		t.Error("Register() returned empty refresh token")
	}
}

// Тест дублирования email
func TestAuthService_Register_DuplicateEmail(t *testing.T) {
	service, _ := setupTestAuthService(t)

	req := &RegisterRequest{
		Email:    "duplicate@example.com",
		Password: "password123",
	}

	// Первая регистрация
	_, err := service.Register(req)
	if err != nil {
		t.Fatalf("First register failed: %v", err)
	}

	// Вторая регистрация с тем же email
	_, err = service.Register(req)
	if err != model.ErrUserAlreadyExists {
		t.Errorf("Expected ErrUserAlreadyExists, got %v", err)
	}
}

// Тест логина (теперь полностью изолирован)
func TestAuthService_Login(t *testing.T) {
	service, _ := setupTestAuthService(t) // <-- Новая БД для каждого запуска!

	// 1. Регистрируем пользователя
	registerReq := &RegisterRequest{
		Email:    "login@example.com",
		Password: "password123",
	}
	_, err := service.Register(registerReq)
	if err != nil {
		t.Fatalf("Register failed: %v", err)
	}

	// 2. Пытаемся войти
	loginReq := &LoginRequest{
		Email:    "login@example.com",
		Password: "password123",
	}
	loginResp, err := service.Login(loginReq)
	if err != nil {
		t.Fatalf("Login() error = %v", err) // <-- Используем Fatalf, чтобы не падать с nil
	}

	// 3. Проверки
	if loginResp.UserID == 0 {
		t.Error("Login() returned zero user ID")
	}
	if loginResp.AccessToken == "" {
		t.Error("Login() returned empty access token")
	}
	if loginResp.RefreshToken == "" {
		t.Error("Login() returned empty refresh token")
	}
}

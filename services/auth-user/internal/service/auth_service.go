package service

import (
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/NatalyaAsh/ads-platform/services/auth-user/internal/model"
	"github.com/NatalyaAsh/ads-platform/services/auth-user/internal/repository"
	"github.com/NatalyaAsh/ads-platform/services/auth-user/pkg/jwt"
)

var (
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrUserBlocked        = errors.New("user is blocked")
)

// AuthService - сервис для аутентификации
type AuthService struct {
	userRepo    *repository.UserRepository
	refreshRepo *repository.RefreshTokenRepository
	jwtManager  *jwt.JWTManager
}

// NewAuthService создает новый сервис аутентификации
func NewAuthService(
	userRepo *repository.UserRepository,
	refreshRepo *repository.RefreshTokenRepository,
	jwtManager *jwt.JWTManager,
) *AuthService {
	return &AuthService{
		userRepo:    userRepo,
		refreshRepo: refreshRepo,
		jwtManager:  jwtManager,
	}
}

// RegisterRequest - запрос на регистрацию
type RegisterRequest struct {
	Email    string
	Password string
}

// RegisterResponse - ответ на регистрацию
type RegisterResponse struct {
	UserID       uint
	AccessToken  string
	RefreshToken string
}

// Register регистрирует нового пользователя
func (s *AuthService) Register(req *RegisterRequest) (*RegisterResponse, error) {
	// Проверяем, существует ли пользователь
	exists, err := s.userRepo.ExistsByEmail(req.Email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, model.ErrUserAlreadyExists
	}

	// Хешируем пароль
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	// Создаем пользователя
	user := &model.User{
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		Role:         model.RoleUser,
		IsBlocked:    false,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	// Отзываем старые токены (на всякий случай, хотя пользователь новый)
	if err := s.refreshRepo.DeleteByUserID(user.ID); err != nil {
		return nil, errors.New("failed to clean old tokens")
	}

	// Генерируем токены
	accessToken, err := s.jwtManager.GenerateAccessToken(user.ID, string(user.Role))
	if err != nil {
		return nil, errors.New("failed to generate access token")
	}

	refreshToken, err := s.jwtManager.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, errors.New("failed to generate refresh token")
	}

	// Сохраняем refresh токен в БД
	refreshTokenModel := &model.RefreshToken{
		UserID:    user.ID,
		Token:     refreshToken,
		ExpiresAt: time.Now().Add(720 * time.Hour), // 30 дней
		Revoked:   false,
	}

	if err := s.refreshRepo.Create(refreshTokenModel); err != nil {
		return nil, errors.New("failed to save refresh token")
	}

	return &RegisterResponse{
		UserID:       user.ID,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// LoginRequest - запрос на вход
type LoginRequest struct {
	Email    string
	Password string
}

// LoginResponse - ответ на вход
type LoginResponse struct {
	UserID       uint
	Role         string
	AccessToken  string
	RefreshToken string
}

// Login выполняет вход пользователя
func (s *AuthService) Login(req *LoginRequest) (*LoginResponse, error) {
	// Находим пользователя
	user, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		if errors.Is(err, model.ErrUserNotFound) {
			return nil, errors.New("invalid email or password")
		}
		return nil, err
	}

	// Проверяем, не заблокирован ли пользователь
	if user.IsBlocked {
		return nil, model.ErrUserBlocked
	}

	// Проверяем пароль
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid email or password")
	}

	// Отзываем все старые refresh токены пользователя
	if err := s.refreshRepo.DeleteByUserID(user.ID); err != nil {
		return nil, errors.New("failed to clean old tokens")
	}

	// Генерируем токены
	accessToken, err := s.jwtManager.GenerateAccessToken(user.ID, string(user.Role))
	if err != nil {
		return nil, errors.New("failed to generate access token")
	}

	refreshToken, err := s.jwtManager.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, errors.New("failed to generate refresh token")
	}

	// Сохраняем refresh токен в БД
	refreshTokenModel := &model.RefreshToken{
		UserID:    user.ID,
		Token:     refreshToken,
		ExpiresAt: time.Now().Add(720 * time.Hour),
		Revoked:   false,
	}

	if err := s.refreshRepo.Create(refreshTokenModel); err != nil {
		return nil, errors.New("failed to save refresh token")
	}

	return &LoginResponse{
		UserID:       user.ID,
		Role:         string(user.Role),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// ValidateToken валидирует access токен
func (s *AuthService) ValidateToken(tokenString string) (*jwt.Claims, error) {
	return s.jwtManager.ValidateAccessToken(tokenString)
}

// RefreshTokens обновляет пару токенов
func (s *AuthService) RefreshTokens(refreshToken string) (string, string, error) {
	// Валидируем refresh токен
	claims, err := s.jwtManager.ValidateRefreshToken(refreshToken)
	if err != nil {
		return "", "", errors.New("invalid refresh token")
	}

	// Проверяем, существует ли токен в БД и не отозван ли он
	savedToken, err := s.refreshRepo.FindByToken(refreshToken)
	if err != nil {
		return "", "", errors.New("refresh token not found")
	}

	if savedToken.Revoked {
		return "", "", errors.New("refresh token revoked")
	}

	if savedToken.IsExpired() {
		return "", "", errors.New("refresh token expired")
	}

	// Находим пользователя
	user, err := s.userRepo.FindByID(claims.UserID)
	if err != nil {
		return "", "", err
	}

	// Отзываем старый токен
	if err := s.refreshRepo.Revoke(savedToken.ID); err != nil {
		return "", "", err
	}

	// Генерируем новую пару токенов
	newAccessToken, err := s.jwtManager.GenerateAccessToken(user.ID, string(user.Role))
	if err != nil {
		return "", "", err
	}

	newRefreshToken, err := s.jwtManager.GenerateRefreshToken(user.ID)
	if err != nil {
		return "", "", err
	}

	// Сохраняем новый refresh токен
	newRefreshTokenModel := &model.RefreshToken{
		UserID:    user.ID,
		Token:     newRefreshToken,
		ExpiresAt: time.Now().Add(720 * time.Hour),
		Revoked:   false,
	}

	if err := s.refreshRepo.Create(newRefreshTokenModel); err != nil {
		return "", "", err
	}

	return newAccessToken, newRefreshToken, nil
}

// Logout выполняет выход пользователя (отзывает refresh токен)
func (s *AuthService) Logout(refreshToken string) error {
	// Находим токен
	token, err := s.refreshRepo.FindByToken(refreshToken)
	if err != nil {
		return err
	}

	// Отзываем токен
	return s.refreshRepo.Revoke(token.ID)
}

// LogoutAll выполняет выход со всех устройств
func (s *AuthService) LogoutAll(userID uint) error {
	return s.refreshRepo.RevokeAllByUserID(userID)
}

// GetUserByID возвращает пользователя по ID
func (s *AuthService) GetUserByID(userID uint) (*model.User, error) {
	return s.userRepo.FindByID(userID)
}

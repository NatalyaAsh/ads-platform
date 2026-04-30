package jwt

import (
	"errors"
	"fmt"
	"time"

	"github.com/NatalyaAsh/ads-platform/services/auth-user/pkg/config"
	"github.com/golang-jwt/jwt/v5"
)

// Claims - структура утверждений JWT токена
type Claims struct {
	UserID uint   `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// JWTManager - менеджер для работы с JWT токенами
type JWTManager struct {
	config *config.JWTConfig
}

// NewJWTManager создает новый менеджер JWT
func NewJWTManager(cfg *config.JWTConfig) *JWTManager {
	return &JWTManager{
		config: cfg,
	}
}

// GenerateAccessToken генерирует access токен
func (m *JWTManager) GenerateAccessToken(userID uint, role string) (string, error) {
	claims := &Claims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(m.config.AccessTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    m.config.Issuer,
			ID:        fmt.Sprintf("%d-%d", userID, time.Now().UnixNano()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(m.config.AccessSecret))
}

// GenerateRefreshToken генерирует refresh токен (без роли)
func (m *JWTManager) GenerateRefreshToken(userID uint) (string, error) {
	claims := &Claims{
		UserID: userID,
		Role:   "", // refresh токен не содержит роль
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(m.config.RefreshTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    m.config.Issuer,
			ID:        fmt.Sprintf("%d-%d", userID, time.Now().UnixNano()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(m.config.RefreshSecret))
}

// ValidateAccessToken валидирует access токен
func (m *JWTManager) ValidateAccessToken(tokenString string) (*Claims, error) {
	return m.validateToken(tokenString, m.config.AccessSecret)
}

// ValidateRefreshToken валидирует refresh токен
func (m *JWTManager) ValidateRefreshToken(tokenString string) (*Claims, error) {
	return m.validateToken(tokenString, m.config.RefreshSecret)
}

// validateToken общая логика валидации токена
func (m *JWTManager) validateToken(tokenString string, secret string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		// Проверяем метод подписи
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, errors.New("token expired")
		}
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

// GetUserIDFromAccessToken извлекает user_id из access токена
func (m *JWTManager) GetUserIDFromAccessToken(tokenString string) (uint, error) {
	claims, err := m.ValidateAccessToken(tokenString)
	if err != nil {
		return 0, err
	}
	return claims.UserID, nil
}

// GetUserIDFromRefreshToken извлекает user_id из refresh токена
func (m *JWTManager) GetUserIDFromRefreshToken(tokenString string) (uint, error) {
	claims, err := m.ValidateRefreshToken(tokenString)
	if err != nil {
		return 0, err
	}
	return claims.UserID, nil
}

// RefreshTokens генерирует новую пару токенов по refresh токену
func (m *JWTManager) RefreshTokens(refreshToken string) (string, string, error) {
	// Валидируем refresh токен
	claims, err := m.ValidateRefreshToken(refreshToken)
	if err != nil {
		return "", "", err
	}

	// Генерируем новую пару токенов
	newAccessToken, err := m.GenerateAccessToken(claims.UserID, claims.Role)
	if err != nil {
		return "", "", err
	}

	newRefreshToken, err := m.GenerateRefreshToken(claims.UserID)
	if err != nil {
		return "", "", err
	}

	return newAccessToken, newRefreshToken, nil
}

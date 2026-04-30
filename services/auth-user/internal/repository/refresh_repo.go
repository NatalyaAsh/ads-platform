package repository

import (
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"

	"github.com/NatalyaAsh/ads-platform/services/auth-user/internal/model"
)

// RefreshTokenRepository - репозиторий для работы с refresh токенами
type RefreshTokenRepository struct {
	db *gorm.DB
}

// NewRefreshTokenRepository создает новый репозиторий refresh токенов
func NewRefreshTokenRepository(db *gorm.DB) *RefreshTokenRepository {
	return &RefreshTokenRepository{db: db}
}

// Create сохраняет новый refresh токен
func (r *RefreshTokenRepository) Create(token *model.RefreshToken) error {
	if err := r.db.Create(token).Error; err != nil {
		return fmt.Errorf("create refresh token: %w", err)
	}
	return nil
}

// FindByToken находит токен по его значению
func (r *RefreshTokenRepository) FindByToken(token string) (*model.RefreshToken, error) {
	var refreshToken model.RefreshToken

	err := r.db.Where("token = ?", token).First(&refreshToken).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, model.ErrTokenNotFound
		}
		return nil, fmt.Errorf("find token by value: %w", err)
	}

	return &refreshToken, nil
}

// FindActiveByUserID находит все активные (не истекшие и не отозванные) токены пользователя
func (r *RefreshTokenRepository) FindActiveByUserID(userID uint) ([]model.RefreshToken, error) {
	var tokens []model.RefreshToken

	err := r.db.
		Where("user_id = ? AND revoked = ? AND expires_at > ?", userID, false, time.Now()).
		Find(&tokens).Error
	if err != nil {
		return nil, fmt.Errorf("find active tokens: %w", err)
	}

	return tokens, nil
}

// Revoke отзывает конкретный токен
func (r *RefreshTokenRepository) Revoke(id uint) error {
	result := r.db.Model(&model.RefreshToken{}).
		Where("id = ?", id).
		Update("revoked", true)

	if result.Error != nil {
		return fmt.Errorf("revoke token: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return model.ErrTokenNotFound
	}
	return nil
}

// RevokeByTokenValue отзывает токен по его значению
func (r *RefreshTokenRepository) RevokeByTokenValue(token string) error {
	result := r.db.Model(&model.RefreshToken{}).
		Where("token = ?", token).
		Update("revoked", true)

	if result.Error != nil {
		return fmt.Errorf("revoke token by value: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return model.ErrTokenNotFound
	}
	return nil
}

// RevokeAllByUserID отзывает ВСЕ токены пользователя (при смене пароля, выходе со всех устройств)
func (r *RefreshTokenRepository) RevokeAllByUserID(userID uint) error {
	result := r.db.Model(&model.RefreshToken{}).
		Where("user_id = ? AND revoked = ?", userID, false).
		Update("revoked", true)

	if result.Error != nil {
		return fmt.Errorf("revoke all tokens: %w", result.Error)
	}
	return nil
}

// DeleteExpired удаляет все просроченные токены (можно запускать по расписанию)
func (r *RefreshTokenRepository) DeleteExpired() error {
	result := r.db.Where("expires_at < ?", time.Now()).Delete(&model.RefreshToken{})
	if result.Error != nil {
		return fmt.Errorf("delete expired tokens: %w", result.Error)
	}
	return nil
}

// CleanupUserTokens удаляет все токены пользователя (для полной очистки)
func (r *RefreshTokenRepository) CleanupUserTokens(userID uint) error {
	result := r.db.Where("user_id = ?", userID).Delete(&model.RefreshToken{})
	if result.Error != nil {
		return fmt.Errorf("cleanup user tokens: %w", result.Error)
	}
	return nil
}

// DeleteByUserID удаляет все токены пользователя (полная очистка)
func (r *RefreshTokenRepository) DeleteByUserID(userID uint) error {
	result := r.db.Where("user_id = ?", userID).Delete(&model.RefreshToken{})
	if result.Error != nil {
		return fmt.Errorf("delete user tokens: %w", result.Error)
	}
	return nil
}

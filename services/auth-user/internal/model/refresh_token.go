package model

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

// RefreshToken - модель refresh токена
type RefreshToken struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"not null;index" json:"user_id"`
	Token     string    `gorm:"uniqueIndex;size:512;not null" json:"token"`
	ExpiresAt time.Time `gorm:"not null;index" json:"expires_at"`
	Revoked   bool      `gorm:"default:false;not null;index" json:"revoked"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`

	// Связь с пользователем (опционально)
	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// TableName указывает имя таблицы в БД
func (RefreshToken) TableName() string {
	return "refresh_tokens"
}

// Validate проверяет корректность модели
func (rt *RefreshToken) Validate() error {
	if rt.UserID == 0 {
		return errors.New("user_id is required")
	}
	if rt.Token == "" {
		return errors.New("token is required")
	}
	if rt.ExpiresAt.IsZero() {
		return errors.New("expires_at is required")
	}
	return nil
}

// IsExpired проверяет, истек ли срок действия токена
func (rt *RefreshToken) IsExpired() bool {
	return time.Now().After(rt.ExpiresAt)
}

// IsValid проверяет, валиден ли токен (не истек и не отозван)
func (rt *RefreshToken) IsValid() bool {
	return !rt.IsExpired() && !rt.Revoked
}

// BeforeCreate - GORM хук перед созданием
func (rt *RefreshToken) BeforeCreate(tx *gorm.DB) error {
	return rt.Validate()
}

// Revoke отзывает токен
func (rt *RefreshToken) Revoke() {
	rt.Revoked = true
}

// Extend продлевает срок действия токена
func (rt *RefreshToken) Extend(duration time.Duration) {
	rt.ExpiresAt = time.Now().Add(duration)
}

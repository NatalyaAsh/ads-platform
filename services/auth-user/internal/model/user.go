package model

import (
	"errors"
	"fmt"
	"regexp"
	"time"

	"gorm.io/gorm"
)

// Role - тип для ролей пользователя
type Role string

const (
	RoleUser  Role = "user"
	RoleAdmin Role = "admin"
)

// String возвращает строковое представление роли
func (r Role) String() string {
	return string(r)
}

// IsValid проверяет, корректная ли роль
func (r Role) IsValid() bool {
	switch r {
	case RoleUser, RoleAdmin:
		return true
	}
	return false
}

// User - модель пользователя
type User struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	Email        string    `gorm:"uniqueIndex;size:255;not null" json:"email"`
	PasswordHash string    `gorm:"not null" json:"-"` // "-" скрывает поле в JSON
	Role         Role      `gorm:"type:varchar(20);default:'user';not null" json:"role"`
	IsBlocked    bool      `gorm:"default:false;not null" json:"is_blocked"`
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// TableName указывает имя таблицы в БД
func (User) TableName() string {
	return "users"
}

// ValidateEmail проверяет корректность email
func (u *User) ValidateEmail() error {
	if u.Email == "" {
		return errors.New("email is required")
	}

	// Простая валидация email через regex
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`)
	if !emailRegex.MatchString(u.Email) {
		return errors.New("invalid email format")
	}

	return nil
}

// ValidateRole проверяет корректность роли
func (u *User) ValidateRole() error {
	if !u.Role.IsValid() {
		return fmt.Errorf("invalid role: %s. Must be 'user' or 'admin'", u.Role)
	}
	return nil
}

// Validate проверяет все поля пользователя
func (u *User) Validate() error {
	if err := u.ValidateEmail(); err != nil {
		return err
	}
	if err := u.ValidateRole(); err != nil {
		return err
	}
	return nil
}

// IsAdmin проверяет, является ли пользователь администратором
func (u *User) IsAdmin() bool {
	return u.Role == RoleAdmin && !u.IsBlocked
}

// IsActive проверяет, активен ли пользователь (не заблокирован)
func (u *User) IsActive() bool {
	return !u.IsBlocked
}

// BeforeCreate - GORM хук, вызываемый перед созданием записи
func (u *User) BeforeCreate(tx *gorm.DB) error {
	return u.Validate()
}

// BeforeUpdate - GORM хук, вызываемый перед обновлением
func (u *User) BeforeUpdate(tx *gorm.DB) error {
	return u.Validate()
}

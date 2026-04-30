package repository

import (
	"errors"
	"fmt"

	"gorm.io/gorm"

	"github.com/NatalyaAsh/ads-platform/services/auth-user/internal/model"
)

// UserRepository - репозиторий для работы с пользователями
type UserRepository struct {
	db *gorm.DB
}

// NewUserRepository создает новый репозиторий пользователей
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Create создает нового пользователя
func (r *UserRepository) Create(user *model.User) error {
	// Проверяем, не существует ли пользователь с таким email
	var existing model.User
	if err := r.db.Where("email = ?", user.Email).First(&existing).Error; err == nil {
		return model.ErrUserAlreadyExists
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("check existing user: %w", err)
	}

	// Создаем пользователя
	if err := r.db.Create(user).Error; err != nil {
		return fmt.Errorf("create user: %w", err)
	}

	return nil
}

// FindByEmail находит пользователя по email
func (r *UserRepository) FindByEmail(email string) (*model.User, error) {
	var user model.User

	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, model.ErrUserNotFound
		}
		return nil, fmt.Errorf("find user by email: %w", err)
	}

	return &user, nil
}

// FindByID находит пользователя по ID
func (r *UserRepository) FindByID(id uint) (*model.User, error) {
	var user model.User

	err := r.db.First(&user, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, model.ErrUserNotFound
		}
		return nil, fmt.Errorf("find user by id: %w", err)
	}

	return &user, nil
}

// Update обновляет данные пользователя
func (r *UserRepository) Update(user *model.User) error {
	result := r.db.Save(user)
	if result.Error != nil {
		return fmt.Errorf("update user: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return model.ErrUserNotFound
	}
	return nil
}

// Block блокирует пользователя
func (r *UserRepository) Block(id uint) error {
	result := r.db.Model(&model.User{}).Where("id = ?", id).Update("is_blocked", true)
	if result.Error != nil {
		return fmt.Errorf("block user: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return model.ErrUserNotFound
	}
	return nil
}

// Unblock разблокирует пользователя
func (r *UserRepository) Unblock(id uint) error {
	result := r.db.Model(&model.User{}).Where("id = ?", id).Update("is_blocked", false)
	if result.Error != nil {
		return fmt.Errorf("unblock user: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return model.ErrUserNotFound
	}
	return nil
}

// UpdateRole обновляет роль пользователя
func (r *UserRepository) UpdateRole(id uint, role model.Role) error {
	if !role.IsValid() {
		return model.ErrInvalidRole
	}

	result := r.db.Model(&model.User{}).Where("id = ?", id).Update("role", role)
	if result.Error != nil {
		return fmt.Errorf("update role: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return model.ErrUserNotFound
	}
	return nil
}

// ExistsByEmail проверяет существование пользователя по email
func (r *UserRepository) ExistsByEmail(email string) (bool, error) {
	var count int64
	err := r.db.Model(&model.User{}).Where("email = ?", email).Count(&count).Error
	if err != nil {
		return false, fmt.Errorf("check email exists: %w", err)
	}
	return count > 0, nil
}

// List возвращает список пользователей с пагинацией
func (r *UserRepository) List(offset, limit int) ([]model.User, int64, error) {
	var users []model.User
	var total int64

	// Подсчет общего количества
	if err := r.db.Model(&model.User{}).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count users: %w", err)
	}

	// Получение пользователей с пагинацией
	err := r.db.Offset(offset).Limit(limit).Find(&users).Error
	if err != nil {
		return nil, 0, fmt.Errorf("list users: %w", err)
	}

	return users, total, nil
}

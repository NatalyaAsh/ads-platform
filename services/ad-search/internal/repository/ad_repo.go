package repository

import (
	"errors"
	"fmt"

	"github.com/NatalyaAsh/ads-platform/services/ad-search/internal/model"
	"gorm.io/gorm"
)

type AdRepository struct {
	db *gorm.DB
}

func NewAdRepository(db *gorm.DB) *AdRepository {
	return &AdRepository{db: db}
}

func (r *AdRepository) Create(ad *model.Ad) error {
	// 1. Проверить обязательные поля
	if ad.Title == "" {
		return fmt.Errorf("title must not be empty")
	}
	if len(ad.Title) > 255 {
		return fmt.Errorf("title must not exceed 255 characters")
	}
	if ad.Price <= 0 {
		return fmt.Errorf("price must be greater than zero")
	}
	if ad.UserID == 0 {
		return fmt.Errorf("user_id must be not null")
	}
	if ad.CategoryID == 0 {
		return fmt.Errorf("category_id must be not null")
	}

	// 2. Проверить существование категории
	var category model.Category
	if err := r.db.First(&category, ad.CategoryID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.ErrCategoryNotFound
		}
		return fmt.Errorf("check category: %w", err)
	}

	// 3. Создать объявление
	if err := r.db.Create(ad).Error; err != nil {
		return fmt.Errorf("create advert: %w", err)
	}
	return nil
}

func (r *AdRepository) Update(ad *model.Ad) error {
	var existing model.Ad
	// 1. Проверить, что объявление существует
	if err := r.db.First(&existing, ad.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.ErrAdNotFound
		}
		return fmt.Errorf("check ad: %w", err)
	}

	// 2. Проверить обязательные поля
	if ad.Title == "" {
		return fmt.Errorf("title must not be empty")
	}
	if len(ad.Title) > 255 {
		return fmt.Errorf("title must not exceed 255 characters")
	}
	if ad.Price <= 0 {
		return fmt.Errorf("price must be greater than zero")
	}
	if ad.UserID == 0 {
		return fmt.Errorf("user_id must be not null")
	}
	if ad.CategoryID == 0 {
		return fmt.Errorf("category_id must be not null")
	}

	// 3. Проверить существование категории
	var category model.Category
	if err := r.db.First(&category, ad.CategoryID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.ErrCategoryNotFound
		}
		return fmt.Errorf("check category: %w", err)
	}

	// 4. Обновить объявление
	result := r.db.Model(&existing).Updates(ad)
	if result.Error != nil {
		return fmt.Errorf("update advertisement: %w", result.Error)
	}
	// 5. Вернуть ошибку, если ничего не обновилось
	if result.RowsAffected == 0 {
		return model.ErrAdNotFound
	}

	return nil
}

func (r *AdRepository) Delete(ad *model.Ad) error {
	var existing model.Ad
	// 1. Проверить, что объявление существует
	if err := r.db.First(&existing, ad.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.ErrAdNotFound
		}
		return fmt.Errorf("check ad: %w", err)
	}

	// 2. Проверить статус объявления, что оно ещё не удалено
	if existing.Status == model.AdStatusDeleted {
		return fmt.Errorf("ad already deleted")
	}

	// 3. Обновить статус на deleted
	result := r.db.Model(&existing).Update("status", model.AdStatusDeleted)
	if result.Error != nil {
		return fmt.Errorf("update advertisement: %w", result.Error)
	}
	// 4. Вернуть ошибку, если ничего не обновилось
	if result.RowsAffected == 0 {
		return model.ErrAdNotFound
	}

	return nil
}

func (r *AdRepository) FindByID(id uint) (*model.Ad, error) {
	var ad model.Ad
	err := r.db.Where("id = ? AND status != ?", id, model.AdStatusDeleted).
		Preload("Category").
		First(&ad).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, model.ErrAdNotFound
		}
		return nil, fmt.Errorf("find ad by id: %w", err)
	}
	return &ad, nil
}

func (r *AdRepository) FindByUserID(userID uint) ([]model.Ad, error) {
	var ads []model.Ad
	err := r.db.Where("user_id = ? AND status != ?", userID, model.AdStatusDeleted).Find(&ads).Error
	if err != nil {
		return nil, fmt.Errorf("find ads by user_id: %w", err)
	}
	return ads, nil
}

func (r *AdRepository) List(filters model.AdListFilters) ([]model.Ad, int64, error) {
	var ads []model.Ad
	var total int64

	// 1. Базовый запрос (не показываем удалённые)
	db := r.db.Model(&model.Ad{}).Where("status != ? AND status != ?", model.AdStatusDeleted, model.AdStatusDeleted)

	// 2. Применяем фильтры
	if filters.CategoryID != 0 {
		db = db.Where("category_id = ?", filters.CategoryID)
	}
	if filters.UserID != 0 {
		db = db.Where("user_id = ?", filters.UserID)
	}
	if filters.Status != "" {
		db = db.Where("status = ?", filters.Status)
	}
	if filters.MinPrice != nil {
		db = db.Where("price >= ?", *filters.MinPrice)
	}
	if filters.MaxPrice != nil {
		db = db.Where("price <= ?", *filters.MaxPrice)
	}

	// 3. Считаем общее количество (до пагинации)
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count ads: %w", err)
	}

	// 4. Применяем сортировку
	if filters.SortBy != "" {
		order := filters.SortBy
		if filters.SortOrder == "desc" {
			order += " DESC"
		}
		db = db.Order(order)
	} else {
		db = db.Order("created_at DESC") // сортировка по умолчанию
	}

	// 5. Пагинация
	if filters.Limit > 0 {
		db = db.Limit(filters.Limit)
	}
	if filters.Offset > 0 {
		db = db.Offset(filters.Offset)
	}

	// 6. Выполняем запрос с подгрузкой категории
	if err := db.Preload("Category").Find(&ads).Error; err != nil {
		return nil, 0, fmt.Errorf("list ads: %w", err)
	}

	return ads, total, nil
}

// DeleteByID мягко удаляет объявление по ID
func (r *AdRepository) DeleteByID(id uint) error {
	// Сначала проверяем, существует ли объявление
	var ad model.Ad
	if err := r.db.First(&ad, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.ErrAdNotFound
		}
		return fmt.Errorf("find ad: %w", err)
	}

	// Проверяем, не удалено ли уже
	if ad.Status == model.AdStatusDeleted {
		return model.ErrAdNotFound
	}

	// Обновляем статус
	result := r.db.Model(&model.Ad{}).Where("id = ?", id).Update("status", model.AdStatusDeleted)
	if result.Error != nil {
		return fmt.Errorf("delete ad: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return model.ErrAdNotFound
	}
	return nil
}

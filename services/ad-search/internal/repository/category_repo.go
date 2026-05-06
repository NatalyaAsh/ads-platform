package repository

import (
	"errors"
	"fmt"

	"github.com/NatalyaAsh/ads-platform/services/ad-search/internal/model"
	"gorm.io/gorm"
)

type CategoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

func (r *CategoryRepository) FindByID(id uint) (*model.Category, error) {
	var category model.Category
	if err := r.db.First(&category, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, model.ErrCategoryNotFound
		}
		return nil, fmt.Errorf("find category by id: %w", err)
	}
	return &category, nil
}

func (r *CategoryRepository) FindByName(name string) (*model.Category, error) {
	var category model.Category
	err := r.db.Where("name = ?", name).First(&category).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, model.ErrCategoryNotFound
		}
		return nil, fmt.Errorf("find category by name: %w", err)
	}
	return &category, nil
}

func (r *CategoryRepository) FindBySlug(slug string) (*model.Category, error) {
	var category model.Category
	err := r.db.Where("slug = ?", slug).First(&category).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, model.ErrCategoryNotFound
		}
		return nil, fmt.Errorf("find category by slug: %w", err)
	}
	return &category, nil
}

func (r *CategoryRepository) Create(category *model.Category) error {
	// 1. Проверить обязательные поля
	if category.Name == "" {
		return fmt.Errorf("name must not be empty")
	}
	if len(category.Name) > 100 {
		return fmt.Errorf("name must not exceed 100 characters")
	}

	if category.Slug == "" {
		return fmt.Errorf("slug must not be empty")
	}
	if len(category.Slug) > 100 {
		return fmt.Errorf("slug must not exceed 100 characters")
	}

	// 2. Проверить существование категории
	var existing model.Category
	if err := r.db.Where("name = ?", category.Name).First(&existing).Error; err == nil {
		return model.ErrCategoryAlreadyExists
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("check name: %w", err)
	}

	// 3. Проверить уникальность slug
	if err := r.db.Where("slug = ?", category.Slug).First(&existing).Error; err == nil {
		return model.ErrCategoryAlreadyExists
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("check slug: %w", err)
	}

	// 4. Создать объявление
	if err := r.db.Create(category).Error; err != nil {
		return fmt.Errorf("create category: %w", err)
	}
	return nil
}

func (r *CategoryRepository) Update(category *model.Category) error {
	var existing model.Category
	// 1. Проверить, что категория существует
	if err := r.db.First(&existing, category.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.ErrCategoryNotFound
		}
		return fmt.Errorf("check category: %w", err)
	}

	// 2. Проверить обязательные поля
	if category.Name == "" {
		return fmt.Errorf("name must not be empty")
	}
	if len(category.Name) > 100 {
		return fmt.Errorf("name must not exceed 100 characters")
	}

	if category.Slug == "" {
		return fmt.Errorf("slug must not be empty")
	}
	if len(category.Slug) > 100 {
		return fmt.Errorf("slug must not exceed 100 characters")
	}

	// 3. Проверить уникальность имени (если изменилось)
	if category.Name != existing.Name {
		var other model.Category
		if err := r.db.Where("name = ?", category.Name).First(&other).Error; err == nil {
			return model.ErrCategoryAlreadyExists
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("check name uniqueness: %w", err)
		}
	}

	// 4. Проверить уникальность slug (если изменился)
	if category.Slug != existing.Slug {
		var other model.Category
		if err := r.db.Where("slug = ?", category.Slug).First(&other).Error; err == nil {
			return model.ErrCategoryAlreadyExists
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("check slug uniqueness: %w", err)
		}
	}

	// 5. Обновить категорию
	result := r.db.Model(&existing).Updates(category)
	if result.Error != nil {
		return fmt.Errorf("update category: %w", result.Error)
	}
	// 6. Вернуть ошибку, если ничего не обновилось
	if result.RowsAffected == 0 {
		return model.ErrCategoryNotFound
	}

	return nil
}

func (r *CategoryRepository) Delete(id uint) error {
	// 1. Проверить, что категория существует
	var category model.Category
	if err := r.db.First(&category, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.ErrCategoryNotFound
		}
		return fmt.Errorf("check category: %w", err)
	}

	// 2. Проверить, есть ли объявления с этой категорией
	var count int64
	if err := r.db.Model(&model.Ad{}).Where("category_id = ?", id).Count(&count).Error; err != nil {
		return fmt.Errorf("check ads in category: %w", err)
	}
	if count > 0 {
		return model.ErrCategoryHasAds
	}

	// 3. Удалить категорию (жёсткое удаление)
	if err := r.db.Delete(&category).Error; err != nil {
		return fmt.Errorf("delete category: %w", err)
	}

	return nil
}

func (r *CategoryRepository) List() ([]model.Category, error) {
	var categories []model.Category
	err := r.db.Order("name ASC").Find(&categories).Error
	if err != nil {
		return nil, fmt.Errorf("list categories: %w", err)
	}
	return categories, nil
}

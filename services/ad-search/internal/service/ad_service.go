package service

import (
	"errors"
	"fmt"

	"github.com/NatalyaAsh/ads-platform/services/ad-search/internal/model"
	"github.com/NatalyaAsh/ads-platform/services/ad-search/internal/repository"
)

type AdService struct {
	adRepo       *repository.AdRepository
	categoryRepo *repository.CategoryRepository
}

func NewAdService(
	adRepo *repository.AdRepository,
	categoryRepo *repository.CategoryRepository,
) *AdService {
	return &AdService{
		adRepo:       adRepo,
		categoryRepo: categoryRepo,
	}
}

func (s *AdService) CreateAd(req *model.CreateAdRequest) (*model.Ad, error) {
	// 1. Валидация полей
	if req.Title == "" {
		return nil, errors.New("title is required")
	}
	if len(req.Title) > 200 {
		return nil, errors.New("title must not exceed 200 characters")
	}
	if req.Price <= 0 {
		return nil, model.ErrInvalidPrice
	}
	if req.UserID == 0 {
		return nil, errors.New("user_id is required")
	}
	if req.CategoryID == 0 {
		return nil, errors.New("category_id is required")
	}

	// 2. Проверка существования категории
	_, err := s.categoryRepo.FindByID(req.CategoryID)
	if err != nil {
		return nil, model.ErrCategoryNotFound
	}

	// 3. Создание объявления
	ad := &model.Ad{
		Title:       req.Title,
		Description: req.Description,
		Price:       req.Price,
		UserID:      req.UserID,
		CategoryID:  req.CategoryID,
		Status:      model.AdStatusDraft,
		Views:       0,
	}

	if err := s.adRepo.Create(ad); err != nil {
		return nil, fmt.Errorf("failed to create ad: %w", err)
	}

	// 4. TODO: Отправить событие в RabbitMQ
	// 5. TODO: Индексировать в Elasticsearch

	return ad, nil
}

func (s *AdService) GetAdByID(id uint) (*model.Ad, error) {
	ad, err := s.adRepo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("ad not found: %w", err)
	}
	return ad, nil
}

// GetUserAds возвращает все объявления пользователя
func (s *AdService) GetUserAds(userID uint) ([]model.Ad, error) {
	if userID == 0 {
		return nil, errors.New("user_id is required")
	}

	ads, err := s.adRepo.FindByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user ads: %w", err)
	}
	return ads, nil
}

// UpdateAd обновляет существующее объявление
func (s *AdService) UpdateAd(req *model.UpdateAdRequest) (*model.Ad, error) {
	// 1. Проверить, что объявление существует
	existingAd, err := s.adRepo.FindByID(req.ID)
	if err != nil {
		return nil, model.ErrAdNotFound
	}

	// 2. Валидация полей
	if req.Title == "" {
		return nil, errors.New("title is required")
	}
	if len(req.Title) > 255 {
		return nil, errors.New("title must not exceed 200 characters")
	}
	if req.Price <= 0 {
		return nil, model.ErrInvalidPrice
	}
	if req.CategoryID == 0 {
		return nil, errors.New("category_id is required")
	}

	// 3. Проверить существование категории (если изменилась)
	if req.CategoryID != existingAd.CategoryID {
		_, err := s.categoryRepo.FindByID(req.CategoryID)
		if err != nil {
			return nil, model.ErrCategoryNotFound
		}
	}

	// 4. Обновить поля
	existingAd.Title = req.Title
	existingAd.Description = req.Description
	existingAd.Price = req.Price
	existingAd.CategoryID = req.CategoryID

	// 5. Сохранить изменения
	if err := s.adRepo.Update(existingAd); err != nil {
		return nil, fmt.Errorf("failed to update ad: %w", err)
	}

	return existingAd, nil
}

// DeleteAd мягко удаляет объявление (меняет статус на deleted)
func (s *AdService) DeleteAd(id uint) error {
	if err := s.adRepo.DeleteByID(id); err != nil {
		return model.ErrAdNotFound
	}
	return nil
}

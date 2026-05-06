package service

import (
	"testing"

	"github.com/NatalyaAsh/ads-platform/services/ad-search/internal/model"
	"github.com/NatalyaAsh/ads-platform/services/ad-search/internal/repository"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// setupTestDB создаёт тестовую БД
func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open db: %v", err)
	}

	// Создаём таблицы
	if err := db.AutoMigrate(&model.Ad{}, &model.Category{}); err != nil {
		t.Fatalf("Failed to migrate: %v", err)
	}

	return db
}

// createTestCategory создаёт тестовую категорию
func createTestCategory(t *testing.T, db *gorm.DB) *model.Category {
	category := &model.Category{
		Name: "Electronics",
		Slug: "electronics",
	}
	if err := db.Create(category).Error; err != nil {
		t.Fatalf("Failed to create test category: %v", err)
	}
	return category
}

func TestAdService_CreateAd(t *testing.T) {
	db := setupTestDB(t)

	adRepo := repository.NewAdRepository(db)
	categoryRepo := repository.NewCategoryRepository(db)
	adService := NewAdService(adRepo, categoryRepo)

	// Создаём тестовую категорию
	category := createTestCategory(t, db)

	tests := []struct {
		name    string
		req     *model.CreateAdRequest
		wantErr bool
		errMsg  string
	}{
		{
			name: "successful creation",
			req: &model.CreateAdRequest{
				Title:       "Test Ad",
				Description: "Test Description",
				Price:       100.50,
				UserID:      1,
				CategoryID:  category.ID,
			},
			wantErr: false,
		},
		{
			name: "empty title",
			req: &model.CreateAdRequest{
				Title:       "",
				Description: "Test",
				Price:       100,
				UserID:      1,
				CategoryID:  category.ID,
			},
			wantErr: true,
			errMsg:  "title is required",
		},
		{
			name: "title too long",
			req: &model.CreateAdRequest{
				Title:       "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat.",
				Description: "Test",
				Price:       100,
				UserID:      1,
				CategoryID:  category.ID,
			},
			wantErr: true,
			errMsg:  "title must not exceed 200 characters",
		},
		{
			name: "price zero",
			req: &model.CreateAdRequest{
				Title:       "Test",
				Description: "Test",
				Price:       0,
				UserID:      1,
				CategoryID:  category.ID,
			},
			wantErr: true,
			errMsg:  "price must be greater than zero",
		},
		{
			name: "price negative",
			req: &model.CreateAdRequest{
				Title:       "Test",
				Description: "Test",
				Price:       -10,
				UserID:      1,
				CategoryID:  category.ID,
			},
			wantErr: true,
			errMsg:  "price must be greater than zero",
		},
		{
			name: "user_id zero",
			req: &model.CreateAdRequest{
				Title:       "Test",
				Description: "Test",
				Price:       100,
				UserID:      0,
				CategoryID:  category.ID,
			},
			wantErr: true,
			errMsg:  "user_id is required",
		},
		{
			name: "category_id zero",
			req: &model.CreateAdRequest{
				Title:       "Test",
				Description: "Test",
				Price:       100,
				UserID:      1,
				CategoryID:  0,
			},
			wantErr: true,
			errMsg:  "category_id is required",
		},
		{
			name: "category not exists",
			req: &model.CreateAdRequest{
				Title:       "Test",
				Description: "Test",
				Price:       100,
				UserID:      1,
				CategoryID:  999,
			},
			wantErr: true,
			errMsg:  "category not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ad, err := adService.CreateAd(tt.req)

			if tt.wantErr {
				if err == nil {
					t.Errorf("CreateAd() expected error, got nil")
				}
				if tt.errMsg != "" && err.Error() != tt.errMsg {
					t.Errorf("CreateAd() error = %v, want %v", err.Error(), tt.errMsg)
				}
				if ad != nil {
					t.Errorf("CreateAd() returned non-nil ad on error")
				}
			} else {
				if err != nil {
					t.Errorf("CreateAd() error = %v, want nil", err)
				}
				if ad == nil {
					t.Errorf("CreateAd() returned nil ad")
				}
				if ad.Title != tt.req.Title {
					t.Errorf("CreateAd() title = %v, want %v", ad.Title, tt.req.Title)
				}
				if ad.Price != tt.req.Price {
					t.Errorf("CreateAd() price = %v, want %v", ad.Price, tt.req.Price)
				}
				if ad.Status != model.AdStatusDraft {
					t.Errorf("CreateAd() status = %v, want %v", ad.Status, model.AdStatusDraft)
				}
			}
		})
	}
}

func TestAdService_GetAdByID(t *testing.T) {
	db := setupTestDB(t)

	adRepo := repository.NewAdRepository(db)
	categoryRepo := repository.NewCategoryRepository(db)
	adService := NewAdService(adRepo, categoryRepo)

	// Создаём категорию
	category := createTestCategory(t, db)

	// Создаём объявление через сервис
	createReq := &model.CreateAdRequest{
		Title:       "Test Ad for Get",
		Description: "Test Description",
		Price:       100,
		UserID:      1,
		CategoryID:  category.ID,
	}
	createdAd, err := adService.CreateAd(createReq)
	if err != nil {
		t.Fatalf("Failed to create test ad: %v", err)
	}

	tests := []struct {
		name    string
		id      uint
		wantErr bool
		errMsg  string
	}{
		{
			name:    "existing ad",
			id:      createdAd.ID,
			wantErr: false,
		},
		{
			name:    "non-existing ad",
			id:      999,
			wantErr: true,
			errMsg:  "ad not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ad, err := adService.GetAdByID(tt.id)

			if tt.wantErr {
				if err == nil {
					t.Errorf("GetAdByID() expected error, got nil")
				}
				if ad != nil {
					t.Errorf("GetAdByID() returned non-nil ad on error")
				}
			} else {
				if err != nil {
					t.Errorf("GetAdByID() error = %v, want nil", err)
				}
				if ad == nil {
					t.Errorf("GetAdByID() returned nil ad")
				}
				if ad.ID != tt.id {
					t.Errorf("GetAdByID() ID = %v, want %v", ad.ID, tt.id)
				}
				if ad.Title != createReq.Title {
					t.Errorf("GetAdByID() Title = %v, want %v", ad.Title, createReq.Title)
				}
			}
		})
	}
}

func TestAdService_GetUserAds(t *testing.T) {
	db := setupTestDB(t)

	adRepo := repository.NewAdRepository(db)
	categoryRepo := repository.NewCategoryRepository(db)
	adService := NewAdService(adRepo, categoryRepo)

	// Создаём категорию
	category := createTestCategory(t, db)

	// Создаём 2 объявления для userID=1
	req1 := &model.CreateAdRequest{
		Title:       "User 1 Ad 1",
		Description: "Description 1",
		Price:       100,
		UserID:      1,
		CategoryID:  category.ID,
	}
	req2 := &model.CreateAdRequest{
		Title:       "User 1 Ad 2",
		Description: "Description 2",
		Price:       200,
		UserID:      1,
		CategoryID:  category.ID,
	}

	_, err := adService.CreateAd(req1)
	if err != nil {
		t.Fatalf("Failed to create ad 1: %v", err)
	}
	_, err = adService.CreateAd(req2)
	if err != nil {
		t.Fatalf("Failed to create ad 2: %v", err)
	}

	// Создаём объявление для другого пользователя (userID=2)
	req3 := &model.CreateAdRequest{
		Title:       "User 2 Ad",
		Description: "Description 3",
		Price:       300,
		UserID:      2,
		CategoryID:  category.ID,
	}
	_, err = adService.CreateAd(req3)
	if err != nil {
		t.Fatalf("Failed to create ad for user 2: %v", err)
	}

	tests := []struct {
		name      string
		userID    uint
		wantCount int
		wantErr   bool
		errMsg    string
	}{
		{
			name:      "user with 2 ads",
			userID:    1,
			wantCount: 2,
			wantErr:   false,
		},
		{
			name:      "user with 0 ads",
			userID:    3,
			wantCount: 0,
			wantErr:   false,
		},
		{
			name:      "invalid user_id",
			userID:    0,
			wantCount: 0,
			wantErr:   true,
			errMsg:    "user_id is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ads, err := adService.GetUserAds(tt.userID)

			if tt.wantErr {
				if err == nil {
					t.Errorf("GetUserAds() expected error, got nil")
				}
				if tt.errMsg != "" && err.Error() != tt.errMsg {
					t.Errorf("GetUserAds() error = %v, want %v", err.Error(), tt.errMsg)
				}
				if ads != nil {
					t.Errorf("GetUserAds() returned non-nil ads on error")
				}
			} else {
				if err != nil {
					t.Errorf("GetUserAds() error = %v, want nil", err)
				}
				if len(ads) != tt.wantCount {
					t.Errorf("GetUserAds() got %d ads, want %d", len(ads), tt.wantCount)
				}
				if tt.userID == 1 {
					// Проверяем, что все объявления принадлежат пользователю
					for _, ad := range ads {
						if ad.UserID != tt.userID {
							t.Errorf("GetUserAds() ad belongs to user %d, want %d", ad.UserID, tt.userID)
						}
					}
				}
			}
		})
	}
}

func TestAdService_UpdateAd(t *testing.T) {
	db := setupTestDB(t)

	adRepo := repository.NewAdRepository(db)
	categoryRepo := repository.NewCategoryRepository(db)
	adService := NewAdService(adRepo, categoryRepo)

	// Создаём категории
	category1 := createTestCategory(t, db)

	// Создаём вторую категорию
	category2 := &model.Category{
		Name: "Books",
		Slug: "books",
	}
	if err := db.Create(category2).Error; err != nil {
		t.Fatalf("Failed to create category2: %v", err)
	}

	// Создаём объявление
	createReq := &model.CreateAdRequest{
		Title:       "Original Title",
		Description: "Original Description",
		Price:       100,
		UserID:      1,
		CategoryID:  category1.ID,
	}
	createdAd, err := adService.CreateAd(createReq)
	if err != nil {
		t.Fatalf("Failed to create test ad: %v", err)
	}

	tests := []struct {
		name    string
		req     *model.UpdateAdRequest
		wantErr bool
		errMsg  string
	}{
		{
			name: "successful update",
			req: &model.UpdateAdRequest{
				ID:          createdAd.ID,
				Title:       "Updated Title",
				Description: "Updated Description",
				Price:       200,
				CategoryID:  category2.ID,
			},
			wantErr: false,
		},
		{
			name: "ad not found",
			req: &model.UpdateAdRequest{
				ID:          999,
				Title:       "New Title",
				Description: "New Desc",
				Price:       100,
				CategoryID:  category1.ID,
			},
			wantErr: true,
			errMsg:  "ad not found",
		},
		{
			name: "empty title",
			req: &model.UpdateAdRequest{
				ID:          createdAd.ID,
				Title:       "",
				Description: "Desc",
				Price:       100,
				CategoryID:  category1.ID,
			},
			wantErr: true,
			errMsg:  "title is required",
		},
		{
			name: "title too long",
			req: &model.UpdateAdRequest{
				ID:          createdAd.ID,
				Title:       "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.",
				Description: "Desc",
				Price:       100,
				CategoryID:  category1.ID,
			},
			wantErr: true,
			errMsg:  "title must not exceed 200 characters",
		},
		{
			name: "price zero",
			req: &model.UpdateAdRequest{
				ID:          createdAd.ID,
				Title:       "Valid Title",
				Description: "Desc",
				Price:       0,
				CategoryID:  category1.ID,
			},
			wantErr: true,
			errMsg:  "price must be greater than zero",
		},
		{
			name: "category_id zero",
			req: &model.UpdateAdRequest{
				ID:          createdAd.ID,
				Title:       "Valid Title",
				Description: "Desc",
				Price:       100,
				CategoryID:  0,
			},
			wantErr: true,
			errMsg:  "category_id is required",
		},
		{
			name: "category not found",
			req: &model.UpdateAdRequest{
				ID:          createdAd.ID,
				Title:       "Valid Title",
				Description: "Desc",
				Price:       100,
				CategoryID:  999,
			},
			wantErr: true,
			errMsg:  "category not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			updatedAd, err := adService.UpdateAd(tt.req)

			if tt.wantErr {
				if err == nil {
					t.Errorf("UpdateAd() expected error, got nil")
				}
				if tt.errMsg != "" && err.Error() != tt.errMsg {
					t.Errorf("UpdateAd() error = %v, want %v", err.Error(), tt.errMsg)
				}
				if updatedAd != nil {
					t.Errorf("UpdateAd() returned non-nil ad on error")
				}
			} else {
				if err != nil {
					t.Errorf("UpdateAd() error = %v, want nil", err)
				}
				if updatedAd == nil {
					t.Errorf("UpdateAd() returned nil ad")
				}
				if updatedAd.Title != tt.req.Title {
					t.Errorf("UpdateAd() Title = %v, want %v", updatedAd.Title, tt.req.Title)
				}
				if updatedAd.Price != tt.req.Price {
					t.Errorf("UpdateAd() Price = %v, want %v", updatedAd.Price, tt.req.Price)
				}
				if updatedAd.CategoryID != tt.req.CategoryID {
					t.Errorf("UpdateAd() CategoryID = %v, want %v", updatedAd.CategoryID, tt.req.CategoryID)
				}
			}
		})
	}
}

func TestAdService_DeleteAd(t *testing.T) {
	db := setupTestDB(t)

	adRepo := repository.NewAdRepository(db)
	categoryRepo := repository.NewCategoryRepository(db)
	adService := NewAdService(adRepo, categoryRepo)

	// Создаём категорию
	category := createTestCategory(t, db)

	// Создаём объявление
	createReq := &model.CreateAdRequest{
		Title:       "Ad to Delete",
		Description: "This ad will be deleted",
		Price:       100,
		UserID:      1,
		CategoryID:  category.ID,
	}
	createdAd, err := adService.CreateAd(createReq)
	if err != nil {
		t.Fatalf("Failed to create test ad: %v", err)
	}

	tests := []struct {
		name    string
		id      uint
		wantErr bool
		errMsg  string
	}{
		{
			name:    "successful delete",
			id:      createdAd.ID,
			wantErr: false,
		},
		{
			name:    "ad not found",
			id:      999,
			wantErr: true,
			errMsg:  "ad not found",
		},
		{
			name:    "delete already deleted ad",
			id:      createdAd.ID,
			wantErr: true,
			errMsg:  "ad not found", // после удаления объявление не найдётся
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := adService.DeleteAd(tt.id)

			if tt.wantErr {
				if err == nil {
					t.Errorf("DeleteAd() expected error, got nil")
				}
				if tt.errMsg != "" && err.Error() != tt.errMsg {
					t.Errorf("DeleteAd() error = %v, want %v", err.Error(), tt.errMsg)
				}
			} else {
				if err != nil {
					t.Errorf("DeleteAd() error = %v, want nil", err)
				}

				// Проверяем, что объявление действительно удалено
				_, err := adService.GetAdByID(tt.id)
				if err == nil {
					t.Errorf("DeleteAd() ad still exists after deletion")
				}
			}
		})
	}
}

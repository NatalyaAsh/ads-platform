package repository

import (
	"context"
	"testing"

	"github.com/NatalyaAsh/ads-platform/services/ad-search/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// setupTestMongoDB создаёт тестовую in-memory MongoDB
func setupTestMongoDB(t *testing.T) (*mongo.Client, *mongo.Collection, func()) {
	// Используем локальную MongoDB для тестов
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://localhost:27017"))
	require.NoError(t, err)

	db := client.Database("ad_media_test")
	collection := db.Collection("ad_media")

	// Очистка после тестов
	cleanup := func() {
		collection.Drop(context.Background())
		client.Disconnect(context.Background())
	}

	return client, collection, cleanup
}

func TestMediaRepository_Create(t *testing.T) {
	_, collection, cleanup := setupTestMongoDB(t)
	defer cleanup()

	repo := &MediaRepository{collection: collection}

	media := &model.AdMedia{
		AdID:      1,
		FilePath:  "uploads/ads/1/photo.jpg",
		FileName:  "photo.jpg",
		FileSize:  1024,
		MimeType:  "image/jpeg",
		IsPrimary: true,
	}

	err := repo.Create(context.Background(), media)
	assert.NoError(t, err)
	assert.NotNil(t, media.ID)
	assert.NotZero(t, media.CreatedAt)
}

func TestMediaRepository_FindByAdID(t *testing.T) {
	_, collection, cleanup := setupTestMongoDB(t)
	defer cleanup()

	repo := &MediaRepository{collection: collection}

	// Создаём два медиа для ad_id=1
	media1 := &model.AdMedia{AdID: 1, FilePath: "path1.jpg", FileName: "file1.jpg"}
	media2 := &model.AdMedia{AdID: 1, FilePath: "path2.jpg", FileName: "file2.jpg"}
	media3 := &model.AdMedia{AdID: 2, FilePath: "path3.jpg", FileName: "file3.jpg"}

	repo.Create(context.Background(), media1)
	repo.Create(context.Background(), media2)
	repo.Create(context.Background(), media3)

	// Ищем все медиа для ad_id=1
	result, err := repo.FindByAdID(context.Background(), 1)
	assert.NoError(t, err)
	assert.Len(t, result, 2)
}

func TestMediaRepository_DeleteByAdID(t *testing.T) {
	_, collection, cleanup := setupTestMongoDB(t)
	defer cleanup()

	repo := &MediaRepository{collection: collection}

	// Создаём медиа
	media := &model.AdMedia{AdID: 1, FilePath: "path.jpg", FileName: "file.jpg"}
	repo.Create(context.Background(), media)

	// Проверяем, что создалось
	result, _ := repo.FindByAdID(context.Background(), 1)
	assert.Len(t, result, 1)

	// Удаляем все медиа для ad_id=1
	err := repo.DeleteByAdID(context.Background(), 1)
	assert.NoError(t, err)

	// Проверяем, что удалилось
	result, _ = repo.FindByAdID(context.Background(), 1)
	assert.Len(t, result, 0)
}

func TestMediaRepository_DeleteByID(t *testing.T) {
	_, collection, cleanup := setupTestMongoDB(t)
	defer cleanup()

	repo := &MediaRepository{collection: collection}

	// Создаём медиа
	media := &model.AdMedia{AdID: 1, FilePath: "path.jpg", FileName: "file.jpg"}
	repo.Create(context.Background(), media)

	// Сохраняем ID
	id := media.ID.Hex()

	// Проверяем, что создалось
	result, _ := repo.FindByAdID(context.Background(), 1)
	assert.Len(t, result, 1)

	// Удаляем по ID
	err := repo.DeleteByID(context.Background(), id)
	assert.NoError(t, err)

	// Проверяем, что удалилось
	result, _ = repo.FindByAdID(context.Background(), 1)
	assert.Len(t, result, 0)
}

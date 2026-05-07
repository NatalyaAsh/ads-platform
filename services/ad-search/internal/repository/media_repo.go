package repository

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/NatalyaAsh/ads-platform/services/ad-search/internal/model"
)

type MediaRepository struct {
	collection *mongo.Collection
}

func NewMediaRepository(db *mongo.Database) *MediaRepository {
	return &MediaRepository{
		collection: db.Collection("ad_media"),
	}
}

// Create сохраняет метаданные файла
func (r *MediaRepository) Create(ctx context.Context, media *model.AdMedia) error {
	media.ID = primitive.NewObjectID()
	media.CreatedAt = time.Now()

	_, err := r.collection.InsertOne(ctx, media)
	if err != nil {
		return fmt.Errorf("failed to create media metadata: %w", err)
	}
	return nil
}

// FindByAdID возвращает все метаданные файлов объявления
func (r *MediaRepository) FindByAdID(ctx context.Context, adID uint) ([]model.AdMedia, error) {
	filter := bson.M{"ad_id": adID}
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to find media: %w", err)
	}
	defer cursor.Close(ctx)

	var media []model.AdMedia
	if err := cursor.All(ctx, &media); err != nil {
		return nil, fmt.Errorf("failed to decode media: %w", err)
	}
	return media, nil
}

// DeleteByAdID удаляет все метаданные файлов объявления
func (r *MediaRepository) DeleteByAdID(ctx context.Context, adID uint) error {
	filter := bson.M{"ad_id": adID}
	_, err := r.collection.DeleteMany(ctx, filter)
	if err != nil {
		return fmt.Errorf("failed to delete media metadata: %w", err)
	}
	return nil
}

// DeleteByID удаляет метаданные одного файла по ID
func (r *MediaRepository) DeleteByID(ctx context.Context, id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid media id: %w", err)
	}

	filter := bson.M{"_id": objID}
	_, err = r.collection.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("failed to delete media: %w", err)
	}
	return nil
}

package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AdMedia struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	AdID      uint               `bson:"ad_id"`
	FilePath  string             `bson:"file_path"`  // путь к файлу на диске
	FileName  string             `bson:"file_name"`  // оригинальное имя
	FileSize  int64              `bson:"file_size"`  // размер в байтах
	MimeType  string             `bson:"mime_type"`  // image/jpeg
	IsPrimary bool               `bson:"is_primary"` // главное фото
	CreatedAt time.Time          `bson:"created_at"`
}

package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AdMedia struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	AdID      uint               `bson:"ad_id"`
	URL       string             `bson:"url"`
	FileName  string             `bson:"file_name"`
	FileSize  int64              `bson:"file_size"`
	MimeType  string             `bson:"mime_type"`
	IsPrimary bool               `bson:"is_primary"`
	CreatedAt time.Time          `bson:"created_at"`
}

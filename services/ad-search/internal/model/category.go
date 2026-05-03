package model

import "time"

// Category - категория объявления
type Category struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"size:100;not null;uniqueIndex" json:"name"`
	Slug      string    `gorm:"size:100;not null;uniqueIndex" json:"slug"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (Category) TableName() string {
	return "categories"
}

package model

import (
	"time"
)

// AdStatus - статус объявления
type AdStatus string

const (
	AdStatusDraft     AdStatus = "draft"     // черновик
	AdStatusPending   AdStatus = "pending"   // на модерации
	AdStatusPublished AdStatus = "published" // опубликовано
	AdStatusRejected  AdStatus = "rejected"  // отклонено
	AdStatusDeleted   AdStatus = "deleted"   // удалено
)

// Ad - модель объявления
type Ad struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Title       string    `gorm:"size:200;not null;index" json:"title"`
	Description string    `gorm:"type:text" json:"description"`
	Price       float64   `gorm:"not null;index" json:"price"`
	UserID      uint      `gorm:"not null;index" json:"user_id"`
	CategoryID  uint      `gorm:"not null;index" json:"category_id"`
	Status      AdStatus  `gorm:"size:20;default:'draft';index" json:"status"`
	Views       int       `gorm:"default:0" json:"views"`
	CreatedAt   time.Time `gorm:"autoCreateTime;index" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	// Связи (не хранятся в БД, для GORM)
	Category Category `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
}

func (Ad) TableName() string {
	return "ads"
}

// IsPublished - проверка, опубликовано ли объявление
func (a *Ad) IsPublished() bool {
	return a.Status == AdStatusPublished
}

// CanBeEditedBy - может ли пользователь редактировать
func (a *Ad) CanBeEditedBy(userID uint, isAdmin bool) bool {
	return isAdmin || a.UserID == userID
}

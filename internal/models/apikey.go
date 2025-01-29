package models

import "time"

type ApiKey struct {
	UserID    int       `gorm:"primaryKey;not null"`
	ApiKey    string    `gorm:"type:uuid;primaryKey;not null"`
	Status    bool      `gorm:"default:true"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
	DeletedAt time.Time `gorm:"index"`
}

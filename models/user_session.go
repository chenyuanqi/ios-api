package models

import (
	"time"
)

// UserSession 用户会话模型
type UserSession struct {
	ID        uint       `json:"id" gorm:"primaryKey"`
	UserID    uint       `json:"user_id" gorm:"index"`
	Token     string     `json:"token" gorm:"uniqueIndex;size:255;not null"`
	ExpiredAt *time.Time `json:"expired_at"`
	CreatedAt time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
	User      User       `json:"-" gorm:"foreignKey:UserID"`
} 
package models

import (
	"time"
)

// OAuthAccount 第三方账号绑定模型
type OAuthAccount struct {
	ID             uint      `json:"id" gorm:"primaryKey"`
	UserID         uint      `json:"user_id" gorm:"index"`
	Provider       string    `json:"provider" gorm:"size:50;not null"`
	ProviderUserID string    `json:"provider_user_id" gorm:"size:255;not null"`
	CreatedAt      time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt      time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	User           User      `json:"-" gorm:"foreignKey:UserID"`
}

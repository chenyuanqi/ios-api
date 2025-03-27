package models

import (
	"time"
)

// User 用户模型
type User struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Email     string    `json:"email" gorm:"uniqueIndex;size:255;default:null"`
	Password  string    `json:"-" gorm:"size:255;default:null"` // 不返回给前端
	Nickname  string    `json:"nickname" gorm:"size:255;default:null"`
	Avatar    string    `json:"avatar" gorm:"size:255;default:null"`
	Signature string    `json:"signature" gorm:"type:text;default:null"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

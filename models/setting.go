package models

import (
	"time"
)

// Setting 设置模型，对应settings表
type Setting struct {
	Key       string    `json:"key" gorm:"primaryKey;size:64;comment:键，唯一标识"`
	Value     string    `json:"value" gorm:"type:text;comment:内容"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime;comment:记录创建时间"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime;comment:记录更新时间"`
}

// TableName 指定表名
func (Setting) TableName() string {
	return "settings"
}

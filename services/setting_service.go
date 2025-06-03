package services

import (
	"crypto/md5"
	"fmt"
	"ios-api/models"

	"gorm.io/gorm"
)

// SettingService 设置服务
type SettingService struct {
	DB   *gorm.DB // 通用数据库连接
	Salt string   // MD5校验盐值
}

// GetSetting 获取指定key的设置
func (s *SettingService) GetSetting(key string) (*models.Setting, error) {
	var setting models.Setting

	err := s.DB.Where("`key` = ?", key).First(&setting).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil // 返回nil表示未找到
		}
		return nil, fmt.Errorf("查询设置失败: %w", err)
	}

	return &setting, nil
}

// SetSetting 设置/更新指定key的值，需要MD5校验
func (s *SettingService) SetSetting(key, value, keyMD5 string) (*models.Setting, error) {
	// 验证key的MD5值
	if !s.validateKeyMD5(key, keyMD5) {
		return nil, fmt.Errorf("key的MD5校验失败")
	}

	// 验证key格式（只允许字母、数字、下划线、点号）
	if !s.validateKeyFormat(key) {
		return nil, fmt.Errorf("key格式不正确，只允许字母、数字、下划线、点号")
	}

	// 验证key长度
	if len(key) < 1 || len(key) > 64 {
		return nil, fmt.Errorf("key长度必须在1-64字符之间")
	}

	// 验证value长度（限制为10KB）
	if len(value) > 10240 {
		return nil, fmt.Errorf("value长度不能超过10KB")
	}

	var setting models.Setting

	// 尝试查找已存在的设置
	err := s.DB.Where("`key` = ?", key).First(&setting).Error
	if err == gorm.ErrRecordNotFound {
		// 不存在，创建新设置
		setting = models.Setting{
			Key:   key,
			Value: value,
		}
		err = s.DB.Create(&setting).Error
		if err != nil {
			return nil, fmt.Errorf("创建设置失败: %w", err)
		}
	} else if err != nil {
		return nil, fmt.Errorf("查询设置失败: %w", err)
	} else {
		// 存在，更新设置
		setting.Value = value
		err = s.DB.Save(&setting).Error
		if err != nil {
			return nil, fmt.Errorf("更新设置失败: %w", err)
		}
	}

	return &setting, nil
}

// validateKeyMD5 验证key的MD5值（使用配置的盐值）
func (s *SettingService) validateKeyMD5(key, keyMD5 string) bool {
	// 计算 key + salt 的MD5值
	expectedMD5 := fmt.Sprintf("%x", md5.Sum([]byte(key+s.Salt)))
	return expectedMD5 == keyMD5
}

// validateKeyFormat 验证key格式（只允许字母、数字、下划线、点号）
func (s *SettingService) validateKeyFormat(key string) bool {
	for _, char := range key {
		if !((char >= 'a' && char <= 'z') ||
			(char >= 'A' && char <= 'Z') ||
			(char >= '0' && char <= '9') ||
			char == '_' || char == '.') {
			return false
		}
	}
	return true
}

package services

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"ios-api/models"
	"path/filepath"

	"github.com/syndtr/goleveldb/leveldb"
	"gorm.io/gorm"
)

// SettingService 设置服务
type SettingService struct {
	DB    *gorm.DB    // 通用数据库连接
	Salt  string      // MD5校验盐值
	Cache *leveldb.DB // LevelDB缓存
}

// NewSettingService 创建新的设置服务实例
func NewSettingService(db *gorm.DB, salt string, cacheDir string) (*SettingService, error) {
	// 创建缓存目录路径
	cachePath := filepath.Join(cacheDir, "settings_cache")

	// 打开LevelDB数据库
	cache, err := leveldb.OpenFile(cachePath, nil)
	if err != nil {
		return nil, fmt.Errorf("无法打开LevelDB缓存: %w", err)
	}

	return &SettingService{
		DB:    db,
		Salt:  salt,
		Cache: cache,
	}, nil
}

// Close 关闭缓存连接
func (s *SettingService) Close() error {
	if s.Cache != nil {
		return s.Cache.Close()
	}
	return nil
}

// getCacheKey 生成缓存键
func (s *SettingService) getCacheKey(key string) string {
	return fmt.Sprintf("setting:%s", key)
}

// GetSetting 获取指定key的设置（支持缓存）
func (s *SettingService) GetSetting(key string) (*models.Setting, error) {
	cacheKey := s.getCacheKey(key)

	// 首先尝试从缓存读取
	if s.Cache != nil {
		cachedData, err := s.Cache.Get([]byte(cacheKey), nil)
		if err == nil {
			// 缓存命中，反序列化数据
			var setting models.Setting
			if err := json.Unmarshal(cachedData, &setting); err == nil {
				return &setting, nil
			}
			// 如果反序列化失败，删除无效缓存
			s.Cache.Delete([]byte(cacheKey), nil)
		}
	}

	// 缓存未命中，从数据库查询
	var setting models.Setting
	err := s.DB.Where("`key` = ?", key).First(&setting).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil // 返回nil表示未找到
		}
		return nil, fmt.Errorf("查询设置失败: %w", err)
	}

	// 将结果存入缓存
	if s.Cache != nil {
		if cachedData, err := json.Marshal(setting); err == nil {
			s.Cache.Put([]byte(cacheKey), cachedData, nil)
		}
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

	// 删除缓存
	if s.Cache != nil {
		cacheKey := s.getCacheKey(key)
		s.Cache.Delete([]byte(cacheKey), nil)
	}

	return &setting, nil
}

// ClearCache 清除指定key的缓存
func (s *SettingService) ClearCache(key string) error {
	if s.Cache != nil {
		cacheKey := s.getCacheKey(key)
		return s.Cache.Delete([]byte(cacheKey), nil)
	}
	return nil
}

// ClearAllCache 清除所有设置缓存
func (s *SettingService) ClearAllCache() error {
	if s.Cache == nil {
		return nil
	}

	iter := s.Cache.NewIterator(nil, nil)
	defer iter.Release()

	for iter.Next() {
		key := iter.Key()
		// 只删除以"setting:"开头的缓存键
		if len(key) > 8 && string(key[:8]) == "setting:" {
			if err := s.Cache.Delete(key, nil); err != nil {
				return fmt.Errorf("清除缓存失败: %w", err)
			}
		}
	}

	return iter.Error()
}

// GetCacheStats 获取缓存统计信息
func (s *SettingService) GetCacheStats() map[string]interface{} {
	stats := make(map[string]interface{})

	if s.Cache == nil {
		stats["cache_enabled"] = false
		return stats
	}

	stats["cache_enabled"] = true

	// 统计缓存中的设置数量
	count := 0
	iter := s.Cache.NewIterator(nil, nil)
	defer iter.Release()

	for iter.Next() {
		key := iter.Key()
		if len(key) > 8 && string(key[:8]) == "setting:" {
			count++
		}
	}

	stats["cached_settings_count"] = count
	return stats
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

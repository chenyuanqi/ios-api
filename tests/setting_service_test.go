package tests

import (
	"crypto/md5"
	"fmt"
	"ios-api/models"
	"ios-api/services"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func setupTestSettingDB() *gorm.DB {
	// 使用测试数据库进行测试
	dsn := "root:@tcp(localhost:3306)/test_yuanqi_general?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		// 如果连接失败，跳过测试
		return nil
	}

	// 自动迁移
	db.AutoMigrate(&models.Setting{})

	return db
}

func TestSettingService_SetSetting_WithSalt(t *testing.T) {
	db := setupTestSettingDB()
	if db == nil {
		t.Skip("跳过测试：无法连接到测试数据库")
		return
	}

	// 清理测试数据
	db.Where("1 = 1").Delete(&models.Setting{})

	// 测试盐值
	testSalt := "test_salt_12345"
	service := &services.SettingService{
		DB:   db,
		Salt: testSalt,
	}

	// 测试MD5校验（使用测试盐值）
	key := "test.setting"
	value := "test value"
	correctMD5 := fmt.Sprintf("%x", md5.Sum([]byte(key+testSalt)))
	wrongMD5 := "wrong_md5"

	// 测试错误的MD5
	_, err := service.SetSetting(key, value, wrongMD5)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "MD5校验失败")

	// 测试正确的MD5 - 创建新设置
	setting, err := service.SetSetting(key, value, correctMD5)
	assert.NoError(t, err)
	assert.NotNil(t, setting)
	assert.Equal(t, key, setting.Key)
	assert.Equal(t, value, setting.Value)

	// 测试更新现有设置
	newValue := "updated value"
	setting, err = service.SetSetting(key, newValue, correctMD5)
	assert.NoError(t, err)
	assert.NotNil(t, setting)
	assert.Equal(t, key, setting.Key)
	assert.Equal(t, newValue, setting.Value)
}

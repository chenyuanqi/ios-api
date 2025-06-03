package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config 应用配置
type Config struct {
	DBHost     string
	DBUser     string
	DBPassword string
	DBPort     int
	DBName     string

	// 通用数据库配置（用于settings表）
	GeneralDBHost     string
	GeneralDBUser     string
	GeneralDBPassword string
	GeneralDBPort     int
	GeneralDBName     string

	JWTSecret string
	AppPort   int

	// 设置管理配置
	SettingSalt string
	CacheDir    string // LevelDB缓存目录

	// 微信登录配置
	WechatAppID     string
	WechatAppSecret string

	// 苹果登录配置
	AppleTeamID     string
	AppleKeyID      string
	ApplePrivateKey string
	AppleBundleID   string
}

// LoadConfig 从环境变量加载配置
func LoadConfig() (*Config, error) {
	// 加载 .env 文件
	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("无法加载 .env 文件: %w", err)
	}

	// 从环境变量获取配置值
	dbPort, _ := strconv.Atoi(getEnv("DB_PORT", "3306"))
	generalDBPort, _ := strconv.Atoi(getEnv("GENERAL_DB_PORT", "3306"))
	appPort, _ := strconv.Atoi(getEnv("APP_PORT", "8080"))

	return &Config{
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBUser:     getEnv("DB_USER", "root"),
		DBPassword: getEnv("DB_PASSWORD", ""),
		DBPort:     dbPort,
		DBName:     getEnv("DB_NAME", "yuanqi_ios"),

		// 通用数据库配置
		GeneralDBHost:     getEnv("GENERAL_DB_HOST", getEnv("DB_HOST", "localhost")),
		GeneralDBUser:     getEnv("GENERAL_DB_USER", getEnv("DB_USER", "root")),
		GeneralDBPassword: getEnv("GENERAL_DB_PASSWORD", getEnv("DB_PASSWORD", "")),
		GeneralDBPort:     generalDBPort,
		GeneralDBName:     getEnv("GENERAL_DB_NAME", "yuanqi_general"),

		JWTSecret: getEnv("JWT_SECRET", "default_jwt_secret"),
		AppPort:   appPort,

		// 设置管理配置
		SettingSalt: getEnv("SETTING_SALT", "default_setting_salt"),
		CacheDir:    getEnv("CACHE_DIR", "./cache"), // 默认缓存目录

		// 微信登录配置
		WechatAppID:     getEnv("WECHAT_APP_ID", ""),
		WechatAppSecret: getEnv("WECHAT_APP_SECRET", ""),

		// 苹果登录配置
		AppleTeamID:     getEnv("APPLE_TEAM_ID", ""),
		AppleKeyID:      getEnv("APPLE_KEY_ID", ""),
		ApplePrivateKey: getEnv("APPLE_PRIVATE_KEY", ""),
		AppleBundleID:   getEnv("APPLE_BUNDLE_ID", ""),
	}, nil
}

// GetDSN 获取数据库连接字符串
func (c *Config) GetDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		c.DBUser, c.DBPassword, c.DBHost, c.DBPort, c.DBName)
}

// GetGeneralDSN 获取通用数据库连接字符串
func (c *Config) GetGeneralDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		c.GeneralDBUser, c.GeneralDBPassword, c.GeneralDBHost, c.GeneralDBPort, c.GeneralDBName)
}

// 获取环境变量，如果不存在则返回默认值
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

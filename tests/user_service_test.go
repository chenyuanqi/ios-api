package tests

import (
	"fmt"
	"log"
	"testing"
	"time"

	"ios-api/config"
	"ios-api/models"
	"ios-api/services"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// 获取测试配置（如果 .env 加载失败，则使用硬编码的测试配置）
func getTestConfig() *config.Config {
	cfg, err := config.LoadConfig()
	if err != nil {
		// 使用硬编码的测试配置
		return &config.Config{
			DBHost:     "132.232.48.66",
			DBUser:     "root",
			DBPassword: "6368656e7975616e7169",
			DBPort:     3306,
			DBName:     "yuanqi_ios_test",
			JWTSecret:  "test_jwt_secret",
		}
	}
	// 确保使用测试数据库
	cfg.DBName = "yuanqi_ios_test"
	return cfg
}

// 创建测试数据库
func setupTestDB() *gorm.DB {
	// 加载测试配置
	cfg := getTestConfig()

	// 使用配置中的数据库信息，但使用测试数据库名
	dsn := getTestDSN(cfg)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("无法连接到测试数据库: %v", err)
	}

	// 禁用外键检查
	db.Exec("SET FOREIGN_KEY_CHECKS = 0")

	// 清空测试数据
	db.Exec("DROP TABLE IF EXISTS user_sessions")
	db.Exec("DROP TABLE IF EXISTS oauth_accounts")
	db.Exec("DROP TABLE IF EXISTS users")

	// 启用外键检查
	db.Exec("SET FOREIGN_KEY_CHECKS = 1")

	// 迁移表结构
	err = db.AutoMigrate(&models.User{}, &models.OAuthAccount{}, &models.UserSession{})
	if err != nil {
		log.Fatalf("迁移表结构失败: %v", err)
	}

	return db
}

// 获取测试数据库DSN
func getTestDSN(cfg *config.Config) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)
}

// 创建测试用户
func createTestUser(db *gorm.DB) (*models.User, error) {
	// 生成随机邮箱，避免冲突
	randomEmail := fmt.Sprintf("test%d@example.com", time.Now().UnixNano())
	
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("testpassword"), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Email:     randomEmail,
		Password:  string(hashedPassword),
		Nickname:  "测试用户",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := db.Create(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

// 测试用户注册
func TestUserRegister(t *testing.T) {
	// 加载测试配置
	cfg := getTestConfig()

	db := setupTestDB()
	userService := &services.UserService{
		DB:        db,
		JWTSecret: cfg.JWTSecret,
	}

	// 测试正常注册
	params := services.RegisterParams{
		Email:    "new@example.com",
		Password: "password123",
		Nickname: "新用户",
	}

	user, token, err := userService.Register(params)
	if err != nil {
		t.Errorf("注册失败: %v", err)
		return
	}
	if user == nil {
		t.Error("用户不应为nil")
		return
	}
	if token == "" {
		t.Error("token不应为空")
		return
	}
	if user.Email != params.Email {
		t.Errorf("邮箱不匹配，期望 %s，实际 %s", params.Email, user.Email)
	}
	if user.Nickname != params.Nickname {
		t.Errorf("昵称不匹配，期望 %s，实际 %s", params.Nickname, user.Nickname)
	}

	// 测试重复邮箱注册
	_, _, err = userService.Register(params)
	if err != services.ErrEmailExists {
		t.Errorf("应返回邮箱已存在错误，实际返回 %v", err)
	}
}

// 测试用户登录
func TestUserLogin(t *testing.T) {
	// 加载测试配置
	cfg := getTestConfig()

	db := setupTestDB()
	userService := &services.UserService{
		DB:        db,
		JWTSecret: cfg.JWTSecret,
	}

	// 创建测试用户
	testUser, err := createTestUser(db)
	if err != nil {
		t.Errorf("创建测试用户失败: %v", err)
		return
	}

	// 测试正常登录
	params := services.LoginParams{
		Email:    testUser.Email,
		Password: "testpassword",
	}

	user, token, err := userService.Login(params)
	if err != nil {
		t.Errorf("登录失败: %v", err)
		return
	}
	if user == nil {
		t.Error("用户不应为nil")
		return
	}
	if token == "" {
		t.Error("token不应为空")
		return
	}
	if user.ID != testUser.ID {
		t.Errorf("用户ID不匹配，期望 %d，实际 %d", testUser.ID, user.ID)
	}

	// 测试错误密码
	params.Password = "wrongpassword"
	_, _, err = userService.Login(params)
	if err != services.ErrInvalidPassword {
		t.Errorf("应返回密码错误，实际返回 %v", err)
	}

	// 测试不存在的用户
	params.Email = "nonexistent@example.com"
	_, _, err = userService.Login(params)
	if err != services.ErrUserNotFound {
		t.Errorf("应返回用户不存在错误，实际返回 %v", err)
	}
}

// 测试第三方登录
func TestOAuthLogin(t *testing.T) {
	// 加载测试配置
	cfg := getTestConfig()

	db := setupTestDB()
	userService := &services.UserService{
		DB:        db,
		JWTSecret: cfg.JWTSecret,
	}

	// 测试第一次第三方登录
	params := services.OAuthLoginParams{
		Provider:       "wechat",
		ProviderUserID: "wechat123",
		Nickname:       "微信用户",
		Avatar:         "https://example.com/avatar.jpg",
	}

	user, token, err := userService.OAuthLogin(params)
	if err != nil {
		t.Errorf("第三方登录失败: %v", err)
		return
	}
	if user == nil {
		t.Error("用户不应为nil")
		return
	}
	if token == "" {
		t.Error("token不应为空")
		return
	}
	if user.Nickname != params.Nickname {
		t.Errorf("昵称不匹配，期望 %s，实际 %s", params.Nickname, user.Nickname)
	}
	if user.Avatar != params.Avatar {
		t.Errorf("头像不匹配，期望 %s，实际 %s", params.Avatar, user.Avatar)
	}

	// 测试再次使用相同第三方账号登录
	user2, token2, err := userService.OAuthLogin(params)
	if err != nil {
		t.Errorf("第二次第三方登录失败: %v", err)
		return
	}
	if user2 == nil {
		t.Error("用户不应为nil")
		return
	}
	if token2 == "" {
		t.Error("token不应为空")
		return
	}
	if user.ID != user2.ID {
		t.Errorf("用户ID不匹配，期望 %d，实际 %d", user.ID, user2.ID)
	}
}

// 测试退出登录
func TestLogout(t *testing.T) {
	// 加载测试配置
	cfg := getTestConfig()

	db := setupTestDB()
	userService := &services.UserService{
		DB:        db,
		JWTSecret: cfg.JWTSecret,
	}

	// 创建测试用户
	testUser, err := createTestUser(db)
	if err != nil {
		t.Errorf("创建测试用户失败: %v", err)
		return
	}

	// 生成token
	token, err := userService.GenerateToken(testUser.ID)
	if err != nil {
		t.Errorf("生成token失败: %v", err)
		return
	}
	if token == "" {
		t.Errorf("token不应为空")
		return
	}

	// 测试退出登录
	err = userService.Logout(token)
	if err != nil {
		t.Errorf("退出登录失败: %v", err)
		return
	}

	// 测试退出已退出的登录
	err = userService.Logout(token)
	if err != services.ErrSessionNotFound {
		t.Errorf("应返回会话不存在错误，实际返回 %v", err)
	}
}

// 测试获取用户信息
func TestGetUserInfo(t *testing.T) {
	// 加载测试配置
	cfg := getTestConfig()

	db := setupTestDB()
	userService := &services.UserService{
		DB:        db,
		JWTSecret: cfg.JWTSecret,
	}

	// 创建测试用户
	testUser, err := createTestUser(db)
	if err != nil {
		t.Errorf("创建测试用户失败: %v", err)
		return
	}

	// 测试获取存在的用户信息
	user, err := userService.GetUserByID(testUser.ID)
	if err != nil {
		t.Errorf("获取用户信息失败: %v", err)
		return
	}
	if user == nil {
		t.Errorf("用户不应为nil")
		return
	}
	if user.ID != testUser.ID {
		t.Errorf("用户ID不匹配，期望 %d，实际 %d", testUser.ID, user.ID)
	}
	if user.Email != testUser.Email {
		t.Errorf("用户邮箱不匹配，期望 %s，实际 %s", testUser.Email, user.Email)
	}

	// 测试获取不存在的用户信息
	_, err = userService.GetUserByID(999)
	if err != services.ErrUserNotFound {
		t.Errorf("应返回用户不存在错误，实际返回 %v", err)
	}
}

// 测试更新用户信息
func TestUpdateUser(t *testing.T) {
	// 加载测试配置
	cfg := getTestConfig()

	db := setupTestDB()
	userService := &services.UserService{
		DB:        db,
		JWTSecret: cfg.JWTSecret,
	}

	// 创建测试用户
	testUser, err := createTestUser(db)
	if err != nil {
		t.Errorf("创建测试用户失败: %v", err)
		return
	}

	// 测试更新用户信息
	params := services.UpdateUserParams{
		Nickname:  "新昵称",
		Avatar:    "https://example.com/new-avatar.jpg",
		Signature: "新个性签名",
	}

	user, err := userService.UpdateUser(testUser.ID, params)
	if err != nil {
		t.Errorf("更新用户信息失败: %v", err)
		return
	}
	if user == nil {
		t.Errorf("用户不应为nil")
		return
	}
	if user.Nickname != params.Nickname {
		t.Errorf("用户昵称不匹配，期望 %s，实际 %s", params.Nickname, user.Nickname)
	}
	if user.Avatar != params.Avatar {
		t.Errorf("用户头像不匹配，期望 %s，实际 %s", params.Avatar, user.Avatar)
	}
	if user.Signature != params.Signature {
		t.Errorf("用户签名不匹配，期望 %s，实际 %s", params.Signature, user.Signature)
	}

	// 测试更新不存在的用户
	_, err = userService.UpdateUser(999, params)
	if err != services.ErrUserNotFound {
		t.Errorf("应返回用户不存在错误，实际返回 %v", err)
	}
}

// 测试验证Token
func TestVerifyToken(t *testing.T) {
	// 加载测试配置
	cfg := getTestConfig()

	db := setupTestDB()
	userService := &services.UserService{
		DB:        db,
		JWTSecret: cfg.JWTSecret,
	}

	// 创建测试用户
	testUser, err := createTestUser(db)
	if err != nil {
		t.Errorf("创建测试用户失败: %v", err)
		return
	}

	// 生成token
	token, err := userService.GenerateToken(testUser.ID)
	if err != nil {
		t.Errorf("生成token失败: %v", err)
		return
	}
	if token == "" {
		t.Errorf("token不应为空")
		return
	}

	// 测试验证有效token
	userID, err := userService.VerifyToken(token)
	if err != nil {
		t.Errorf("验证token失败: %v", err)
		return
	}
	if userID != testUser.ID {
		t.Errorf("用户ID不匹配，期望 %d，实际 %d", testUser.ID, userID)
	}

	// 测试验证无效token
	_, err = userService.VerifyToken("invalid_token")
	if err == nil {
		t.Errorf("应返回错误，实际没有返回错误")
	}
}

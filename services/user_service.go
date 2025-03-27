package services

import (
	"errors"
	"time"

	"ios-api/config"
	"ios-api/models"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// UserService 用户服务
type UserService struct {
	DB        *gorm.DB
	JWTSecret string
	Config    *config.Config
}

// 用户注册参数
type RegisterParams struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Nickname string `json:"nickname"`
}

// 用户登录参数
type LoginParams struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// 第三方登录参数
type OAuthLoginParams struct {
	Provider       string `json:"provider" binding:"required"`
	ProviderUserID string `json:"provider_user_id" binding:"required"`
	Nickname       string `json:"nickname"`
	Avatar         string `json:"avatar"`
	Email          string `json:"email"`
}

// 更新用户信息参数
type UpdateUserParams struct {
	Nickname  string `json:"nickname"`
	Avatar    string `json:"avatar"`
	Signature string `json:"signature"`
}

// 自定义错误
var (
	ErrUserNotFound    = errors.New("用户不存在")
	ErrEmailExists     = errors.New("邮箱已被注册")
	ErrInvalidPassword = errors.New("密码错误")
	ErrInvalidToken    = errors.New("无效的令牌")
	ErrTokenExpired    = errors.New("令牌已过期")
	ErrOAuthBound      = errors.New("第三方账号已绑定其他用户")
	ErrSessionNotFound = errors.New("会话不存在")
)

// 生成JWT
func (s *UserService) GenerateToken(userID uint) (string, error) {
	// 创建一个token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 24 * 7).Unix(), // 7天过期
	})

	// 签名并获得完整的编码后的字符串token
	tokenString, err := token.SignedString([]byte(s.JWTSecret))
	if err != nil {
		return "", err
	}

	// 保存会话
	expiredAt := time.Now().Add(time.Hour * 24 * 7)
	session := models.UserSession{
		UserID:    userID,
		Token:     tokenString,
		ExpiredAt: &expiredAt,
	}
	if err := s.DB.Create(&session).Error; err != nil {
		return "", err
	}

	return tokenString, nil
}

// Register 用户注册
func (s *UserService) Register(params RegisterParams) (*models.User, string, error) {
	// 检查邮箱是否已存在
	var existingUser models.User
	if err := s.DB.Where("email = ?", params.Email).First(&existingUser).Error; err == nil {
		return nil, "", ErrEmailExists
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, "", err
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, "", err
	}

	// 创建用户
	user := models.User{
		Email:    params.Email,
		Password: string(hashedPassword),
		Nickname: params.Nickname,
	}

	if err := s.DB.Create(&user).Error; err != nil {
		return nil, "", err
	}

	// 生成token
	token, err := s.GenerateToken(user.ID)
	if err != nil {
		return nil, "", err
	}

	return &user, token, nil
}

// Login 用户登录
func (s *UserService) Login(params LoginParams) (*models.User, string, error) {
	// 查找用户
	var user models.User
	if err := s.DB.Where("email = ?", params.Email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, "", ErrUserNotFound
		}
		return nil, "", err
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(params.Password)); err != nil {
		return nil, "", ErrInvalidPassword
	}

	// 生成token
	token, err := s.GenerateToken(user.ID)
	if err != nil {
		return nil, "", err
	}

	return &user, token, nil
}

// OAuthLogin 第三方登录
func (s *UserService) OAuthLogin(params OAuthLoginParams) (*models.User, string, error) {
	var oauthAccount models.OAuthAccount
	tx := s.DB.Begin()

	// 查找第三方账号是否已存在
	if err := tx.Where("provider = ? AND provider_user_id = ?", params.Provider, params.ProviderUserID).First(&oauthAccount).Error; err == nil {
		// 账号已存在，查找对应的用户
		var user models.User
		if err := tx.First(&user, oauthAccount.UserID).Error; err != nil {
			tx.Rollback()
			return nil, "", err
		}

		// 生成token
		token, err := s.GenerateToken(user.ID)
		if err != nil {
			tx.Rollback()
			return nil, "", err
		}

		tx.Commit()
		return &user, token, nil
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		tx.Rollback()
		return nil, "", err
	}

	// 账号不存在，创建新用户和账号绑定
	user := models.User{
		Nickname: params.Nickname,
		Avatar:   params.Avatar,
		Email:    params.Email,
	}

	if err := tx.Create(&user).Error; err != nil {
		tx.Rollback()
		return nil, "", err
	}

	// 创建第三方账号绑定
	oauthAccount = models.OAuthAccount{
		UserID:         user.ID,
		Provider:       params.Provider,
		ProviderUserID: params.ProviderUserID,
	}

	if err := tx.Create(&oauthAccount).Error; err != nil {
		tx.Rollback()
		return nil, "", err
	}

	// 生成token
	token, err := s.GenerateToken(user.ID)
	if err != nil {
		tx.Rollback()
		return nil, "", err
	}

	tx.Commit()
	return &user, token, nil
}

// Logout 用户退出登录
func (s *UserService) Logout(token string) error {
	// 删除用户会话
	result := s.DB.Where("token = ?", token).Delete(&models.UserSession{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrSessionNotFound
	}
	return nil
}

// GetUserByID 获取用户信息
func (s *UserService) GetUserByID(userID uint) (*models.User, error) {
	var user models.User
	if err := s.DB.First(&user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

// UpdateUser 更新用户信息
func (s *UserService) UpdateUser(userID uint, params UpdateUserParams) (*models.User, error) {
	var user models.User
	if err := s.DB.First(&user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	// 更新用户信息
	updates := map[string]interface{}{}
	if params.Nickname != "" {
		updates["nickname"] = params.Nickname
	}
	if params.Avatar != "" {
		updates["avatar"] = params.Avatar
	}
	if params.Signature != "" {
		updates["signature"] = params.Signature
	}

	if len(updates) > 0 {
		if err := s.DB.Model(&user).Updates(updates).Error; err != nil {
			return nil, err
		}
	}

	// 重新获取用户信息
	if err := s.DB.First(&user, userID).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

// VerifyToken 验证token
func (s *UserService) VerifyToken(tokenString string) (uint, error) {
	// 查找会话
	var session models.UserSession
	if err := s.DB.Where("token = ?", tokenString).First(&session).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, ErrInvalidToken
		}
		return 0, err
	}

	// 检查会话是否过期
	if session.ExpiredAt != nil && session.ExpiredAt.Before(time.Now()) {
		return 0, ErrTokenExpired
	}

	// 解析JWT
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// 验证算法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return []byte(s.JWTSecret), nil
	})

	if err != nil {
		return 0, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID := uint(claims["user_id"].(float64))
		return userID, nil
	}

	return 0, ErrInvalidToken
}

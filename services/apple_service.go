package services

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// AppleService 苹果服务
type AppleService struct {
	TeamID     string
	KeyID      string
	PrivateKey string
	BundleID   string
}

// AppleIdTokenPayload 苹果ID令牌载荷
type AppleIdTokenPayload struct {
	// 标准JWT字段
	Iss   string `json:"iss"`
	Aud   string `json:"aud"`
	Exp   int64  `json:"exp"`
	Iat   int64  `json:"iat"`
	Sub   string `json:"sub"`
	Nonce string `json:"nonce,omitempty"`
	// 苹果特有字段
	Email          string `json:"email,omitempty"`
	EmailVerified  string `json:"email_verified,omitempty"`
	IsPrivateEmail string `json:"is_private_email,omitempty"`
	RealUserStatus int    `json:"real_user_status,omitempty"`
	Name           struct {
		FirstName string `json:"firstName,omitempty"`
		LastName  string `json:"lastName,omitempty"`
	} `json:"name,omitempty"`
}

// AppleTokenResponse 苹果访问令牌响应
type AppleTokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	IdToken      string `json:"id_token"`
	Error        string `json:"error,omitempty"`
}

// 自定义错误
var (
	ErrAppleAuthFailed      = errors.New("苹果授权失败")
	ErrAppleCodeInvalid     = errors.New("无效的苹果授权码")
	ErrAppleTokenInvalid    = errors.New("无效的苹果令牌")
	ErrAppleServerError     = errors.New("苹果服务器错误")
	ErrApplePrivateKeyError = errors.New("苹果私钥解析错误")
)

// GenerateClientSecret 生成客户端密钥
func (s *AppleService) GenerateClientSecret() (string, error) {
	// 解析私钥
	var privateKey *ecdsa.PrivateKey

	// 检查私钥是否为文件路径
	if strings.HasPrefix(s.PrivateKey, "/") || strings.HasPrefix(s.PrivateKey, "./") {
		// 从文件读取私钥
		keyData, err := ioutil.ReadFile(s.PrivateKey)
		if err != nil {
			return "", fmt.Errorf("读取私钥文件失败: %w", err)
		}

		// 解析PEM格式
		block, _ := pem.Decode(keyData)
		if block == nil {
			return "", fmt.Errorf("%w: 无法解码PEM格式", ErrApplePrivateKeyError)
		}

		// 解析私钥
		key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
		if err != nil {
			return "", fmt.Errorf("%w: %v", ErrApplePrivateKeyError, err)
		}

		var ok bool
		privateKey, ok = key.(*ecdsa.PrivateKey)
		if !ok {
			return "", fmt.Errorf("%w: 不是ECDSA私钥", ErrApplePrivateKeyError)
		}
	} else {
		// 直接解析私钥内容
		// 假设内容已经是标准的PEM格式
		block, _ := pem.Decode([]byte(s.PrivateKey))
		if block == nil {
			return "", fmt.Errorf("%w: 无法解码PEM格式", ErrApplePrivateKeyError)
		}

		key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
		if err != nil {
			return "", fmt.Errorf("%w: %v", ErrApplePrivateKeyError, err)
		}

		var ok bool
		privateKey, ok = key.(*ecdsa.PrivateKey)
		if !ok {
			return "", fmt.Errorf("%w: 不是ECDSA私钥", ErrApplePrivateKeyError)
		}
	}

	// 创建JWT
	now := time.Now()
	claims := jwt.MapClaims{
		"iss": s.TeamID,                       // 发行者是Team ID
		"iat": now.Unix(),                     // 发行时间
		"exp": now.Add(time.Hour * 24).Unix(), // 过期时间（24小时）
		"aud": "https://appleid.apple.com",    // 目标受众
		"sub": s.BundleID,                     // 主题（应用的Bundle ID）
	}

	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	token.Header["kid"] = s.KeyID

	// 签名token
	clientSecret, err := token.SignedString(privateKey)
	if err != nil {
		return "", fmt.Errorf("签名客户端密钥失败: %w", err)
	}

	return clientSecret, nil
}

// ValidateIdToken 验证苹果ID令牌
func (s *AppleService) ValidateIdToken(idToken string) (*AppleIdTokenPayload, error) {
	// 解析JWT令牌而不验证签名
	parts := strings.Split(idToken, ".")
	if len(parts) != 3 {
		return nil, ErrAppleTokenInvalid
	}

	// 解码JWT载荷
	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrAppleTokenInvalid, err)
	}

	// 解析为JSON
	var tokenPayload AppleIdTokenPayload
	if err := json.Unmarshal(payload, &tokenPayload); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrAppleTokenInvalid, err)
	}

	// 验证token是否过期
	if tokenPayload.Exp < time.Now().Unix() {
		return nil, fmt.Errorf("%w: 令牌已过期", ErrAppleTokenInvalid)
	}

	// 验证发行者
	if tokenPayload.Iss != "https://appleid.apple.com" {
		return nil, fmt.Errorf("%w: 发行者无效", ErrAppleTokenInvalid)
	}

	return &tokenPayload, nil
}

// ExchangeAuthCodeForToken 使用授权码交换访问令牌
func (s *AppleService) ExchangeAuthCodeForToken(code string) (*AppleTokenResponse, error) {
	// 生成客户端密钥
	clientSecret, err := s.GenerateClientSecret()
	if err != nil {
		return nil, err
	}

	// 创建请求数据
	data := url.Values{}
	data.Set("client_id", s.BundleID)
	data.Set("client_secret", clientSecret)
	data.Set("code", code)
	data.Set("grant_type", "authorization_code")

	// 发送POST请求
	resp, err := http.PostForm("https://appleid.apple.com/auth/token", data)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrAppleServerError, err)
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrAppleServerError, err)
	}

	// 解析响应
	var tokenResp AppleTokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrAppleServerError, err)
	}

	// 检查错误
	if tokenResp.Error != "" {
		return nil, fmt.Errorf("%w: %s", ErrAppleAuthFailed, tokenResp.Error)
	}

	return &tokenResp, nil
}

// HandleCallback 处理苹果授权回调
func (s *AppleService) HandleCallback(code, idToken, name, email string) (*OAuthLoginParams, error) {
	var tokenPayload *AppleIdTokenPayload

	// 如果提供了授权码，则交换访问令牌
	if code != "" {
		tokenResp, err := s.ExchangeAuthCodeForToken(code)
		if err != nil {
			return nil, err
		}

		// 验证ID令牌
		tokenPayload, err = s.ValidateIdToken(tokenResp.IdToken)
		if err != nil {
			return nil, err
		}
	} else if idToken != "" {
		// 直接验证ID令牌
		var err error
		tokenPayload, err = s.ValidateIdToken(idToken)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, ErrAppleCodeInvalid
	}

	// 构建用户昵称
	nickname := ""
	if name != "" {
		// 如果前端提供了用户名，使用前端提供的
		nickname = name
	} else if tokenPayload.Name.FirstName != "" || tokenPayload.Name.LastName != "" {
		// 否则尝试从令牌中获取
		nickname = tokenPayload.Name.FirstName + " " + tokenPayload.Name.LastName
		nickname = strings.TrimSpace(nickname)
	} else {
		// 如果没有名字，使用"Apple User"
		nickname = "Apple User"
	}

	// 构建OAuth登录参数
	params := &OAuthLoginParams{
		Provider:       "apple",
		ProviderUserID: tokenPayload.Sub, // 使用Sub作为用户标识
		Nickname:       nickname,
		Email:          email, // 如果前端提供了邮箱（首次登录时），则使用前端提供的
	}

	// 如果没有提供邮箱，但令牌中有邮箱，则使用令牌中的
	if params.Email == "" && tokenPayload.Email != "" {
		params.Email = tokenPayload.Email
	}

	return params, nil
}

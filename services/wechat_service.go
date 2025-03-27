package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

// WechatService 微信服务
type WechatService struct {
	AppID       string
	AppSecret   string
	RedirectURI string
}

// WechatAccessTokenResponse 微信访问令牌响应
type WechatAccessTokenResponse struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	OpenID       string `json:"openid"`
	Scope        string `json:"scope"`
	ErrCode      int    `json:"errcode"`
	ErrMsg       string `json:"errmsg"`
}

// WechatUserInfoResponse 微信用户信息响应
type WechatUserInfoResponse struct {
	OpenID     string   `json:"openid"`
	Nickname   string   `json:"nickname"`
	Sex        int      `json:"sex"`
	Province   string   `json:"province"`
	City       string   `json:"city"`
	Country    string   `json:"country"`
	HeadImgURL string   `json:"headimgurl"`
	Privilege  []string `json:"privilege"`
	UnionID    string   `json:"unionid"`
	ErrCode    int      `json:"errcode"`
	ErrMsg     string   `json:"errmsg"`
}

// 自定义错误
var (
	ErrWechatAuthFailed     = errors.New("微信授权失败")
	ErrWechatCodeInvalid    = errors.New("无效的微信授权码")
	ErrWechatServerError    = errors.New("微信服务器错误")
	ErrWechatUserInfoFailed = errors.New("获取微信用户信息失败")
)

// GetAuthURL 获取微信授权链接
func (s *WechatService) GetAuthURL(state string) string {
	authURL := fmt.Sprintf(
		"https://open.weixin.qq.com/connect/oauth2/authorize?appid=%s&redirect_uri=%s&response_type=code&scope=snsapi_userinfo&state=%s#wechat_redirect",
		s.AppID,
		url.QueryEscape(s.RedirectURI),
		state,
	)
	return authURL
}

// GetAccessToken 通过授权码获取访问令牌
func (s *WechatService) GetAccessToken(code string) (*WechatAccessTokenResponse, error) {
	// 构建接口URL
	tokenURL := fmt.Sprintf(
		"https://api.weixin.qq.com/sns/oauth2/access_token?appid=%s&secret=%s&code=%s&grant_type=authorization_code",
		s.AppID,
		s.AppSecret,
		code,
	)

	// 发送请求
	resp, err := http.Get(tokenURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// 解析JSON
	var tokenResp WechatAccessTokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return nil, err
	}

	// 检查错误
	if tokenResp.ErrCode != 0 {
		return nil, fmt.Errorf("%w: %d %s", ErrWechatAuthFailed, tokenResp.ErrCode, tokenResp.ErrMsg)
	}

	return &tokenResp, nil
}

// GetUserInfo 获取微信用户信息
func (s *WechatService) GetUserInfo(accessToken, openID string) (*WechatUserInfoResponse, error) {
	// 构建接口URL
	userInfoURL := fmt.Sprintf(
		"https://api.weixin.qq.com/sns/userinfo?access_token=%s&openid=%s&lang=zh_CN",
		accessToken,
		openID,
	)

	// 发送请求
	resp, err := http.Get(userInfoURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// 解析JSON
	var userInfo WechatUserInfoResponse
	if err := json.Unmarshal(body, &userInfo); err != nil {
		return nil, err
	}

	// 检查错误
	if userInfo.ErrCode != 0 {
		return nil, fmt.Errorf("%w: %d %s", ErrWechatUserInfoFailed, userInfo.ErrCode, userInfo.ErrMsg)
	}

	return &userInfo, nil
}

// HandleCallback 处理微信授权回调
func (s *WechatService) HandleCallback(code string) (*OAuthLoginParams, error) {
	if code == "" {
		return nil, ErrWechatCodeInvalid
	}

	// 获取访问令牌
	tokenResp, err := s.GetAccessToken(code)
	if err != nil {
		return nil, err
	}

	// 获取用户信息
	userInfo, err := s.GetUserInfo(tokenResp.AccessToken, tokenResp.OpenID)
	if err != nil {
		return nil, err
	}

	// 构建OAuth登录参数
	params := &OAuthLoginParams{
		Provider:       "wechat",
		ProviderUserID: userInfo.OpenID, // 使用OpenID作为用户标识
		Nickname:       userInfo.Nickname,
		Avatar:         userInfo.HeadImgURL,
	}

	return params, nil
}

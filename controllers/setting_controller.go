package controllers

import (
	"ios-api/services"
	"ios-api/utils"

	"github.com/gin-gonic/gin"
)

// SettingController 设置控制器
type SettingController struct {
	SettingService *services.SettingService
}

// GetSetting 获取指定key的设置
// GET /api/v1/settings/:key
func (sc *SettingController) GetSetting(c *gin.Context) {
	key := c.Param("key")

	// 验证key参数
	if key == "" {
		utils.ParamError(c, "key参数不能为空")
		return
	}

	// 获取设置
	setting, err := sc.SettingService.GetSetting(key)
	if err != nil {
		utils.ServerError(c, "获取设置失败: "+err.Error())
		return
	}

	if setting == nil {
		utils.NotFound(c, "设置不存在")
		return
	}

	utils.Success(c, "获取设置成功", setting)
}

// SetSetting 设置/更新指定key的值
// PUT /api/v1/settings/:key
func (sc *SettingController) SetSetting(c *gin.Context) {
	key := c.Param("key")

	// 验证key参数
	if key == "" {
		utils.ParamError(c, "key参数不能为空")
		return
	}

	// 解析请求体
	var req struct {
		Value  string `json:"value" binding:"required"`
		KeyMD5 string `json:"key_md5" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ParamError(c, "请求参数错误: "+err.Error())
		return
	}

	// 设置/更新设置
	setting, err := sc.SettingService.SetSetting(key, req.Value, req.KeyMD5)
	if err != nil {
		// 根据错误类型返回不同的响应
		if err.Error() == "key的MD5校验失败" {
			utils.Unauthorized(c, "MD5校验失败，无权限操作此设置")
			return
		}
		utils.ParamError(c, err.Error())
		return
	}

	utils.Success(c, "设置保存成功", setting)
}

// ClearCache 清除指定key的缓存
// DELETE /api/v1/settings/:key/cache
func (sc *SettingController) ClearCache(c *gin.Context) {
	key := c.Param("key")

	// 验证key参数
	if key == "" {
		utils.ParamError(c, "key参数不能为空")
		return
	}

	// 清除缓存
	err := sc.SettingService.ClearCache(key)
	if err != nil {
		utils.ServerError(c, "清除缓存失败: "+err.Error())
		return
	}

	utils.Success(c, "缓存清除成功", nil)
}

// ClearAllCache 清除所有设置缓存
// DELETE /api/v1/settings/cache
func (sc *SettingController) ClearAllCache(c *gin.Context) {
	// 清除所有缓存
	err := sc.SettingService.ClearAllCache()
	if err != nil {
		utils.ServerError(c, "清除所有缓存失败: "+err.Error())
		return
	}

	utils.Success(c, "所有缓存清除成功", nil)
}

// GetCacheStats 获取缓存统计信息
// GET /api/v1/settings/cache/stats
func (sc *SettingController) GetCacheStats(c *gin.Context) {
	stats := sc.SettingService.GetCacheStats()
	utils.Success(c, "获取缓存统计成功", stats)
}

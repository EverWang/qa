package controllers

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"qaminiprogram/config"
	"qaminiprogram/middleware"
	"qaminiprogram/models"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// WechatLoginRequest 微信登录请求
type WechatLoginRequest struct {
	Code string `json:"code" binding:"required"`
}

// PasswordLoginRequest 密码登录请求
type PasswordLoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Type     string `json:"type" binding:"required"`
}

// WechatLoginResponse 微信登录响应
type WechatLoginResponse struct {
	Token    string      `json:"token"`
	User     models.User `json:"user"`
	ExpireAt int64       `json:"expire_at"`
}

// WechatSessionResponse 微信会话响应
type WechatSessionResponse struct {
	OpenID     string `json:"openid"`
	SessionKey string `json:"session_key"`
	UnionID    string `json:"unionid,omitempty"`
	ErrCode    int    `json:"errcode,omitempty"`
	ErrMsg     string `json:"errmsg,omitempty"`
}

// RefreshTokenRequest 刷新令牌请求
type RefreshTokenRequest struct {
	Token string `json:"token" binding:"required"`
}

// AdminLoginRequest 管理员登录请求
type AdminLoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// GuestLoginResponse 游客登录响应
type GuestLoginResponse struct {
	Token    string      `json:"token"`
	User     models.User `json:"user"`
	ExpireAt int64       `json:"expire_at"`
}

// UnifiedLogin 统一登录接口
func UnifiedLogin(c *gin.Context) {
	// 先尝试解析为密码登录请求
	var passwordReq PasswordLoginRequest
	if err := c.ShouldBindJSON(&passwordReq); err == nil && passwordReq.Type == "password" {
		// 处理密码登录
		PasswordLogin(c, passwordReq)
		return
	}

	// 尝试解析为微信登录请求
	var wechatReq WechatLoginRequest
	if err := c.ShouldBindJSON(&wechatReq); err == nil {
		// 处理微信登录
		WechatLogin(c, wechatReq)
		return
	}

	ErrorResponse(c, http.StatusBadRequest, "参数错误")
}

// PasswordLogin 密码登录
func PasswordLogin(c *gin.Context, req PasswordLoginRequest) {
	db := config.GetDB()
	var user models.User

	// 查找用户
	if err := db.Where("username = ?", req.Username).First(&user).Error; err != nil {
		ErrorResponse(c, http.StatusUnauthorized, "用户名或密码错误")
		return
	}

	// 验证密码
	if user.Password == "" {
		ErrorResponse(c, http.StatusUnauthorized, "用户未设置密码")
		return
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		ErrorResponse(c, http.StatusUnauthorized, "用户名或密码错误")
		return
	}

	// 生成JWT令牌
	var openIDStr string
	if user.OpenID != nil {
		openIDStr = *user.OpenID
	}
	token, err := middleware.GenerateJWT(user.ID, openIDStr, user.Role)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "生成令牌失败")
		return
	}

	// 返回登录结果
	SuccessResponse(c, WechatLoginResponse{
		Token:    token,
		User:     user,
		ExpireAt: time.Now().Add(24 * time.Hour).Unix(),
	})
}

// WechatLogin 微信登录
func WechatLogin(c *gin.Context, req WechatLoginRequest) {

	// 调用微信API获取openid
	openID, err := getWechatOpenID(req.Code)
	if err != nil {
		ErrorResponse(c, http.StatusBadRequest, "微信登录失败："+err.Error())
		return
	}

	db := config.GetDB()
	var user models.User

	// 查找或创建用户
	result := db.Where("open_id = ?", openID).First(&user)
	if result.Error != nil {
		// 用户不存在，创建新用户
		user = models.User{
			OpenID:   &openID,
			Nickname: "微信用户" + generateRandomString(6),
			Avatar:   "",
			Role:     "user",
		}
		if err := db.Create(&user).Error; err != nil {
			ErrorResponse(c, http.StatusInternalServerError, "创建用户失败")
			return
		}
	}

	// 生成JWT令牌
	var openIDStr string
	if user.OpenID != nil {
		openIDStr = *user.OpenID
	}
	token, err := middleware.GenerateJWT(user.ID, openIDStr, user.Role)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "生成令牌失败")
		return
	}

	// 返回登录结果
	SuccessResponse(c, WechatLoginResponse{
		Token:    token,
		User:     user,
		ExpireAt: time.Now().Add(24 * time.Hour).Unix(),
	})
}

// RefreshToken 刷新令牌
func RefreshToken(c *gin.Context) {
	var req RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ErrorResponse(c, http.StatusBadRequest, "参数错误")
		return
	}

	// 解析旧令牌
	claims, err := middleware.ParseJWT(req.Token)
	if err != nil {
		ErrorResponse(c, http.StatusUnauthorized, "无效的令牌")
		return
	}

	// 验证用户是否存在
	db := config.GetDB()
	var user models.User
	if err := db.Where("id = ?", claims.UserID).First(&user).Error; err != nil {
		ErrorResponse(c, http.StatusUnauthorized, "用户不存在")
		return
	}

	// 生成新令牌
	var openIDStr string
	if user.OpenID != nil {
		openIDStr = *user.OpenID
	}
	newToken, err := middleware.GenerateJWT(user.ID, openIDStr, user.Role)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "生成令牌失败")
		return
	}

	// 返回新令牌
	SuccessResponse(c, gin.H{
		"token":     newToken,
		"expire_at": time.Now().Add(24 * time.Hour).Unix(),
	})
}

// AdminLoginResponse 管理员登录响应
type AdminLoginResponse struct {
	Token    string       `json:"token"`
	User     models.Admin `json:"user"`
	ExpireAt int64        `json:"expire_at"`
}

// AdminLogin 管理员登录
func AdminLogin(c *gin.Context) {
	var req AdminLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ErrorResponse(c, http.StatusBadRequest, "参数错误")
		return
	}

	db := config.GetDB()
	var admin models.Admin

	// 查找管理员
	if err := db.Where("username = ?", req.Username).First(&admin).Error; err != nil {
		ErrorResponse(c, http.StatusUnauthorized, "用户名或密码错误")
		return
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(admin.PasswordHash), []byte(req.Password)); err != nil {
		ErrorResponse(c, http.StatusUnauthorized, "用户名或密码错误")
		return
	}

	// 生成JWT令牌（使用管理员信息）
	token, err := middleware.GenerateJWT(admin.ID, "", admin.Role)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "生成令牌失败")
		return
	}

	// 记录操作日志
	LogOperation(c, "LOGIN", "ADMIN", fmt.Sprintf("管理员登录: %s", admin.Username))

	// 返回登录结果，使用user字段以匹配前端期望
	SuccessResponse(c, AdminLoginResponse{
		Token:    token,
		User:     admin,
		ExpireAt: time.Now().Add(24 * time.Hour).Unix(),
	})
}

// GetAdminProfile 获取管理员个人信息
func GetAdminProfile(c *gin.Context) {
	// 从JWT中获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		ErrorResponse(c, http.StatusUnauthorized, "未授权")
		return
	}

	db := config.GetDB()
	var admin models.Admin

	// 查找管理员信息
	if err := db.Where("id = ?", userID).First(&admin).Error; err != nil {
		ErrorResponse(c, http.StatusNotFound, "管理员不存在")
		return
	}

	// 返回管理员信息
	SuccessResponse(c, admin)
}

// AdminLogout 管理员登出
func AdminLogout(c *gin.Context) {
	// 这里可以实现token黑名单等逻辑
	// 目前只返回成功响应，实际的token清除由前端处理
	SuccessResponse(c, gin.H{
		"message": "登出成功",
	})
}

// GuestLoginRequest 游客登录请求
type GuestLoginRequest struct {
	DeviceID string `json:"device_id"` // 设备唯一标识
}

// GuestLogin 游客登录
func GuestLogin(c *gin.Context) {
	db := config.GetDB()
	
	// 尝试从请求体中获取设备ID
	var req GuestLoginRequest
	c.ShouldBindJSON(&req)
	
	// 如果没有提供设备ID，从请求头中获取或生成一个
	deviceID := req.DeviceID
	if deviceID == "" {
		// 尝试从请求头获取设备标识
		deviceID = c.GetHeader("X-Device-ID")
		if deviceID == "" {
			// 使用IP地址和User-Agent生成设备标识
			clientIP := c.ClientIP()
			userAgent := c.GetHeader("User-Agent")
			hash := sha256.Sum256([]byte(clientIP + userAgent))
			// 只取前16个字符的哈希值，确保用户名不超过50字符限制
			deviceID = fmt.Sprintf("%x", hash)[:16]
		}
	}
	
	var user models.User
	
	// 生成游客用户名，确保不超过50字符
	guestUsername := "guest_" + deviceID
	if len(guestUsername) > 50 {
		guestUsername = guestUsername[:50]
	}
	
	// 首先尝试查找已存在的游客用户
	if err := db.Where("username = ? AND is_guest = true", guestUsername).First(&user).Error; err != nil {
		// 如果没有找到，创建新的游客用户
		user = models.User{
			Username: guestUsername,
			Nickname: "游客" + generateRandomString(6),
			Avatar:   "",
			Role:     "guest",
			IsGuest:  true,
		}
		
		if err := db.Create(&user).Error; err != nil {
			ErrorResponse(c, http.StatusInternalServerError, "创建游客用户失败")
			return
		}
	}
	
	// 生成JWT令牌
	token, err := middleware.GenerateJWT(user.ID, "", user.Role)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "生成令牌失败")
		return
	}
	
	// 记录操作日志
	LogOperation(c, "LOGIN", "GUEST", fmt.Sprintf("游客登录: %s", user.Username))

	// 返回登录结果
	SuccessResponse(c, GuestLoginResponse{
		Token:    token,
		User:     user,
		ExpireAt: time.Now().Add(24 * time.Hour).Unix(),
	})
}



// getWechatOpenID 获取微信OpenID
func getWechatOpenID(code string) (string, error) {
	appID := os.Getenv("WX_APPID")
	appSecret := os.Getenv("WX_SECRET")
	
	if appID == "" || appSecret == "" {
		return "", fmt.Errorf("微信配置未设置")
	}

	url := fmt.Sprintf("https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code",
		appID, appSecret, code)

	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var wxResp WechatSessionResponse
	if err := json.Unmarshal(body, &wxResp); err != nil {
		return "", err
	}

	if wxResp.ErrCode != 0 {
		return "", fmt.Errorf("微信API错误: %s", wxResp.ErrMsg)
	}

	return wxResp.OpenID, nil
}

// generateRandomString 生成随机字符串
func generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	rand.Read(b)
	for i := range b {
		b[i] = charset[b[i]%byte(len(charset))]
	}
	return string(b)
}
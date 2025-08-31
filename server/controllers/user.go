package controllers

import (
	"fmt"
	"net/http"
	"qaminiprogram/config"
	"qaminiprogram/models"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// UpdateUserProfileRequest 更新用户信息请求
type UpdateUserProfileRequest struct {
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
}

// UpdateUserStatusRequest 更新用户状态请求
type UpdateUserStatusRequest struct {
	Status string `json:"status"`
}

// AdminUpdateUserRequest 管理员更新用户请求
type AdminUpdateUserRequest struct {
	Username   string `json:"username"`
	Email      string `json:"email"`
	Role       string `json:"role"`
	Status     string `json:"status"`
	IsVerified bool   `json:"isVerified"`
	Password   string `json:"password,omitempty"`
}

// AdminCreateUserRequest 管理员创建用户请求
type AdminCreateUserRequest struct {
	Username   string `json:"username"`
	Email      string `json:"email"`
	Password   string `json:"password"`
	Role       string `json:"role"`
	Status     string `json:"status"`
	IsVerified bool   `json:"isVerified"`
}

// HashPassword 密码加密函数
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// GetUserProfile 获取用户信息
func GetUserProfile(c *gin.Context) {
	userID, exists := GetUserID(c)
	if !exists {
		ErrorResponse(c, http.StatusUnauthorized, "未授权")
		return
	}

	db := config.GetDB()
	var user models.User
	if err := db.Where("id = ?", userID).First(&user).Error; err != nil {
		ErrorResponse(c, http.StatusNotFound, "用户不存在")
		return
	}

	SuccessResponse(c, user)
}

// UpdateUserProfile 更新用户信息
func UpdateUserProfile(c *gin.Context) {
	userID, exists := GetUserID(c)
	if !exists {
		ErrorResponse(c, http.StatusUnauthorized, "未授权")
		return
	}

	var req UpdateUserProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ErrorResponse(c, http.StatusBadRequest, "参数错误")
		return
	}

	db := config.GetDB()
	var user models.User
	if err := db.Where("id = ?", userID).First(&user).Error; err != nil {
		ErrorResponse(c, http.StatusNotFound, "用户不存在")
		return
	}

	// 更新用户信息
	if req.Nickname != "" {
		user.Nickname = req.Nickname
	}
	if req.Avatar != "" {
		user.Avatar = req.Avatar
	}

	if err := db.Save(&user).Error; err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "更新失败")
		return
	}

	// 记录操作日志
	LogOperation(c, "UPDATE", "USER", fmt.Sprintf("更新用户: %s", user.Username))

	SuccessResponse(c, user)
}

// GetUsers 获取用户列表（管理员）
func GetUsers(c *gin.Context) {
	page, size := ParsePageParams(c)
	offset := (page - 1) * size

	// 获取搜索参数
	username := c.Query("username")
	email := c.Query("email")
	role := c.Query("role")
	status := c.Query("status")

	db := config.GetDB()
	var users []models.User
	var total int64

	// 构建查询条件
	query := db.Model(&models.User{})
	if username != "" {
		query = query.Where("username LIKE ?", "%"+username+"%")
	}
	if email != "" {
		query = query.Where("email LIKE ?", "%"+email+"%")
	}
	if role != "" {
		query = query.Where("role = ?", role)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	// 获取总数
	query.Count(&total)

	// 获取用户列表
	if err := query.Offset(offset).Limit(size).Find(&users).Error; err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "获取用户列表失败")
		return
	}

	PageSuccessResponse(c, users, total, page, size)
}

// GetUserByID 根据ID获取用户（管理员）
func GetUserByID(c *gin.Context) {
	id, err := ParseIDParam(c, "id")
	if err != nil {
		ErrorResponse(c, http.StatusBadRequest, "无效的用户ID")
		return
	}

	db := config.GetDB()
	var user models.User
	if err := db.Where("id = ?", id).First(&user).Error; err != nil {
		ErrorResponse(c, http.StatusNotFound, "用户不存在")
		return
	}

	SuccessResponse(c, user)
}

// UpdateUser 更新用户（管理员）
func UpdateUser(c *gin.Context) {
	id, err := ParseIDParam(c, "id")
	if err != nil {
		ErrorResponse(c, http.StatusBadRequest, "无效的用户ID")
		return
	}

	var req AdminUpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ErrorResponse(c, http.StatusBadRequest, "参数错误")
		return
	}

	db := config.GetDB()
	var user models.User
	if err := db.Where("id = ?", id).First(&user).Error; err != nil {
		ErrorResponse(c, http.StatusNotFound, "用户不存在")
		return
	}

	// 更新用户信息
	if req.Username != "" {
		user.Username = req.Username
	}
	if req.Email != "" {
		user.Email = req.Email
	}
	if req.Role != "" {
		user.Role = req.Role
	}
	if req.Status != "" {
		user.Status = req.Status
	}
	user.IsVerified = req.IsVerified

	// 如果提供了密码，则更新密码
	if req.Password != "" {
		hashedPassword, err := HashPassword(req.Password)
		if err != nil {
			ErrorResponse(c, http.StatusInternalServerError, "密码加密失败")
			return
		}
		user.Password = hashedPassword
	}

	if err := db.Save(&user).Error; err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "更新失败")
		return
	}

	SuccessResponse(c, user)
}

// DeleteUser 删除用户（管理员）
func DeleteUser(c *gin.Context) {
	id, err := ParseIDParam(c, "id")
	if err != nil {
		ErrorResponse(c, http.StatusBadRequest, "无效的用户ID")
		return
	}

	db := config.GetDB()
	// 获取用户名用于日志记录
	var user models.User
	db.Where("id = ?", id).First(&user)
	username := user.Username

	if err := db.Delete(&models.User{}, id).Error; err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "删除失败")
		return
	}

	// 记录操作日志
	LogOperation(c, "DELETE", "USER", fmt.Sprintf("删除用户: %s", username))

	SuccessResponse(c, gin.H{"message": "删除成功"})
}

// BatchDeleteUsers 批量删除用户（管理员）
func BatchDeleteUsers(c *gin.Context) {
	var req BatchDeleteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ErrorResponse(c, http.StatusBadRequest, "参数错误")
		return
	}

	if len(req.IDs) == 0 {
		ErrorResponse(c, http.StatusBadRequest, "请选择要删除的用户")
		return
	}

	db := config.GetDB()
	// 获取要删除的用户名列表用于日志记录
	var users []models.User
	db.Where("id IN ?", req.IDs).Find(&users)
	usernames := make([]string, len(users))
	for i, user := range users {
		usernames[i] = user.Username
	}

	if err := db.Where("id IN ?", req.IDs).Delete(&models.User{}).Error; err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "批量删除失败")
		return
	}

	// 记录操作日志
	LogOperation(c, "BATCH_DELETE", "USER", fmt.Sprintf("批量删除用户: %v", usernames))

	SuccessResponse(c, gin.H{"message": "批量删除成功"})
}

// UpdateUserStatus 更新用户状态（管理员）
func UpdateUserStatus(c *gin.Context) {
	id, err := ParseIDParam(c, "id")
	if err != nil {
		ErrorResponse(c, http.StatusBadRequest, "无效的用户ID")
		return
	}

	var req UpdateUserStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ErrorResponse(c, http.StatusBadRequest, "参数错误")
		return
	}

	db := config.GetDB()
	var user models.User
	if err := db.Where("id = ?", id).First(&user).Error; err != nil {
		ErrorResponse(c, http.StatusNotFound, "用户不存在")
		return
	}

	oldStatus := user.Status
	user.Status = req.Status
	if err := db.Save(&user).Error; err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "更新状态失败")
		return
	}

	// 记录操作日志
	LogOperation(c, "UPDATE_STATUS", "USER", fmt.Sprintf("将用户 %s 状态从 %s 更新为 %s", user.Username, oldStatus, req.Status))

	SuccessResponse(c, user)
}

// CreateUser 创建用户（管理员）
func CreateUser(c *gin.Context) {
	var req AdminCreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ErrorResponse(c, http.StatusBadRequest, "参数错误")
		return
	}

	// 验证必填字段
	if req.Username == "" || req.Email == "" || req.Password == "" {
		ErrorResponse(c, http.StatusBadRequest, "用户名、邮箱和密码不能为空")
		return
	}

	db := config.GetDB()

	// 检查用户名是否已存在
	var existingUser models.User
	if err := db.Where("username = ? OR email = ?", req.Username, req.Email).First(&existingUser).Error; err == nil {
		ErrorResponse(c, http.StatusBadRequest, "用户名或邮箱已存在")
		return
	}

	// 加密密码
	hashedPassword, err := HashPassword(req.Password)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "密码加密失败")
		return
	}

	// 创建用户
	user := models.User{
		OpenID:     fmt.Sprintf("admin_created_%s", req.Username), // 管理员创建的用户不设置OpenID
		Username:   req.Username,
		Email:      req.Email,
		Password:   hashedPassword,
		Nickname:   req.Username, // 默认使用用户名作为昵称
		Role:       req.Role,
		Status:     req.Status,
		IsVerified: req.IsVerified,
	}

	if err := db.Create(&user).Error; err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "创建用户失败")
		return
	}

	// 记录操作日志
	LogOperation(c, "CREATE", "USER", fmt.Sprintf("创建用户: %s", user.Username))

	SuccessResponse(c, user)
}
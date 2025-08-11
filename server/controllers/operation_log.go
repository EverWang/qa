package controllers

import (
	"fmt"
	"net/http"
	"qaminiprogram/config"
	"qaminiprogram/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetOperationLogs 获取操作日志列表
func GetOperationLogs(c *gin.Context) {
	db := config.GetDB()

	// 获取查询参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	operator := c.Query("operator")
	action := c.Query("action")
	startTime := c.Query("startTime")
	endTime := c.Query("endTime")

	// 构建查询
	query := db.Model(&models.OperationLog{})

	// 添加筛选条件
	if operator != "" {
		query = query.Where("operator LIKE ?", "%"+operator+"%")
	}
	if action != "" {
		query = query.Where("action = ?", action)
	}
	if startTime != "" {
		query = query.Where("created_at >= ?", startTime+" 00:00:00")
	}
	if endTime != "" {
		query = query.Where("created_at <= ?", endTime+" 23:59:59")
	}

	// 获取总数
	var total int64
	countErr := query.Count(&total).Error
	if countErr != nil {
		ErrorResponse(c, http.StatusInternalServerError, "获取操作日志总数失败")
		return
	}

	// 分页查询
	var logs []models.OperationLog
	offset := (page - 1) * pageSize
	err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&logs).Error
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "获取操作日志失败")
		return
	}

	// 返回结果
	SuccessResponse(c, gin.H{
		"data":     logs,
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
	})
}

// CreateOperationLog 创建操作日志
func CreateOperationLog(operator, action, resource, description, ip, userAgent string) {
	db := config.GetDB()

	log := models.OperationLog{
		Operator:    operator,
		Action:      action,
		Resource:    resource,
		Description: description,
		IP:          ip,
		UserAgent:   userAgent,
	}

	db.Create(&log)
}

// LogOperation 记录操作日志的中间件辅助函数
func LogOperation(c *gin.Context, action, resource, description string) {
	// 获取操作人信息
	operator := "Unknown"
	
	// 尝试从JWT中获取用户ID和角色
	if userID, exists := c.Get("userId"); exists {
		if role, roleExists := c.Get("role"); roleExists {
			db := config.GetDB()
			
			// 确保userID是正确的类型
			var uid uint
			switch v := userID.(type) {
			case uint:
				uid = v
			case float64:
				uid = uint(v)
			default:
				// 如果类型转换失败，使用默认值
				operator = "Unknown(类型错误)"
				goto createLog
			}
			
			if role == "admin" {
				// 管理员用户，查询管理员表
				var admin models.Admin
				if err := db.Where("id = ?", uid).First(&admin).Error; err == nil {
					operator = admin.Username
				} else {
					operator = fmt.Sprintf("Admin_%d", uid)
				}
			} else {
				// 普通用户，查询用户表
				var user models.User
				if err := db.Where("id = ?", uid).First(&user).Error; err == nil {
					if user.Nickname != "" {
						operator = user.Nickname
					} else if user.OpenID != nil && len(*user.OpenID) >= 8 {
						operator = "用户" + (*user.OpenID)[:8] // 使用OpenID前8位作为标识
					} else {
						operator = fmt.Sprintf("User_%d", uid)
					}
				} else {
					operator = fmt.Sprintf("User_%d", uid)
				}
			}
		} else {
			operator = "Unknown(无角色)"
		}
	} else {
		operator = "Unknown(无用户ID)"
	}

createLog:
	// 获取客户端信息
	ip := c.ClientIP()
	userAgent := c.GetHeader("User-Agent")

	// 创建操作日志
	CreateOperationLog(operator, action, resource, description, ip, userAgent)
}
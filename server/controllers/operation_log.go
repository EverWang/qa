package controllers

import (
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
	query.Count(&total)

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
	operator := "unknown"
	if username, exists := c.Get("username"); exists {
		operator = username.(string)
	}

	// 获取客户端信息
	ip := c.ClientIP()
	userAgent := c.GetHeader("User-Agent")

	// 创建操作日志
	CreateOperationLog(operator, action, resource, description, ip, userAgent)
}
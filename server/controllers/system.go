package controllers

import (
	"net/http"
	"qaminiprogram/config"
	"qaminiprogram/models"

	"github.com/gin-gonic/gin"
)

// BasicSettingsRequest 基础设置请求
type BasicSettingsRequest struct {
	SystemName        string `json:"system_name"`
	SystemDescription string `json:"system_description"`
	SystemVersion     string `json:"system_version"`
	ContactEmail      string `json:"contact_email"`
	SystemStatus      string `json:"system_status"`
	MaintenanceNotice string `json:"maintenance_notice"`
}

// QuizSettingsRequest 答题设置请求
type QuizSettingsRequest struct {
	DailyLimit       int      `json:"daily_limit"`
	TimeLimit        int      `json:"time_limit"`
	EnablePoints     bool     `json:"enable_points"`
	CorrectPoints    int      `json:"correct_points"`
	WrongPoints      int      `json:"wrong_points"`
	QuizModes        []string `json:"quiz_modes"`
	ShowExplanation  string   `json:"show_explanation"`
}

// SystemStatistics 系统统计数据
type SystemStatistics struct {
	TotalUsers    int64 `json:"total_users"`
	TotalQuestions int64 `json:"total_questions"`
	TotalAnswers  int64 `json:"total_answers"`
	ActiveUsers   int64 `json:"active_users"`
	UserStats     []UserStatItem   `json:"user_stats"`
	AnswerStats   []AnswerStatItem `json:"answer_stats"`
}

// UserStatItem 用户统计项
type UserStatItem struct {
	Date        string `json:"date"`
	NewUsers    int64  `json:"new_users"`
	ActiveUsers int64  `json:"active_users"`
}

// AnswerStatItem 答题统计项
type AnswerStatItem struct {
	Date         string `json:"date"`
	TotalAnswers int64  `json:"total_answers"`
	CorrectRate  string `json:"correct_rate"`
}

/**
 * 获取基础设置
 */
func GetBasicSettings(c *gin.Context) {
	db := config.GetDB()

	// 从数据库获取系统设置
	var settings models.SystemSetting
	err := db.Where("key = ?", "basic").First(&settings).Error
	if err != nil {
		// 如果没有设置，返回默认值
		defaultSettings := BasicSettingsRequest{
			SystemName:        "刷刷题",
			SystemDescription: "专业的在线刷题平台",
			SystemVersion:     "1.0.0",
			ContactEmail:      "admin@example.com",
			SystemStatus:      "normal",
			MaintenanceNotice: "",
		}
		SuccessResponse(c, defaultSettings)
		return
	}

	// 解析设置值
	var basicSettings BasicSettingsRequest
	if err := parseJSONValue(settings.Value, &basicSettings); err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "解析设置失败")
		return
	}

	SuccessResponse(c, basicSettings)
}

/**
 * 更新基础设置
 */
func UpdateBasicSettings(c *gin.Context) {
	var req BasicSettingsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ErrorResponse(c, http.StatusBadRequest, "参数错误")
		return
	}

	db := config.GetDB()

	// 将设置保存到数据库
	value, err := toJSONString(req)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "序列化设置失败")
		return
	}

	// 更新或创建设置
	var settings models.SystemSetting
	err = db.Where("key = ?", "basic").First(&settings).Error
	if err != nil {
		// 创建新设置
		settings = models.SystemSetting{
			Key:   "basic",
			Value: value,
		}
		if err := db.Create(&settings).Error; err != nil {
			ErrorResponse(c, http.StatusInternalServerError, "保存设置失败")
			return
		}
	} else {
		// 更新现有设置
		settings.Value = value
		if err := db.Save(&settings).Error; err != nil {
			ErrorResponse(c, http.StatusInternalServerError, "更新设置失败")
			return
		}
	}

	SuccessResponse(c, gin.H{"message": "基础设置保存成功"})
}

/**
 * 获取答题设置
 */
func GetQuizSettings(c *gin.Context) {
	db := config.GetDB()

	// 从数据库获取答题设置
	var settings models.SystemSetting
	err := db.Where("key = ?", "quiz").First(&settings).Error
	if err != nil {
		// 如果没有设置，返回默认值
		defaultSettings := QuizSettingsRequest{
			DailyLimit:      0,
			TimeLimit:       0,
			EnablePoints:    true,
			CorrectPoints:   1,
			WrongPoints:     0,
			QuizModes:       []string{"random", "category"},
			ShowExplanation: "after_answer",
		}
		SuccessResponse(c, defaultSettings)
		return
	}

	// 解析设置值
	var quizSettings QuizSettingsRequest
	if err := parseJSONValue(settings.Value, &quizSettings); err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "解析设置失败")
		return
	}

	SuccessResponse(c, quizSettings)
}

/**
 * 更新答题设置
 */
func UpdateQuizSettings(c *gin.Context) {
	var req QuizSettingsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ErrorResponse(c, http.StatusBadRequest, "参数错误")
		return
	}

	db := config.GetDB()

	// 将设置保存到数据库
	value, err := toJSONString(req)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "序列化设置失败")
		return
	}

	// 更新或创建设置
	var settings models.SystemSetting
	err = db.Where("key = ?", "quiz").First(&settings).Error
	if err != nil {
		// 创建新设置
		settings = models.SystemSetting{
			Key:   "quiz",
			Value: value,
		}
		if err := db.Create(&settings).Error; err != nil {
			ErrorResponse(c, http.StatusInternalServerError, "保存设置失败")
			return
		}
	} else {
		// 更新现有设置
		settings.Value = value
		if err := db.Save(&settings).Error; err != nil {
			ErrorResponse(c, http.StatusInternalServerError, "更新设置失败")
			return
		}
	}

	SuccessResponse(c, gin.H{"message": "答题设置保存成功"})
}

/**
 * 获取系统统计数据
 */
func GetSystemStatistics(c *gin.Context) {
	db := config.GetDB()

	// 获取总用户数
	var totalUsers int64
	db.Model(&models.User{}).Count(&totalUsers)

	// 获取总题目数
	var totalQuestions int64
	db.Model(&models.Question{}).Count(&totalQuestions)

	// 获取总答题数
	var totalAnswers int64
	db.Model(&models.AnswerRecord{}).Count(&totalAnswers)

	// 获取活跃用户数（最近7天有答题记录的用户）
	var activeUsers int64
	db.Model(&models.AnswerRecord{}).
		Where("created_at >= DATE_SUB(NOW(), INTERVAL 7 DAY)").
		Distinct("user_id").
		Count(&activeUsers)

	// 构造统计数据
	stats := SystemStatistics{
		TotalUsers:     totalUsers,
		TotalQuestions: totalQuestions,
		TotalAnswers:   totalAnswers,
		ActiveUsers:    activeUsers,
		UserStats:      []UserStatItem{},
		AnswerStats:    []AnswerStatItem{},
	}

	SuccessResponse(c, stats)
}

/**
 * 导出统计报表
 */
func ExportStatistics(c *gin.Context) {
	// 这里可以实现Excel导出功能
	// 暂时返回成功响应
	SuccessResponse(c, gin.H{"message": "导出功能开发中"})
}
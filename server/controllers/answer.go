package controllers

import (
	"net/http"
	"qaminiprogram/config"
	"qaminiprogram/models"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// SubmitAnswerRequest 提交答案请求
type SubmitAnswerRequest struct {
	QuestionID uuid.UUID `json:"questionId" binding:"required"`
	UserAnswer int       `json:"userAnswer" binding:"min=0"`
	TimeSpent  int       `json:"timeSpent" binding:"min=0"`
}

// AnswerStatistics 答题统计
type AnswerStatistics struct {
	TotalAnswered   int64   `json:"totalAnswered"`
	CorrectAnswered int64   `json:"correctAnswered"`
	AccuracyRate    float64 `json:"accuracyRate"`
	TotalTimeSpent  int64   `json:"totalTimeSpent"`
	AverageTime     float64 `json:"averageTime"`
	TodayAnswered   int64   `json:"todayAnswered"`
	WeekAnswered    int64   `json:"weekAnswered"`
	MonthAnswered   int64   `json:"monthAnswered"`
}

// CategoryStatistics 分类统计
type CategoryStatistics struct {
	CategoryID      uuid.UUID `json:"categoryId"`
	CategoryName    string    `json:"categoryName"`
	TotalAnswered   int64     `json:"totalAnswered"`
	CorrectAnswered int64     `json:"correctAnswered"`
	AccuracyRate    float64   `json:"accuracyRate"`
}

// SubmitAnswer 提交答案
func SubmitAnswer(c *gin.Context) {
	userID, exists := GetUserID(c)
	if !exists {
		ErrorResponse(c, http.StatusUnauthorized, "未授权")
		return
	}

	var req SubmitAnswerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ErrorResponse(c, http.StatusBadRequest, "参数错误")
		return
	}

	db := config.GetDB()

	// 验证题目是否存在
	var question models.Question
	if err := db.Where("id = ?", req.QuestionID).First(&question).Error; err != nil {
		ErrorResponse(c, http.StatusNotFound, "题目不存在")
		return
	}

	// 验证答案索引是否有效
	if req.UserAnswer >= len(question.Options) {
		ErrorResponse(c, http.StatusBadRequest, "答案索引无效")
		return
	}

	// 判断答案是否正确
	isCorrect := req.UserAnswer == question.CorrectAnswer

	// 检查是否已经答过这道题
	var existingRecord models.AnswerRecord
	result := db.Where("user_id = ? AND question_id = ?", userID, req.QuestionID).First(&existingRecord)

	if result.Error == nil {
		// 更新已有记录
		existingRecord.UserAnswer = req.UserAnswer
		existingRecord.IsCorrect = isCorrect
		existingRecord.TimeSpent = req.TimeSpent
		if err := db.Save(&existingRecord).Error; err != nil {
			ErrorResponse(c, http.StatusInternalServerError, "更新答题记录失败")
			return
		}
		
		// 如果答错了，添加到错题本
		if !isCorrect {
			addToMistakeBook(db, userID, req.QuestionID)
		} else {
			// 如果答对了，从错题本中移除
			removeFromMistakeBook(db, userID, req.QuestionID)
		}
		
		SuccessResponse(c, gin.H{
			"isCorrect": isCorrect,
			"correctAnswer": question.CorrectAnswer,
			"explanation": question.Explanation,
			"updated": true,
		})
		return
	}

	// 创建新的答题记录
	answerRecord := models.AnswerRecord{
		UserID:     userID,
		QuestionID: req.QuestionID,
		UserAnswer: req.UserAnswer,
		IsCorrect:  isCorrect,
		TimeSpent:  req.TimeSpent,
	}

	if err := db.Create(&answerRecord).Error; err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "保存答题记录失败")
		return
	}

	// 如果答错了，添加到错题本
	if !isCorrect {
		addToMistakeBook(db, userID, req.QuestionID)
	}

	SuccessResponse(c, gin.H{
		"isCorrect": isCorrect,
		"correctAnswer": question.CorrectAnswer,
		"explanation": question.Explanation,
		"updated": false,
	})
}

// GetAnswerHistory 获取答题历史
func GetAnswerHistory(c *gin.Context) {
	userID, exists := GetUserID(c)
	if !exists {
		ErrorResponse(c, http.StatusUnauthorized, "未授权")
		return
	}

	page, size := ParsePageParams(c)
	offset := (page - 1) * size

	// 获取查询参数
	categoryID := c.Query("category_id")
	isCorrect := c.Query("is_correct")

	db := config.GetDB()
	query := db.Model(&models.AnswerRecord{}).Where("user_id = ?", userID)

	// 添加查询条件
	if isCorrect != "" {
		query = query.Where("is_correct = ?", isCorrect == "true")
	}

	// 如果指定了分类，需要关联题目表
	if categoryID != "" {
		query = query.Joins("JOIN questions ON answer_records.question_id = questions.id").Where("questions.category_id = ?", categoryID)
	}

	var total int64
	query.Count(&total)

	var records []models.AnswerRecord
	if err := query.Preload("Question").Preload("Question.Category").Order("created_at DESC").Offset(offset).Limit(size).Find(&records).Error; err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "获取答题历史失败")
		return
	}

	PageSuccessResponse(c, records, total, page, size)
}

// GetAnswerStatistics 获取答题统计
func GetAnswerStatistics(c *gin.Context) {
	userID, exists := GetUserID(c)
	if !exists {
		ErrorResponse(c, http.StatusUnauthorized, "未授权")
		return
	}

	db := config.GetDB()
	var stats AnswerStatistics

	// 总答题数
	db.Model(&models.AnswerRecord{}).Where("user_id = ?", userID).Count(&stats.TotalAnswered)

	// 正确答题数
	db.Model(&models.AnswerRecord{}).Where("user_id = ? AND is_correct = ?", userID, true).Count(&stats.CorrectAnswered)

	// 计算正确率
	if stats.TotalAnswered > 0 {
		stats.AccuracyRate = float64(stats.CorrectAnswered) / float64(stats.TotalAnswered) * 100
	}

	// 总用时
	db.Model(&models.AnswerRecord{}).Where("user_id = ?", userID).Select("COALESCE(SUM(time_spent), 0)").Scan(&stats.TotalTimeSpent)

	// 平均用时
	if stats.TotalAnswered > 0 {
		stats.AverageTime = float64(stats.TotalTimeSpent) / float64(stats.TotalAnswered)
	}

	// 今日答题数
	today := time.Now().Format("2006-01-02")
	db.Model(&models.AnswerRecord{}).Where("user_id = ? AND DATE(created_at) = ?", userID, today).Count(&stats.TodayAnswered)

	// 本周答题数
	weekStart := time.Now().AddDate(0, 0, -int(time.Now().Weekday())).Format("2006-01-02")
	db.Model(&models.AnswerRecord{}).Where("user_id = ? AND DATE(created_at) >= ?", userID, weekStart).Count(&stats.WeekAnswered)

	// 本月答题数
	monthStart := time.Now().AddDate(0, 0, -time.Now().Day()+1).Format("2006-01-02")
	db.Model(&models.AnswerRecord{}).Where("user_id = ? AND DATE(created_at) >= ?", userID, monthStart).Count(&stats.MonthAnswered)

	// 分类统计
	var categoryStats []CategoryStatistics
	db.Raw(`
		SELECT 
			c.id as category_id,
			c.name as category_name,
			COUNT(ar.id) as total_answered,
			SUM(CASE WHEN ar.is_correct = 1 THEN 1 ELSE 0 END) as correct_answered,
			CASE 
				WHEN COUNT(ar.id) > 0 THEN 
					SUM(CASE WHEN ar.is_correct = 1 THEN 1 ELSE 0 END) * 100.0 / COUNT(ar.id)
				ELSE 0 
			END as accuracy_rate
		FROM categories c
		LEFT JOIN questions q ON c.id = q.category_id
		LEFT JOIN answer_records ar ON q.id = ar.question_id AND ar.user_id = ?
		WHERE ar.id IS NOT NULL
		GROUP BY c.id, c.name
		ORDER BY total_answered DESC
	`, userID).Scan(&categoryStats)

	SuccessResponse(c, gin.H{
		"overall": stats,
		"categories": categoryStats,
	})
}

// addToMistakeBook 添加到错题本
func addToMistakeBook(db *gorm.DB, userID uint, questionID uuid.UUID) {
	var existing models.MistakeBook
	result := db.Where("user_id = ? AND question_id = ?", userID, questionID).First(&existing)
	if result.Error != nil {
		// 不存在则创建
		mistake := models.MistakeBook{
			UserID:     userID,
			QuestionID: questionID,
		}
		db.Create(&mistake)
	}
}

// removeFromMistakeBook 从错题本移除
func removeFromMistakeBook(db *gorm.DB, userID uint, questionID uuid.UUID) {
	db.Where("user_id = ? AND question_id = ?", userID, questionID).Delete(&models.MistakeBook{})
}

// GetCategoryProgress 获取用户在特定分类下的进度
func GetCategoryProgress(c *gin.Context) {
	userID, exists := GetUserID(c)
	if !exists {
		ErrorResponse(c, http.StatusUnauthorized, "未授权")
		return
	}

	categoryID, err := ParseIDParam(c, "id")
	if err != nil {
		ErrorResponse(c, http.StatusBadRequest, "无效的分类ID")
		return
	}

	db := config.GetDB()
	var progress CategoryStatistics

	// 查询该分类下用户的答题统计
	err = db.Raw(`
		SELECT 
			c.id as category_id,
			c.name as category_name,
			COUNT(ar.id) as total_answered,
			SUM(CASE WHEN ar.is_correct = 1 THEN 1 ELSE 0 END) as correct_answered,
			CASE 
				WHEN COUNT(ar.id) > 0 THEN 
					SUM(CASE WHEN ar.is_correct = 1 THEN 1 ELSE 0 END) * 100.0 / COUNT(ar.id)
				ELSE 0 
			END as accuracy_rate
		FROM categories c
		LEFT JOIN questions q ON c.id = q.category_id
		LEFT JOIN answer_records ar ON q.id = ar.question_id AND ar.user_id = ?
		WHERE c.id = ?
		GROUP BY c.id, c.name
	`, userID, categoryID).Scan(&progress).Error

	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "获取分类进度失败")
		return
	}

	// 如果没有答题记录，设置默认值
	if progress.CategoryID == uuid.Nil {
		// 获取分类信息
		var category models.Category
		if err := db.Where("id = ?", categoryID).First(&category).Error; err != nil {
			ErrorResponse(c, http.StatusNotFound, "分类不存在")
			return
		}
		progress = CategoryStatistics{
			CategoryID:      categoryID,
			CategoryName:    category.Name,
			TotalAnswered:   0,
			CorrectAnswered: 0,
			AccuracyRate:    0,
		}
	}

	SuccessResponse(c, progress)
}
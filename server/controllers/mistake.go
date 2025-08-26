package controllers

import (
	"net/http"
	"qaminiprogram/config"
	"qaminiprogram/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// AddToMistakeBookRequest 添加到错题本请求
type AddToMistakeBookRequest struct {
	QuestionID uuid.UUID `json:"questionId" binding:"required"`
}

// GetMistakeBooks 获取错题本
func GetMistakeBooks(c *gin.Context) {
	userID, exists := GetUserID(c)
	if !exists {
		ErrorResponse(c, http.StatusUnauthorized, "未授权")
		return
	}

	page, size := ParsePageParams(c)
	offset := (page - 1) * size

	// 获取查询参数
	categoryID := c.Query("categoryId")
	difficulty := c.Query("difficulty")
	isMastered := c.Query("isMastered")

	db := config.GetDB()
	query := db.Model(&models.MistakeBook{}).Where("user_id = ?", userID)

	// 筛选掌握状态
	if isMastered != "" {
		if isMastered == "true" {
			query = query.Where("is_mastered = ?", true)
		} else if isMastered == "false" {
			query = query.Where("is_mastered = ?", false)
		}
	}

	// 如果指定了分类或难度，需要关联题目表
	if categoryID != "" || difficulty != "" {
		query = query.Joins("JOIN questions ON mistake_books.question_id = questions.id")
		if categoryID != "" {
			query = query.Where("questions.category_id = ?", categoryID)
		}
		if difficulty != "" {
			query = query.Where("questions.difficulty = ?", difficulty)
		}
	}

	var total int64
	query.Count(&total)

	var mistakes []models.MistakeBook
	if err := query.Preload("Question").Preload("Question.Category").Order("created_at DESC").Offset(offset).Limit(size).Find(&mistakes).Error; err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "获取错题本失败")
		return
	}

	PageSuccessResponse(c, mistakes, total, page, size)
}

// AddToMistakeBook 添加到错题本
func AddToMistakeBook(c *gin.Context) {
	userID, exists := GetUserID(c)
	if !exists {
		ErrorResponse(c, http.StatusUnauthorized, "未授权")
		return
	}

	var req AddToMistakeBookRequest
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

	// 检查是否已经在错题本中
	var existing models.MistakeBook
	result := db.Where("user_id = ? AND question_id = ?", userID, req.QuestionID).First(&existing)
	if result.Error == nil {
		ErrorResponse(c, http.StatusConflict, "该题目已在错题本中")
		return
	}

	// 添加到错题本
	mistake := models.MistakeBook{
		UserID:     userID,
		QuestionID: req.QuestionID,
	}

	if err := db.Create(&mistake).Error; err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "添加到错题本失败")
		return
	}

	// 预加载关联数据
	db.Preload("Question").Preload("Question.Category").First(&mistake, mistake.ID)

	SuccessResponse(c, mistake)
}

// RemoveFromMistakeBook 从错题本移除
func RemoveFromMistakeBook(c *gin.Context) {
	userID, exists := GetUserID(c)
	if !exists {
		ErrorResponse(c, http.StatusUnauthorized, "未授权")
		return
	}

	id, err := ParseIDParam(c, "id")
	if err != nil {
		ErrorResponse(c, http.StatusBadRequest, "无效的错题本ID")
		return
	}

	db := config.GetDB()

	// 验证错题本记录是否存在且属于当前用户
	var mistake models.MistakeBook
	if err := db.Where("id = ? AND user_id = ?", id, userID).First(&mistake).Error; err != nil {
		ErrorResponse(c, http.StatusNotFound, "错题本记录不存在")
		return
	}

	// 删除记录
	if err := db.Delete(&mistake).Error; err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "移除失败")
		return
	}

	SuccessResponse(c, gin.H{"message": "移除成功"})
}

// ClearMistakeBook 清空错题本
func ClearMistakeBook(c *gin.Context) {
	userID, exists := GetUserID(c)
	if !exists {
		ErrorResponse(c, http.StatusUnauthorized, "未授权")
		return
	}

	db := config.GetDB()

	// 删除用户的所有错题本记录
	result := db.Where("user_id = ?", userID).Delete(&models.MistakeBook{})
	if result.Error != nil {
		ErrorResponse(c, http.StatusInternalServerError, "清空错题本失败")
		return
	}

	SuccessResponse(c, gin.H{
		"message": "清空成功",
		"deletedCount": result.RowsAffected,
	})
}

// MarkMistakeAsMastered 标记错题为已掌握
func MarkMistakeAsMastered(c *gin.Context) {
	userID, exists := GetUserID(c)
	if !exists {
		ErrorResponse(c, http.StatusUnauthorized, "未授权")
		return
	}

	questionID, err := ParseIDParam(c, "questionId")
	if err != nil {
		ErrorResponse(c, http.StatusBadRequest, "无效的题目ID")
		return
	}

	db := config.GetDB()

	// 查找错题本记录
	var mistake models.MistakeBook
	if err := db.Where("user_id = ? AND question_id = ?", userID, questionID).First(&mistake).Error; err != nil {
		ErrorResponse(c, http.StatusNotFound, "错题本记录不存在")
		return
	}

	// 更新掌握状态
	mistake.IsMastered = true
	if err := db.Save(&mistake).Error; err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "更新掌握状态失败")
		return
	}

	SuccessResponse(c, gin.H{"message": "已标记为掌握"})
}

// ResetMistakeStatus 重置错题状态
func ResetMistakeStatus(c *gin.Context) {
	userID, exists := GetUserID(c)
	if !exists {
		ErrorResponse(c, http.StatusUnauthorized, "未授权")
		return
	}

	questionID, err := ParseIDParam(c, "questionId")
	if err != nil {
		ErrorResponse(c, http.StatusBadRequest, "无效的题目ID")
		return
	}

	db := config.GetDB()

	// 查找错题本记录
	var mistake models.MistakeBook
	if err := db.Where("user_id = ? AND question_id = ?", userID, questionID).First(&mistake).Error; err != nil {
		ErrorResponse(c, http.StatusNotFound, "错题本记录不存在")
		return
	}

	// 重置掌握状态
	mistake.IsMastered = false
	if err := db.Save(&mistake).Error; err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "重置掌握状态失败")
		return
	}

	SuccessResponse(c, gin.H{"message": "已重置掌握状态"})
}

// GetMistakeBookStatistics 获取错题本统计
func GetMistakeBookStatistics(c *gin.Context) {
	userID, exists := GetUserID(c)
	if !exists {
		ErrorResponse(c, http.StatusUnauthorized, "未授权")
		return
	}

	db := config.GetDB()

	// 总错题数
	var totalMistakes int64
	db.Model(&models.MistakeBook{}).Where("user_id = ?", userID).Count(&totalMistakes)

	// 已掌握错题数
	var masteredMistakes int64
	db.Model(&models.MistakeBook{}).Where("user_id = ? AND is_mastered = ?", userID, true).Count(&masteredMistakes)

	// 按分类统计错题数
	type CategoryMistakeStats struct {
		CategoryID   uuid.UUID `json:"categoryId"`
		CategoryName string    `json:"categoryName"`
		MistakeCount int64     `json:"mistakeCount"`
	}

	var categoryStats []CategoryMistakeStats
	db.Raw(`
		SELECT 
			c.id as category_id,
			c.name as category_name,
			COUNT(mb.id) as mistake_count
		FROM categories c
		JOIN questions q ON c.id = q.category_id
		JOIN mistake_books mb ON q.id = mb.question_id
		WHERE mb.user_id = ?
		GROUP BY c.id, c.name
		ORDER BY mistake_count DESC
	`, userID).Scan(&categoryStats)

	// 按难度统计错题数
	type DifficultyMistakeStats struct {
		Difficulty   string `json:"difficulty"`
		MistakeCount int64  `json:"mistakeCount"`
	}

	var difficultyStats []DifficultyMistakeStats
	db.Raw(`
		SELECT 
			q.difficulty,
			COUNT(mb.id) as mistake_count
		FROM questions q
		JOIN mistake_books mb ON q.id = mb.question_id
		WHERE mb.user_id = ?
		GROUP BY q.difficulty
		ORDER BY mistake_count DESC
	`, userID).Scan(&difficultyStats)

	SuccessResponse(c, gin.H{
		"totalMistakes": totalMistakes,
		"masteredMistakes": masteredMistakes,
		"categoryStats": categoryStats,
		"difficultyStats": difficultyStats,
	})
}
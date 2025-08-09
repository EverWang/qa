package controllers

import (
	"qaminiprogram/config"
	"qaminiprogram/models"
	"time"

	"github.com/gin-gonic/gin"
)

// OverviewStatistics 概览统计
type OverviewStatistics struct {
	TotalUsers      int64 `json:"totalUsers"`
	TotalQuestions  int64 `json:"totalQuestions"`
	TotalCategories int64 `json:"totalCategories"`
	TotalAnswers    int64 `json:"totalAnswers"`
	TodayUsers      int64 `json:"todayUsers"`
	TodayAnswers    int64 `json:"todayAnswers"`
	WeekUsers       int64 `json:"weekUsers"`
	WeekAnswers     int64 `json:"weekAnswers"`
	MonthUsers      int64 `json:"monthUsers"`
	MonthAnswers    int64 `json:"monthAnswers"`
}

// QuestionStatistics 题目统计
type QuestionStatistics struct {
	QuestionID      uint    `json:"question_id"`
	QuestionTitle   string  `json:"question_title"`
	CategoryName    string  `json:"category_name"`
	TotalAnswered   int64   `json:"total_answered"`
	CorrectAnswered int64   `json:"correct_answered"`
	AccuracyRate    float64 `json:"accuracy_rate"`
	Difficulty      string  `json:"difficulty"`
}

// UserStatistics 用户统计
type UserStatistics struct {
	UserID          uint    `json:"user_id"`
	Nickname        string  `json:"nickname"`
	TotalAnswered   int64   `json:"total_answered"`
	CorrectAnswered int64   `json:"correct_answered"`
	AccuracyRate    float64 `json:"accuracy_rate"`
	TotalTimeSpent  int64   `json:"total_time_spent"`
	LastActiveTime  string  `json:"last_active_time"`
}

// GetOverviewStatistics 获取概览统计（管理员）
func GetOverviewStatistics(c *gin.Context) {
	db := config.GetDB()
	var stats OverviewStatistics

	// 总用户数
	db.Model(&models.User{}).Count(&stats.TotalUsers)

	// 总题目数
	db.Model(&models.Question{}).Count(&stats.TotalQuestions)

	// 总分类数
	db.Model(&models.Category{}).Count(&stats.TotalCategories)

	// 总答题数
	db.Model(&models.AnswerRecord{}).Count(&stats.TotalAnswers)

	// 今日数据
	today := time.Now().Format("2006-01-02")
	db.Model(&models.User{}).Where("DATE(created_at) = ?", today).Count(&stats.TodayUsers)
	db.Model(&models.AnswerRecord{}).Where("DATE(created_at) = ?", today).Count(&stats.TodayAnswers)

	// 本周数据
	weekStart := time.Now().AddDate(0, 0, -int(time.Now().Weekday())).Format("2006-01-02")
	db.Model(&models.User{}).Where("DATE(created_at) >= ?", weekStart).Count(&stats.WeekUsers)
	db.Model(&models.AnswerRecord{}).Where("DATE(created_at) >= ?", weekStart).Count(&stats.WeekAnswers)

	// 本月数据
	monthStart := time.Now().AddDate(0, 0, -time.Now().Day()+1).Format("2006-01-02")
	db.Model(&models.User{}).Where("DATE(created_at) >= ?", monthStart).Count(&stats.MonthUsers)
	db.Model(&models.AnswerRecord{}).Where("DATE(created_at) >= ?", monthStart).Count(&stats.MonthAnswers)

	SuccessResponse(c, stats)
}

// GetQuestionStatistics 获取题目统计（管理员）
func GetQuestionStatistics(c *gin.Context) {
	page, size := ParsePageParams(c)
	offset := (page - 1) * size

	// 获取查询参数
	categoryID := c.Query("category_id")
	difficulty := c.Query("difficulty")
	sortBy := c.DefaultQuery("sort_by", "total_answered")
	sortOrder := c.DefaultQuery("sort_order", "DESC")

	db := config.GetDB()

	// 构建查询
	query := `
		SELECT 
			q.id as question_id,
			q.title as question_title,
			c.name as category_name,
			q.difficulty,
			COUNT(ar.id) as total_answered,
			SUM(CASE WHEN ar.is_correct = 1 THEN 1 ELSE 0 END) as correct_answered,
			CASE 
				WHEN COUNT(ar.id) > 0 THEN 
					SUM(CASE WHEN ar.is_correct = 1 THEN 1 ELSE 0 END) * 100.0 / COUNT(ar.id)
				ELSE 0 
			END as accuracy_rate
		FROM questions q
		LEFT JOIN categories c ON q.category_id = c.id
		LEFT JOIN answer_records ar ON q.id = ar.question_id
		WHERE 1=1
	`

	args := []interface{}{}

	if categoryID != "" {
		query += " AND q.category_id = ?"
		args = append(args, categoryID)
	}

	if difficulty != "" {
		query += " AND q.difficulty = ?"
		args = append(args, difficulty)
	}

	query += " GROUP BY q.id, q.title, c.name, q.difficulty"

	// 添加排序
	validSortFields := map[string]bool{
		"total_answered":   true,
		"correct_answered": true,
		"accuracy_rate":    true,
	}

	if validSortFields[sortBy] && (sortOrder == "ASC" || sortOrder == "DESC") {
		query += " ORDER BY " + sortBy + " " + sortOrder
	} else {
		query += " ORDER BY total_answered DESC"
	}

	// 获取总数
	countQuery := `
		SELECT COUNT(*) FROM (
			SELECT q.id
			FROM questions q
			LEFT JOIN categories c ON q.category_id = c.id
			WHERE 1=1
	`

	countArgs := []interface{}{}
	if categoryID != "" {
		countQuery += " AND q.category_id = ?"
		countArgs = append(countArgs, categoryID)
	}
	if difficulty != "" {
		countQuery += " AND q.difficulty = ?"
		countArgs = append(countArgs, difficulty)
	}
	countQuery += ") as subquery"

	var total int64
	db.Raw(countQuery, countArgs...).Scan(&total)

	// 添加分页
	query += " LIMIT ? OFFSET ?"
	args = append(args, size, offset)

	var stats []QuestionStatistics
	db.Raw(query, args...).Scan(&stats)

	PageSuccessResponse(c, stats, total, page, size)
}

// GetUserStatistics 获取用户统计（管理员）
func GetUserStatistics(c *gin.Context) {
	page, size := ParsePageParams(c)
	offset := (page - 1) * size

	// 获取查询参数
	sortBy := c.DefaultQuery("sort_by", "total_answered")
	sortOrder := c.DefaultQuery("sort_order", "DESC")

	db := config.GetDB()

	// 构建查询
	query := `
		SELECT 
			u.id as user_id,
			u.nickname,
			COUNT(ar.id) as total_answered,
			SUM(CASE WHEN ar.is_correct = 1 THEN 1 ELSE 0 END) as correct_answered,
			CASE 
				WHEN COUNT(ar.id) > 0 THEN 
					SUM(CASE WHEN ar.is_correct = 1 THEN 1 ELSE 0 END) * 100.0 / COUNT(ar.id)
				ELSE 0 
			END as accuracy_rate,
			COALESCE(SUM(ar.time_spent), 0) as total_time_spent,
			MAX(ar.created_at) as last_active_time
		FROM users u
		LEFT JOIN answer_records ar ON u.id = ar.user_id
		WHERE u.role = 'user'
		GROUP BY u.id, u.nickname
	`

	// 添加排序
	validSortFields := map[string]bool{
		"total_answered":   true,
		"correct_answered": true,
		"accuracy_rate":    true,
		"total_time_spent": true,
		"last_active_time": true,
	}

	if validSortFields[sortBy] && (sortOrder == "ASC" || sortOrder == "DESC") {
		query += " ORDER BY " + sortBy + " " + sortOrder
	} else {
		query += " ORDER BY total_answered DESC"
	}

	// 获取总数
	var total int64
	db.Model(&models.User{}).Where("role = ?", "user").Count(&total)

	// 添加分页
	query += " LIMIT ? OFFSET ?"

	var stats []UserStatistics
	db.Raw(query, size, offset).Scan(&stats)

	// 格式化时间
	for i := range stats {
		if stats[i].LastActiveTime != "" {
			if t, err := time.Parse("2006-01-02 15:04:05", stats[i].LastActiveTime); err == nil {
				stats[i].LastActiveTime = t.Format("2006-01-02 15:04:05")
			}
		}
	}

	PageSuccessResponse(c, stats, total, page, size)
}
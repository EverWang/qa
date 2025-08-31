package controllers

import (
	"fmt"
	"net/http"
	"qaminiprogram/config"
	"qaminiprogram/models"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	
	"gorm.io/gorm"
)

// CreateQuestionRequest 创建题目请求
type CreateQuestionRequest struct {
	Title         string    `json:"title" binding:"required"`
	Content       string    `json:"content" binding:"required"`
	Type          string    `json:"type" binding:"required"`
	Options       []string  `json:"options"`
	CorrectAnswer int       `json:"correctAnswer" binding:"min=0"`
	Explanation   string    `json:"explanation"`
	Difficulty    string    `json:"difficulty"`
	CategoryID    uint `json:"categoryId" binding:"required"`
}

// UpdateQuestionRequest 更新题目请求
type UpdateQuestionRequest struct {
	Title         string     `json:"title"`
	Content       string     `json:"content"`
	Type          string     `json:"type"`
	Options       []string   `json:"options"`
	CorrectAnswer *int       `json:"correctAnswer"`
	Explanation   string     `json:"explanation"`
	Difficulty    string     `json:"difficulty"`
	CategoryID    *uint `json:"categoryId"`
}

// BatchDeleteRequest 批量删除请求
type BatchDeleteRequest struct {
	IDs []uint `json:"ids"`
}

// GetQuestions 获取题目列表
func GetQuestions(c *gin.Context) {
	page, size := ParsePageParams(c)
	offset := (page - 1) * size

	// 获取查询参数
	categoryID := c.Query("categoryId")
	difficulty := c.Query("difficulty")
	keyword := c.Query("keyword")

	db := config.GetDB()
	query := db.Model(&models.Question{})

	// 添加查询条件
	if categoryID != "" {
		query = query.Where("category_id = ?", categoryID)
	}
	if difficulty != "" {
		query = query.Where("difficulty = ?", difficulty)
	}
	if keyword != "" {
		query = query.Where("title LIKE ? OR content LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	var total int64
	query.Count(&total)

	var questions []models.Question
	if err := query.Preload("Category").Preload("Creator").Offset(offset).Limit(size).Find(&questions).Error; err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "获取题目列表失败")
		return
	}

	PageSuccessResponse(c, questions, total, page, size)
}

// GetAdminQuestions 获取管理员题目列表（带更多筛选条件）
func GetAdminQuestions(c *gin.Context) {
	page, size := ParsePageParams(c)
	offset := (page - 1) * size

	// 获取查询参数
	categoryID := c.Query("categoryId")
	difficulty := c.Query("difficulty")
	questionType := c.Query("type") // 添加题目类型参数
	keyword := c.Query("keyword")
	creatorID := c.Query("creatorId")
	startDate := c.Query("startDate")
	endDate := c.Query("endDate")

	db := config.GetDB()
	query := db.Model(&models.Question{})

	// 添加查询条件
	if categoryID != "" {
		query = query.Where("category_id = ?", categoryID)
	}
	if difficulty != "" {
		query = query.Where("difficulty = ?", difficulty)
	}
	if questionType != "" {
		query = query.Where("type = ?", questionType)
	}
	if keyword != "" {
		query = query.Where("title LIKE ? OR content LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}
	if creatorID != "" {
		query = query.Where("creator_id = ?", creatorID)
	}
	if startDate != "" {
		query = query.Where("created_at >= ?", startDate)
	}
	if endDate != "" {
		query = query.Where("created_at <= ?", endDate)
	}

	// 获取总数
	var total int64
	if err := query.Count(&total).Error; err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "获取题目总数失败")
		return
	}

	// 获取分页数据
	var questions []models.Question
	if err := query.Preload("Category").Preload("Creator").Order("created_at DESC").Offset(offset).Limit(size).Find(&questions).Error; err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "获取题目列表失败")
		return
	}

	PageSuccessResponse(c, questions, total, page, size)
}

// GetQuestionByID 根据ID获取题目
func GetQuestionByID(c *gin.Context) {
	id, err := ParseIDParam(c, "id")
	if err != nil {
		ErrorResponse(c, http.StatusBadRequest, "无效的题目ID")
		return
	}

	db := config.GetDB()
	var question models.Question
	if err := db.Preload("Category").Preload("Creator").Where("id = ?", id).First(&question).Error; err != nil {
		ErrorResponse(c, http.StatusNotFound, "题目不存在")
		return
	}

	SuccessResponse(c, question)
}

// GetRandomQuestions 获取随机题目
func GetRandomQuestions(c *gin.Context) {
	countStr := c.DefaultQuery("count", "10")
	count, err := strconv.Atoi(countStr)
	if err != nil || count < 1 || count > 50 {
		count = 10
	}

	// 获取查询参数
	categoryID := c.Query("categoryId")
	difficulty := c.Query("difficulty")

	db := config.GetDB()
	query := db.Model(&models.Question{})

	// 添加查询条件
	if categoryID != "" {
		query = query.Where("category_id = ?", categoryID)
	}
	if difficulty != "" {
		query = query.Where("difficulty = ?", difficulty)
	}

	var questions []models.Question
	if err := query.Preload("Category").Order("RAND()").Limit(count).Find(&questions).Error; err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "获取随机题目失败")
		return
	}

	SuccessResponse(c, questions)
}

// GetQuestionsByCategory 根据分类获取题目
func GetQuestionsByCategory(c *gin.Context) {
	categoryID, err := ParseIDParam(c, "categoryId")
	if err != nil {
		ErrorResponse(c, http.StatusBadRequest, "无效的分类ID")
		return
	}

	page, size := ParsePageParams(c)
	offset := (page - 1) * size

	db := config.GetDB()
	var total int64
	db.Model(&models.Question{}).Where("category_id = ?", categoryID).Count(&total)

	var questions []models.Question
	if err := db.Preload("Category").Where("category_id = ?", categoryID).Offset(offset).Limit(size).Find(&questions).Error; err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "获取题目失败")
		return
	}

	PageSuccessResponse(c, questions, total, page, size)
}

// CreateQuestion 创建题目（管理员）
func CreateQuestion(c *gin.Context) {
	// 绑定请求参数
	var req CreateQuestionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ErrorResponse(c, http.StatusBadRequest, "参数错误")
		return
	}

	// 验证题目类型
	validTypes := []string{"single", "multiple", "judge", "fill"}
	isValidType := false
	for _, validType := range validTypes {
		if req.Type == validType {
			isValidType = true
			break
		}
	}
	if !isValidType {
		ErrorResponse(c, http.StatusBadRequest, "题目类型无效")
		return
	}

	// 验证选项（仅对选择题）
	if (req.Type == "single" || req.Type == "multiple") && len(req.Options) < 2 {
		ErrorResponse(c, http.StatusBadRequest, "选择题至少需要2个选项")
		return
	}

	// 验证正确答案索引（仅对选择题）
	if (req.Type == "single" || req.Type == "multiple") && req.CorrectAnswer >= len(req.Options) {
		ErrorResponse(c, http.StatusBadRequest, "正确答案索引超出范围")
		return
	}

	// 验证分类是否存在
	db := config.GetDB()
	var category models.Category
	if err := db.Where("id = ?", req.CategoryID).First(&category).Error; err != nil {
		ErrorResponse(c, http.StatusBadRequest, "分类不存在")
		return
	}

	// 获取创建者ID
	creatorID, exists := GetUserID(c)
	if !exists {
		ErrorResponse(c, http.StatusUnauthorized, "无法获取用户信息")
		return
	}

	// 创建题目
	question := models.Question{
		Title:         req.Title,
		Content:       req.Content,
		Type:          req.Type,
		Options:       models.JSONArray(req.Options),
		CorrectAnswer: req.CorrectAnswer,
		Explanation:   req.Explanation,
		Difficulty:    req.Difficulty,
		CategoryID:    req.CategoryID,
		CreatorID:     &creatorID,
	}

	if err := db.Create(&question).Error; err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "创建题目失败")
		return
	}

	// 预加载关联数据
	db.Preload("Category").Preload("Creator").First(&question, question.ID)

	// 记录操作日志
	LogOperation(c, "CREATE", "QUESTION", fmt.Sprintf("创建题目: %s", question.Title))

	SuccessResponse(c, question)
}

// UpdateQuestion 更新题目（管理员）
func UpdateQuestion(c *gin.Context) {
	id, err := ParseIDParam(c, "id")
	if err != nil {
		ErrorResponse(c, http.StatusBadRequest, "无效的题目ID")
		return
	}

	var req UpdateQuestionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ErrorResponse(c, http.StatusBadRequest, "参数错误")
		return
	}

	db := config.GetDB()
	var question models.Question
	if err := db.Where("id = ?", id).First(&question).Error; err != nil {
		ErrorResponse(c, http.StatusNotFound, "题目不存在")
		return
	}

	// 更新字段
	if req.Title != "" {
		question.Title = req.Title
	}
	if req.Content != "" {
		question.Content = req.Content
	}
	if req.Type != "" {
		// 验证题目类型
		validTypes := []string{"single", "multiple", "judge", "fill"}
		isValidType := false
		for _, validType := range validTypes {
			if req.Type == validType {
				isValidType = true
				break
			}
		}
		if !isValidType {
			ErrorResponse(c, http.StatusBadRequest, "题目类型无效")
			return
		}
		question.Type = req.Type
	}
	if len(req.Options) > 0 {
		question.Options = models.JSONArray(req.Options)
	}
	if req.CorrectAnswer != nil {
		// 根据题目类型验证正确答案
		questionType := req.Type
		if questionType == "" {
			questionType = question.Type // 使用现有类型
		}
		
		if questionType == "single" || questionType == "judge" {
			// 单选题和判断题：验证答案索引
			if *req.CorrectAnswer >= len(question.Options) {
				ErrorResponse(c, http.StatusBadRequest, "正确答案索引超出范围")
				return
			}
		} else if questionType == "multiple" {
			// 多选题：验证位掩码
			maxValidMask := (1 << len(question.Options)) - 1
			if *req.CorrectAnswer < 0 || *req.CorrectAnswer > maxValidMask {
				ErrorResponse(c, http.StatusBadRequest, "多选题答案格式错误")
				return
			}
			// 确保至少选择了一个答案
			if *req.CorrectAnswer == 0 {
				ErrorResponse(c, http.StatusBadRequest, "多选题至少需要选择一个正确答案")
				return
			}
		} else if questionType == "fill" {
			// 填空题：答案通常为0，实际答案存储在explanation中
			if *req.CorrectAnswer != 0 {
				ErrorResponse(c, http.StatusBadRequest, "填空题答案索引应为0")
				return
			}
		}
		
		question.CorrectAnswer = *req.CorrectAnswer
	}
	if req.Explanation != "" {
		question.Explanation = req.Explanation
	}
	if req.Difficulty != "" {
		question.Difficulty = req.Difficulty
	}
	if req.CategoryID != nil {
		// 验证分类是否存在
		var category models.Category
		if err := db.Where("id = ?", *req.CategoryID).First(&category).Error; err != nil {
			ErrorResponse(c, http.StatusBadRequest, "分类不存在")
			return
		}
		question.CategoryID = *req.CategoryID
	}

	if err := db.Save(&question).Error; err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "更新题目失败")
		return
	}

	// 预加载关联数据
	db.Preload("Category").Preload("Creator").First(&question, question.ID)

	// 记录操作日志
	LogOperation(c, "UPDATE", "QUESTION", fmt.Sprintf("更新题目: %s", question.Title))

	SuccessResponse(c, question)
}

// DeleteQuestion 删除题目（管理员）
func DeleteQuestion(c *gin.Context) {
	id, err := ParseIDParam(c, "id")
	if err != nil {
		ErrorResponse(c, http.StatusBadRequest, "无效的题目ID")
		return
	}

	db := config.GetDB()
	// 获取题目标题用于日志记录
	var question models.Question
	db.Where("id = ?", id).First(&question)
	questionTitle := question.Title

	if err := db.Delete(&models.Question{}, id).Error; err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "删除题目失败")
		return
	}

	// 记录操作日志
	LogOperation(c, "DELETE", "QUESTION", fmt.Sprintf("删除题目: %s", questionTitle))

	SuccessResponse(c, gin.H{"message": "删除成功"})
}

// BatchDeleteQuestions 批量删除题目（管理员）
func BatchDeleteQuestions(c *gin.Context) {
	var req BatchDeleteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ErrorResponse(c, http.StatusBadRequest, "参数错误")
		return
	}

	if len(req.IDs) == 0 {
		ErrorResponse(c, http.StatusBadRequest, "请选择要删除的题目")
		return
	}

	db := config.GetDB()
	// 获取要删除的题目标题列表用于日志记录
	var questions []models.Question
	db.Where("id IN ?", req.IDs).Find(&questions)
	questionTitles := make([]string, len(questions))
	for i, question := range questions {
		questionTitles[i] = question.Title
	}

	if err := db.Where("id IN ?", req.IDs).Delete(&models.Question{}).Error; err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "批量删除失败")
		return
	}

	// 记录操作日志
	LogOperation(c, "BATCH_DELETE", "QUESTION", fmt.Sprintf("批量删除题目: %v", questionTitles))

	SuccessResponse(c, gin.H{"message": "批量删除成功"})
}

// ImportQuestionItem 导入题目项
type ImportQuestionItem struct {
	Content    string    `json:"content" binding:"required"`
	Type       string    `json:"type" binding:"required"`
	Difficulty string    `json:"difficulty" binding:"required"`
	CategoryID uint `json:"category_id" binding:"required"`
	Options    []string  `json:"options"`
	Answer     string    `json:"answer" binding:"required"`
	Explanation string   `json:"explanation"`
}

// ImportQuestionsRequest 批量导入题目请求
type ImportQuestionsRequest struct {
	Questions []ImportQuestionItem `json:"questions" binding:"required,min=1"`
	Options   ImportOptions        `json:"options"`
}

// ImportOptions 导入选项
type ImportOptions struct {
	SkipDuplicates  bool `json:"skip_duplicates"`
	UpdateExisting  bool `json:"update_existing"`
}

// ImportResult 导入结果
type ImportResult struct {
	Imported int      `json:"imported"`
	Skipped  int      `json:"skipped"`
	Errors   []string `json:"errors,omitempty"`
}

// ImportQuestions 批量导入题目（管理员）
func ImportQuestions(c *gin.Context) {
	// 绑定请求参数
	var req ImportQuestionsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ErrorResponse(c, http.StatusBadRequest, "参数错误")
		return
	}

	// 获取创建者ID
	creatorID, exists := GetUserID(c)
	if !exists {
		ErrorResponse(c, http.StatusUnauthorized, "无法获取用户信息")
		return
	}

	db := config.GetDB()
	result := ImportResult{
		Imported: 0,
		Skipped:  0,
		Errors:   []string{},
	}

	// 开始事务
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	for i, item := range req.Questions {
		// 验证题目类型
		validTypes := []string{"single", "multiple", "judge", "fill"}
		isValidType := false
		for _, validType := range validTypes {
			if item.Type == validType {
				isValidType = true
				break
			}
		}
		if !isValidType {
			result.Errors = append(result.Errors, fmt.Sprintf("第%d行：题目类型无效", i+1))
			result.Skipped++
			continue
		}

		// 验证难度等级
		validDifficulties := []string{"easy", "medium", "hard"}
		isValidDifficulty := false
		for _, validDifficulty := range validDifficulties {
			if item.Difficulty == validDifficulty {
				isValidDifficulty = true
				break
			}
		}
		if !isValidDifficulty {
			result.Errors = append(result.Errors, fmt.Sprintf("第%d行：难度等级无效", i+1))
			result.Skipped++
			continue
		}

		// 验证分类是否存在
		var category models.Category
		if err := tx.Where("id = ?", item.CategoryID).First(&category).Error; err != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("第%d行：分类不存在", i+1))
			result.Skipped++
			continue
		}

		// 检查是否重复（基于内容）
		if req.Options.SkipDuplicates {
			var existingQuestion models.Question
			if err := tx.Where("content = ? AND category_id = ?", item.Content, item.CategoryID).First(&existingQuestion).Error; err == nil {
				if req.Options.UpdateExisting {
					// 更新现有题目
					if err := updateExistingQuestion(tx, &existingQuestion, item); err != nil {
						result.Errors = append(result.Errors, fmt.Sprintf("第%d行：更新失败 - %v", i+1, err))
						result.Skipped++
					} else {
						result.Imported++
					}
				} else {
					result.Skipped++
				}
				continue
			}
		}

		// 处理答案和选项
		options := item.Options
		correctAnswer := 0

		switch item.Type {
		case "single":
			if len(options) < 2 {
				result.Errors = append(result.Errors, fmt.Sprintf("第%d行：单选题至少需要2个选项", i+1))
				result.Skipped++
				continue
			}
			// 查找答案索引
			answerIndex := -1
			for j := range options {
				if strings.ToUpper(item.Answer) == string(rune('A'+j)) {
					answerIndex = j
					break
				}
			}
			if answerIndex == -1 {
				result.Errors = append(result.Errors, fmt.Sprintf("第%d行：单选题答案格式错误", i+1))
				result.Skipped++
				continue
			}
			correctAnswer = answerIndex

		case "multiple":
			if len(options) < 2 {
				result.Errors = append(result.Errors, fmt.Sprintf("第%d行：多选题至少需要2个选项", i+1))
				result.Skipped++
				continue
			}
			// 多选题答案存储为位掩码
			answers := strings.Split(strings.ToUpper(item.Answer), ",")
			mask := 0
			for _, ans := range answers {
				ans = strings.TrimSpace(ans)
				if len(ans) == 1 && ans[0] >= 'A' && ans[0] <= 'Z' {
					index := int(ans[0] - 'A')
					if index < len(options) {
						mask |= (1 << index)
					}
				}
			}
			if mask == 0 {
				result.Errors = append(result.Errors, fmt.Sprintf("第%d行：多选题答案格式错误", i+1))
				result.Skipped++
				continue
			}
			correctAnswer = mask

		case "judge":
			options = []string{"正确", "错误"}
			if strings.ToLower(item.Answer) == "true" || item.Answer == "正确" {
				correctAnswer = 0
			} else if strings.ToLower(item.Answer) == "false" || item.Answer == "错误" {
				correctAnswer = 1
			} else {
				result.Errors = append(result.Errors, fmt.Sprintf("第%d行：判断题答案格式错误", i+1))
				result.Skipped++
				continue
			}

		case "fill":
			options = []string{} // 填空题没有选项
			correctAnswer = 0    // 填空题答案存储在explanation中
		}

		// 创建题目
		question := models.Question{
			Title:         item.Content[:min(len(item.Content), 100)], // 截取前100字符作为标题
			Content:       item.Content,
			Options:       models.JSONArray(options),
			CorrectAnswer: correctAnswer,
			Explanation:   item.Explanation,
			Difficulty:    item.Difficulty,
			CategoryID:    item.CategoryID,
			CreatorID:     &creatorID,
		}

		// 填空题的答案存储在explanation中
		if item.Type == "fill" {
			if question.Explanation == "" {
				question.Explanation = item.Answer
			} else {
				question.Explanation = item.Answer + "\n\n" + question.Explanation
			}
		}

		if err := tx.Create(&question).Error; err != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("第%d行：创建失败 - %v", i+1, err))
			result.Skipped++
			continue
		}

		result.Imported++
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "导入失败")
		return
	}

	SuccessResponse(c, result)
}

// updateExistingQuestion 更新现有题目
func updateExistingQuestion(tx *gorm.DB, question *models.Question, item ImportQuestionItem) error {
	// 处理选项和答案
	options := item.Options
	correctAnswer := 0

	switch item.Type {
	case "single":
		for j := range options {
			if strings.ToUpper(item.Answer) == string(rune('A'+j)) {
				correctAnswer = j
				break
			}
		}
	case "multiple":
		answers := strings.Split(strings.ToUpper(item.Answer), ",")
		mask := 0
		for _, ans := range answers {
			ans = strings.TrimSpace(ans)
			if len(ans) == 1 && ans[0] >= 'A' && ans[0] <= 'Z' {
				index := int(ans[0] - 'A')
				if index < len(options) {
					mask |= (1 << index)
				}
			}
		}
		correctAnswer = mask
	case "judge":
		options = []string{"正确", "错误"}
		if strings.ToLower(item.Answer) == "true" || item.Answer == "正确" {
			correctAnswer = 0
		} else {
			correctAnswer = 1
		}
	case "fill":
		options = []string{}
		correctAnswer = 0
	}

	// 更新字段
	question.Title = item.Content[:min(len(item.Content), 100)]
	question.Content = item.Content
	question.Options = models.JSONArray(options)
	question.CorrectAnswer = correctAnswer
	question.Explanation = item.Explanation
	question.Difficulty = item.Difficulty
	question.CategoryID = item.CategoryID

	// 填空题的答案存储在explanation中
	if item.Type == "fill" {
		if question.Explanation == "" {
			question.Explanation = item.Answer
		} else {
			question.Explanation = item.Answer + "\n\n" + question.Explanation
		}
	}

	return tx.Save(question).Error
}

// min 返回两个整数中的较小值
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// ExportQuestions 导出题目（管理员）
func ExportQuestions(c *gin.Context) {
	// 获取查询参数
	categoryID := c.Query("category_id")
	difficulty := c.Query("difficulty")
	keyword := c.Query("keyword")
	creatorID := c.Query("creator_id")
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	db := config.GetDB()
	query := db.Model(&models.Question{})

	// 添加查询条件
	if categoryID != "" {
		query = query.Where("category_id = ?", categoryID)
	}
	if difficulty != "" {
		query = query.Where("difficulty = ?", difficulty)
	}
	if keyword != "" {
		query = query.Where("title LIKE ? OR content LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}
	if creatorID != "" {
		query = query.Where("creator_id = ?", creatorID)
	}
	if startDate != "" {
		query = query.Where("created_at >= ?", startDate)
	}
	if endDate != "" {
		query = query.Where("created_at <= ?", endDate)
	}

	// 获取题目数据
	var questions []models.Question
	if err := query.Preload("Category").Preload("Creator").Order("created_at DESC").Find(&questions).Error; err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "获取题目数据失败")
		return
	}

	// 转换为导出格式
	exportData := make([]map[string]interface{}, 0, len(questions))
	for _, question := range questions {
		// 处理题目类型
		questionType := "single"
		if len(question.Options) == 0 {
			questionType = "fill"
		} else if len(question.Options) == 2 && question.Options[0] == "正确" && question.Options[1] == "错误" {
			questionType = "judge"
		} else {
			// 检查是否为多选题（通过答案是否为位掩码判断）
			if question.CorrectAnswer > len(question.Options)-1 {
				questionType = "multiple"
			}
		}

		// 处理答案格式
		var answer string
		switch questionType {
		case "single":
			if question.CorrectAnswer < len(question.Options) {
				answer = string(rune('A' + question.CorrectAnswer))
			}
		case "multiple":
			var answers []string
			mask := question.CorrectAnswer
			for i := 0; i < len(question.Options); i++ {
				if mask&(1<<i) != 0 {
					answers = append(answers, string(rune('A'+i)))
				}
			}
			answer = strings.Join(answers, ",")
		case "judge":
			if question.CorrectAnswer == 0 {
				answer = "true"
			} else {
				answer = "false"
			}
		case "fill":
			// 填空题答案从explanation中提取
			lines := strings.Split(question.Explanation, "\n")
			if len(lines) > 0 {
				answer = lines[0]
			}
		}

		// 处理选项
		optionMap := make(map[string]string)
		optionKeys := []string{"选项A", "选项B", "选项C", "选项D", "选项E", "选项F"}
		for i, key := range optionKeys {
			if i < len(question.Options) {
				optionMap[key] = question.Options[i]
			} else {
				optionMap[key] = ""
			}
		}

		// 获取分类名称
		categoryName := ""
		if question.Category != nil {
			categoryName = question.Category.Name
		}

		// 处理解析（移除答案部分）
		explanation := question.Explanation
		if questionType == "fill" && explanation != "" {
			lines := strings.Split(explanation, "\n")
			if len(lines) > 2 {
				explanation = strings.Join(lines[2:], "\n")
			} else {
				explanation = ""
			}
		}

		exportItem := map[string]interface{}{
			"题目内容": question.Content,
			"题目类型": questionType,
			"难度等级": question.Difficulty,
			"分类ID":  question.CategoryID,
			"分类名称": categoryName,
			"选项A":   optionMap["选项A"],
			"选项B":   optionMap["选项B"],
			"选项C":   optionMap["选项C"],
			"选项D":   optionMap["选项D"],
			"选项E":   optionMap["选项E"],
			"选项F":   optionMap["选项F"],
			"正确答案": answer,
			"题目解析": explanation,
			"创建时间": question.CreatedAt.Format("2006-01-02 15:04:05"),
		}
		exportData = append(exportData, exportItem)
	}

	// 设置响应头
	c.Header("Content-Type", "application/json")
	c.Header("Content-Disposition", "attachment; filename=questions.json")

	// 返回JSON格式的数据
	SuccessResponse(c, gin.H{
		"data":  exportData,
		"total": len(exportData),
		"exported_at": fmt.Sprintf("%d", time.Now().Unix()),
	})
}
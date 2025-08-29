package controllers

import (
	"fmt"
	"net/http"
	"qaminiprogram/config"
	"qaminiprogram/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// CreateCategoryRequest 创建分类请求
type CreateCategoryRequest struct {
	Name     string     `json:"name" binding:"required"`
	ParentID *uuid.UUID `json:"parentId"`
	Level    int        `json:"level"`
	Sort     int        `json:"sort"`
}

// UpdateCategoryRequest 更新分类请求
type UpdateCategoryRequest struct {
	Name        string     `json:"name"`
	Description string     `json:"description"`
	ParentID    *uuid.UUID `json:"parentId"`
	Level       *int       `json:"level"`
	SortOrder   *int       `json:"sortOrder"`
	Status      *int       `json:"status"`
}

// CategoryTreeNode 分类树节点
type CategoryTreeNode struct {
	models.Category
	Children []*CategoryTreeNode `json:"children"`
}

// GetCategories 获取分类列表
func GetCategories(c *gin.Context) {
	tree := c.DefaultQuery("tree", "false")
	parentID := c.Query("parentId")

	db := config.GetDB()

	if tree == "true" {
		// 返回树形结构
		categoryTree, err := buildCategoryTree(db)
		if err != nil {
			ErrorResponse(c, http.StatusInternalServerError, "获取分类树失败")
			return
		}
		SuccessResponse(c, categoryTree)
		return
	}

	// 返回平铺列表
	var categories []models.Category
	query := db.Model(&models.Category{})

	if parentID != "" {
		if parentID == "0" {
			// 获取顶级分类
			query = query.Where("parent_id IS NULL")
		} else {
			// 获取指定父分类的子分类
			query = query.Where("parent_id = ?", parentID)
		}
	}

	if err := query.Order("level ASC, sort ASC").Find(&categories).Error; err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "获取分类列表失败")
		return
	}

	// 为每个分类添加题目数量
	for i := range categories {
		var questionCount int64
		db.Model(&models.Question{}).Where("category_id = ?", categories[i].ID).Count(&questionCount)
		categories[i].QuestionCount = int(questionCount)
	}

	// 预加载父分类信息
	for i := range categories {
		if categories[i].ParentID != nil {
			db.Where("id = ?", *categories[i].ParentID).First(&categories[i].Parent)
		}
	}

	SuccessResponse(c, categories)
}

// GetAdminCategories 获取管理员分类列表（带分页）
func GetAdminCategories(c *gin.Context) {
	tree := c.DefaultQuery("tree", "false")
	parentID := c.Query("parentId")
	search := c.Query("search")
	status := c.Query("status")

	db := config.GetDB()

	if tree == "true" {
		// 返回树形结构
		categoryTree, err := buildCategoryTreeWithFilter(db, search, status)
		if err != nil {
			ErrorResponse(c, http.StatusInternalServerError, "获取分类树失败")
			return
		}
		SuccessResponse(c, categoryTree)
		return
	}

	// 解析分页参数
	page, size := ParsePageParams(c)
	offset := (page - 1) * size

	// 构建查询
	query := db.Model(&models.Category{})

	// 搜索条件
	if search != "" {
		query = query.Where("name LIKE ?", "%"+search+"%")
	}

	// 状态筛选
	if status != "" {
		query = query.Where("status = ?", status)
	}

	// 父分类筛选
	if parentID != "" {
		if parentID == "0" {
			// 获取顶级分类
			query = query.Where("parent_id IS NULL")
		} else {
			// 获取指定父分类的子分类
			query = query.Where("parent_id = ?", parentID)
		}
	}

	// 获取总数
	var total int64
	if err := query.Count(&total).Error; err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "获取分类总数失败")
		return
	}

	// 获取分页数据
	var categories []models.Category
	if err := query.Preload("Parent").Order("level ASC, sort ASC").Offset(offset).Limit(size).Find(&categories).Error; err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "获取分类列表失败")
		return
	}

	// 为每个分类添加题目数量
	for i := range categories {
		var questionCount int64
		db.Model(&models.Question{}).Where("category_id = ?", categories[i].ID).Count(&questionCount)
		categories[i].QuestionCount = int(questionCount)
	}

	PageSuccessResponse(c, categories, total, page, size)
}

// GetCategoryByID 根据ID获取分类
func GetCategoryByID(c *gin.Context) {
	id, err := ParseIDParam(c, "id")
	if err != nil {
		ErrorResponse(c, http.StatusBadRequest, "无效的分类ID")
		return
	}

	db := config.GetDB()
	var category models.Category
	if err := db.Preload("Parent").Preload("Children").Where("id = ?", id).First(&category).Error; err != nil {
		ErrorResponse(c, http.StatusNotFound, "分类不存在")
		return
	}

	SuccessResponse(c, category)
}

// CreateCategory 创建分类（管理员）
func CreateCategory(c *gin.Context) {
	var req CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ErrorResponse(c, http.StatusBadRequest, "参数错误")
		return
	}

	db := config.GetDB()

	// 验证父分类是否存在
	if req.ParentID != nil {
		var parentCategory models.Category
		if err := db.Where("id = ?", *req.ParentID).First(&parentCategory).Error; err != nil {
			ErrorResponse(c, http.StatusBadRequest, "父分类不存在")
			return
		}
		// 自动设置层级
		if req.Level == 0 {
			req.Level = parentCategory.Level + 1
		}
	} else {
		// 顶级分类
		if req.Level == 0 {
			req.Level = 1
		}
	}

	// 创建分类
	category := models.Category{
		Name:     req.Name,
		ParentID: req.ParentID,
		Level:    req.Level,
		Sort:     req.Sort,
	}

	if err := db.Create(&category).Error; err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "创建分类失败")
		return
	}

	// 预加载关联数据
	db.Preload("Parent").First(&category, category.ID)

	// 记录操作日志
	LogOperation(c, "CREATE", "CATEGORY", fmt.Sprintf("创建分类: %s", category.Name))

	SuccessResponse(c, category)
}

// UpdateCategory 更新分类（管理员）
func UpdateCategory(c *gin.Context) {
	id, err := ParseIDParam(c, "id")
	if err != nil {
		ErrorResponse(c, http.StatusBadRequest, "无效的分类ID")
		return
	}

	var req UpdateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ErrorResponse(c, http.StatusBadRequest, "参数错误")
		return
	}

	db := config.GetDB()
	var category models.Category
	if err := db.Where("id = ?", id).First(&category).Error; err != nil {
		ErrorResponse(c, http.StatusNotFound, "分类不存在")
		return
	}

	// 验证父分类
	if req.ParentID != nil {
		// 不能将分类设置为自己的子分类
		if *req.ParentID == id {
			ErrorResponse(c, http.StatusBadRequest, "不能将分类设置为自己的父分类")
			return
		}
		
		// 验证父分类是否存在
		var parentCategory models.Category
		if err := db.Where("id = ?", *req.ParentID).First(&parentCategory).Error; err != nil {
			ErrorResponse(c, http.StatusBadRequest, "父分类不存在")
			return
		}
	}

	// 更新字段
	if req.Name != "" {
		category.Name = req.Name
	}
	if req.Description != "" {
		category.Description = req.Description
	}
	// 处理ParentID字段，包括设置为nil的情况
	category.ParentID = req.ParentID
	// 根据ParentID自动计算层级
	if req.ParentID == nil {
		// 设置为大类
		category.Level = 1
	} else {
		// 根据父分类层级计算
		var parentCategory models.Category
		if err := db.Where("id = ?", *req.ParentID).First(&parentCategory).Error; err == nil {
			category.Level = parentCategory.Level + 1
		}
	}
	if req.Level != nil {
		category.Level = *req.Level
	}
	if req.SortOrder != nil {
		category.Sort = *req.SortOrder
	}
	if req.Status != nil {
		category.Status = *req.Status
	}

	if err := db.Save(&category).Error; err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "更新分类失败")
		return
	}

	// 预加载关联数据
	db.Preload("Parent").First(&category, category.ID)

	// 记录操作日志
	LogOperation(c, "UPDATE", "CATEGORY", fmt.Sprintf("更新分类: %s", category.Name))

	SuccessResponse(c, category)
}

// DeleteCategory 删除分类（管理员）
func DeleteCategory(c *gin.Context) {
	id, err := ParseIDParam(c, "id")
	if err != nil {
		ErrorResponse(c, http.StatusBadRequest, "无效的分类ID")
		return
	}

	db := config.GetDB()

	// 检查是否有子分类
	var childCount int64
	db.Model(&models.Category{}).Where("parent_id = ?", id).Count(&childCount)
	if childCount > 0 {
		ErrorResponse(c, http.StatusBadRequest, "该分类下还有子分类，无法删除")
		return
	}

	// 检查是否有题目
	var questionCount int64
	db.Model(&models.Question{}).Where("category_id = ?", id).Count(&questionCount)
	if questionCount > 0 {
		ErrorResponse(c, http.StatusBadRequest, "该分类下还有题目，无法删除")
		return
	}

	// 获取分类名称用于日志记录
	var category models.Category
	db.Where("id = ?", id).First(&category)
	categoryName := category.Name

	if err := db.Delete(&models.Category{}, id).Error; err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "删除分类失败")
		return
	}

	// 记录操作日志
	LogOperation(c, "DELETE", "CATEGORY", fmt.Sprintf("删除分类: %s", categoryName))

	SuccessResponse(c, gin.H{"message": "删除成功"})
}

// UpdateCategoryStatus 更新分类状态（管理员）
func UpdateCategoryStatus(c *gin.Context) {
	id, err := ParseIDParam(c, "id")
	if err != nil {
		ErrorResponse(c, http.StatusBadRequest, "无效的分类ID")
		return
	}

	var req struct {
		Status *int `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		ErrorResponse(c, http.StatusBadRequest, "参数错误")
		return
	}

	// 验证状态值
	if req.Status == nil || (*req.Status != 0 && *req.Status != 1) {
		ErrorResponse(c, http.StatusBadRequest, "状态值必须为0或1")
		return
	}

	db := config.GetDB()
	var category models.Category
	if err := db.Where("id = ?", id).First(&category).Error; err != nil {
		ErrorResponse(c, http.StatusNotFound, "分类不存在")
		return
	}

	// 更新状态
	oldStatus := category.Status
	category.Status = *req.Status
	if err := db.Save(&category).Error; err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "更新分类状态失败")
		return
	}

	// 记录操作日志
	statusText := "禁用"
	if *req.Status == 1 {
		statusText = "启用"
	}
	LogOperation(c, "UPDATE_STATUS", "CATEGORY", fmt.Sprintf("将分类 %s 状态从 %d 更新为 %s", category.Name, oldStatus, statusText))

	SuccessResponse(c, gin.H{"message": "状态更新成功"})
}

// buildCategoryTree 构建分类树
func buildCategoryTree(db *gorm.DB) ([]CategoryTreeNode, error) {
	return buildCategoryTreeWithFilter(db, "", "")
}

// buildCategoryTreeWithFilter 构建带筛选条件的分类树
func buildCategoryTreeWithFilter(db *gorm.DB, search, status string) ([]CategoryTreeNode, error) {
	var categories []models.Category
	query := db.Model(&models.Category{})
	
	// 添加搜索条件
	if search != "" {
		query = query.Where("name LIKE ?", "%"+search+"%")
	}
	
	// 添加状态筛选
	if status != "" {
		query = query.Where("status = ?", status)
	}
	
	if err := query.Order("level ASC, sort ASC").Find(&categories).Error; err != nil {
		return nil, err
	}

	// 为每个分类添加题目数量
	for i := range categories {
		var questionCount int64
		db.Model(&models.Question{}).Where("category_id = ?", categories[i].ID).Count(&questionCount)
		categories[i].QuestionCount = int(questionCount)
	}

	// 使用指针映射来确保修改能够正确反映
	categoryMap := make(map[uuid.UUID]*CategoryTreeNode)
	var allNodes []*CategoryTreeNode
	
	// 创建所有节点
	for _, category := range categories {
		node := &CategoryTreeNode{
			Category: category,
			Children: make([]*CategoryTreeNode, 0),
		}
		categoryMap[category.ID] = node
		allNodes = append(allNodes, node)
	}

	// 构建父子关系
	for _, node := range allNodes {
		if node.Category.ParentID != nil {
			if parentNode, exists := categoryMap[*node.Category.ParentID]; exists {
				parentNode.Children = append(parentNode.Children, node)
			}
		}
	}

	// 收集根节点
	var rootNodes []CategoryTreeNode
	for _, node := range allNodes {
		if node.Category.ParentID == nil {
			rootNodes = append(rootNodes, *node)
		}
	}

	return rootNodes, nil
}
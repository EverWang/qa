package models

import (
	"time"
	"database/sql/driver"
	"encoding/json"
	"errors"
)

// User 用户模型
type User struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	OpenID     *string   `json:"openid" gorm:"index;size:100"`
	Username   string    `json:"username" gorm:"size:50"`
	Email      string    `json:"email" gorm:"size:100"`
	Password   string    `json:"-" gorm:"size:255"`
	Nickname   string    `json:"nickname" gorm:"size:50;not null"`
	Avatar     string    `json:"avatar" gorm:"size:255"`
	Role       string    `json:"role" gorm:"type:enum('user','admin','guest');default:'user'"`
	Status     string    `json:"status" gorm:"type:enum('active','inactive','banned');default:'active'"`
	IsVerified bool      `json:"isVerified" gorm:"default:false"`
	IsGuest    bool      `json:"isGuest" gorm:"default:false"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}

// Category 分类模型
type Category struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name" gorm:"size:50;not null"`
	Description string    `json:"description" gorm:"size:255"`
	ParentID    *uint     `json:"parentId" gorm:"index"`
	Level       int       `json:"level" gorm:"default:1"`
	Sort        int       `json:"sortOrder" gorm:"default:0"`
	Status      int       `json:"status" gorm:"default:1;comment:状态 1-启用 0-禁用"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	
	// 关联
	Parent   *Category  `json:"parent,omitempty" gorm:"foreignKey:ParentID"`
	Children []Category `json:"children,omitempty" gorm:"foreignKey:ParentID"`
	
	// 计算字段（不存储在数据库中）
	QuestionCount int `json:"questionCount" gorm:"-"`
}

// JSONArray 自定义JSON数组类型
type JSONArray []string

// Scan 实现Scanner接口
func (j *JSONArray) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}
	
	switch v := value.(type) {
	case []byte:
		return json.Unmarshal(v, j)
	case string:
		return json.Unmarshal([]byte(v), j)
	default:
		return errors.New("cannot scan into JSONArray")
	}
}

// Value 实现Valuer接口
func (j JSONArray) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return json.Marshal(j)
}

// Question 题目模型
type Question struct {
	ID            uint      `json:"id" gorm:"primaryKey"`
	Title         string    `json:"title" gorm:"size:200;not null"`
	Content       string    `json:"content" gorm:"type:text;not null"`
	Type          string    `json:"type" gorm:"type:enum('single','multiple','judge','fill');default:'single'"`
	Options       JSONArray `json:"options" gorm:"type:json;not null"`
	CorrectAnswer int       `json:"correctAnswer" gorm:"not null"`
	Explanation   string    `json:"explanation" gorm:"type:text"`
	Difficulty    string    `json:"difficulty" gorm:"type:enum('easy','medium','hard');default:'medium'"`
	CategoryID    uint      `json:"categoryId" gorm:"not null;index"`
	CreatorID     *uint     `json:"creatorId" gorm:"index"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
	
	// 关联
	Category *Category `json:"category,omitempty" gorm:"foreignKey:CategoryID"`
	Creator  *User     `json:"creator,omitempty" gorm:"foreignKey:CreatorID"`
}

// AnswerRecord 答题记录模型
type AnswerRecord struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	UserID     uint      `json:"userId" gorm:"not null;index"`
	QuestionID uint      `json:"questionId" gorm:"not null;index"`
	UserAnswer int       `json:"userAnswer" gorm:"not null"`
	IsCorrect  bool      `json:"isCorrect" gorm:"not null"`
	TimeSpent  int       `json:"timeSpent" gorm:"default:0"`
	CreatedAt  time.Time `json:"createdAt"`
	
	// 关联
	User     *User     `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Question *Question `json:"question,omitempty" gorm:"foreignKey:QuestionID"`
}

// MistakeBook 错题本模型
type MistakeBook struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	UserID     uint      `json:"userId" gorm:"not null;index"`
	QuestionID uint      `json:"questionId" gorm:"not null;index"`
	IsMastered bool      `json:"isMastered" gorm:"default:false"`
	CreatedAt  time.Time `json:"createdAt"`
	
	// 关联
	User     *User     `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Question *Question `json:"question,omitempty" gorm:"foreignKey:QuestionID"`
}

// Admin 管理员模型
type Admin struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	Username     string    `json:"username" gorm:"uniqueIndex;size:50;not null"`
	PasswordHash string    `json:"-" gorm:"size:255;not null"`
	Email        string    `json:"email" gorm:"size:100"`
	Role         string    `json:"role" gorm:"type:enum('admin','editor');default:'admin'"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

// OperationLog 操作日志模型
type OperationLog struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	Operator     string    `json:"operator" gorm:"size:50;not null"`
	Action       string    `json:"action" gorm:"size:20;not null"`
	Resource     string    `json:"resource" gorm:"size:50;not null"`
	Description  string    `json:"description" gorm:"size:255"`
	IP           string    `json:"ip" gorm:"size:45"`
	UserAgent    string    `json:"userAgent" gorm:"size:500"`
	RequestData  string    `json:"requestData" gorm:"type:text"`
	ResponseData string    `json:"responseData" gorm:"type:text"`
	CreatedAt    time.Time `json:"createdAt"`
}

// TableName 设置表名
func (User) TableName() string {
	return "users"
}

func (Category) TableName() string {
	return "categories"
}

func (Question) TableName() string {
	return "questions"
}

func (AnswerRecord) TableName() string {
	return "answer_records"
}

func (MistakeBook) TableName() string {
	return "mistake_books"
}

func (Admin) TableName() string {
	return "admins"
}

func (OperationLog) TableName() string {
	return "operation_logs"
}
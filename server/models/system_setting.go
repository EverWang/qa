package models

import (
	"time"

	"gorm.io/gorm"
)

// SystemSetting 系统设置模型
type SystemSetting struct {
	ID        uint           `json:"id" gorm:"primarykey"`
	Key       string         `json:"key" gorm:"type:varchar(255);uniqueIndex;not null;comment:设置键名"`
	Value     string         `json:"value" gorm:"type:text;comment:设置值(JSON格式)"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"deletedAt" gorm:"index"`
}

// TableName 指定表名
func (SystemSetting) TableName() string {
	return "system_settings"
}
package main

import (
	"fmt"
	"log"
	"os"
	"qaminiprogram/config"
	"qaminiprogram/models"
)

// 数据库迁移程序
func main() {
	// 设置环境变量
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_PORT", "3306")
	os.Setenv("DB_USER", "root")
	os.Setenv("DB_PASSWORD", "123456")
	os.Setenv("DB_NAME", "qaminiprogram")
	
	// 初始化数据库
	config.InitDatabase()
	
	// 获取数据库实例
	db := config.GetDB()
	if db == nil {
		log.Fatal("数据库连接失败")
	}
	
	// 自动迁移模型（包含新的is_mastered字段）
	err := db.AutoMigrate(
		&models.User{},
		&models.Category{},
		&models.Question{},
		&models.AnswerRecord{},
		&models.MistakeBook{},
		&models.Admin{},
		&models.OperationLog{},
	)
	
	if err != nil {
		log.Fatal("数据库迁移失败:", err)
	}
	
	fmt.Println("数据库迁移完成，已添加is_mastered字段到mistake_books表")
}
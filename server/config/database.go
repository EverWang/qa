package config

import (
	"fmt"
	"log"
	"os"
	"qaminiprogram/models"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// InitDatabase 初始化数据库连接
func InitDatabase() {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	// 构建DSN
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user, password, host, port, dbname)

	// 连接数据库
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	log.Println("Database connected successfully")

	// 自动迁移数据库表结构
	migrateDatabase()

	// 初始化默认数据
	initializeDefaultData()
}

// GetDB 获取数据库实例
func GetDB() *gorm.DB {
	return DB
}

// migrateDatabase 自动迁移数据库表结构
func migrateDatabase() {
	err := DB.AutoMigrate(
		&models.User{},
		&models.Category{},
		&models.Question{},
		&models.AnswerRecord{},
		&models.MistakeBook{},
		&models.Admin{},
		&models.SystemSetting{},
		&models.OperationLog{},
	)
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}
	log.Println("Database migration completed")
}

// initializeDefaultData 初始化默认数据
func initializeDefaultData() {
	// 检查是否已存在管理员
	var adminCount int64
	DB.Model(&models.Admin{}).Count(&adminCount)
	if adminCount == 0 {
		// 创建默认管理员账号
		passwordHash, err := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost)
		if err != nil {
			log.Printf("Failed to hash password: %v", err)
		} else {
			admin := models.Admin{
				Username:     "admin",
				PasswordHash: string(passwordHash),
				Email:        "admin@example.com",
				Role:         "admin",
			}

			if err := DB.Create(&admin).Error; err != nil {
				log.Printf("Failed to create default admin: %v", err)
			} else {
				log.Println("Default admin user created successfully (username: admin, password: 123456)")
			}
		}
	}

	// 检查是否已存在系统设置
	var settingsCount int64
	DB.Model(&models.SystemSetting{}).Count(&settingsCount)
	if settingsCount == 0 {
		// 创建默认系统设置
		basicSettings := models.SystemSetting{
			Key:   "basic",
			Value: `{"system_name":"刷刷题","system_description":"专业的在线刷题平台","system_version":"1.0.0","contact_email":"admin@example.com","system_status":"normal","maintenance_notice":""}`,
		}
		quizSettings := models.SystemSetting{
			Key:   "quiz",
			Value: `{"daily_limit":0,"time_limit":0,"enable_points":true,"correct_points":1,"wrong_points":0,"quiz_modes":["random","category"],"show_explanation":"after_answer"}`,
		}

		if err := DB.Create(&basicSettings).Error; err != nil {
			log.Printf("Failed to create basic settings: %v", err)
		}
		if err := DB.Create(&quizSettings).Error; err != nil {
			log.Printf("Failed to create quiz settings: %v", err)
		}
		log.Println("Default system settings created successfully")
	}

	// 检查是否已存在操作日志
	var logCount int64
	DB.Model(&models.OperationLog{}).Count(&logCount)
	if logCount == 0 {
		// 创建一些测试操作日志
		testLogs := []models.OperationLog{
			{
				Operator:    "admin",
				Action:      "login",
				Resource:    "系统",
				Description: "管理员登录系统",
				IP:          "127.0.0.1",
				UserAgent:   "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36",
			},
			{
				Operator:    "admin",
				Action:      "create",
				Resource:    "题目",
				Description: "创建题目：测试题目",
				IP:          "127.0.0.1",
				UserAgent:   "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36",
			},
			{
				Operator:    "admin",
				Action:      "update",
				Resource:    "分类",
				Description: "更新分类信息",
				IP:          "127.0.0.1",
				UserAgent:   "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36",
			},
		}

		for _, testLog := range testLogs {
			if err := DB.Create(&testLog).Error; err != nil {
				log.Printf("Failed to create test operation log: %v", err)
			}
		}
		log.Println("Test operation logs created successfully")
	}
}
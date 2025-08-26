package config

import (
	"fmt"
	"log"
	"os"
	"qaminiprogram/models"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// InitDatabase 初始化数据库连接
func InitDatabase() {
	// 从环境变量获取Supabase数据库连接信息
	host := os.Getenv("SUPABASE_DB_HOST")
	port := os.Getenv("SUPABASE_DB_PORT")
	user := os.Getenv("SUPABASE_DB_USER")
	password := os.Getenv("SUPABASE_DB_PASSWORD")
	dbname := os.Getenv("SUPABASE_DB_NAME")

	// 如果环境变量未设置，使用默认的Supabase连接信息
	if host == "" {
		host = "db.lvafyknbbxmcatbmzeij.supabase.co"
	}
	if port == "" {
		port = "5432"
	}
	if user == "" {
		user = "postgres"
	}
	if password == "" {
		password = os.Getenv("SUPABASE_DB_PASSWORD") // 需要设置实际密码
	}
	if dbname == "" {
		dbname = "postgres"
	}

	// 构建PostgreSQL DSN
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=require TimeZone=Asia/Shanghai",
		host, port, user, password, dbname)

	// 连接数据库
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
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
// 注意：Supabase数据库表结构已通过SQL迁移文件创建，这里只做检查
func migrateDatabase() {
	// 检查必要的表是否存在
	if !DB.Migrator().HasTable(&models.Category{}) {
		log.Println("Warning: categories table not found, please run Supabase migrations")
	}
	if !DB.Migrator().HasTable(&models.Question{}) {
		log.Println("Warning: questions table not found, please run Supabase migrations")
	}
	if !DB.Migrator().HasTable(&models.AnswerRecord{}) {
		log.Println("Warning: answer_records table not found, please run Supabase migrations")
	}
	if !DB.Migrator().HasTable(&models.MistakeBook{}) {
		log.Println("Warning: mistake_books table not found, please run Supabase migrations")
	}
	if !DB.Migrator().HasTable(&models.Admin{}) {
		log.Println("Warning: admins table not found, please run Supabase migrations")
	}
	log.Println("Database table check completed")
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

	// 检查是否已存在admin用户在User表中
	var adminUserCount int64
	DB.Model(&models.User{}).Where("username = ?", "admin").Count(&adminUserCount)
	if adminUserCount == 0 {
		// 创建admin用户在User表中
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost)
		if err != nil {
			log.Printf("Failed to hash password: %v", err)
		} else {
			adminUser := models.User{
				Username: "admin",
				Email:    "admin@example.com",
				Password: string(hashedPassword),
				Nickname: "管理员",
				Role:     "admin",
				Status:   "active",
			}

			if err := DB.Create(&adminUser).Error; err != nil {
				log.Printf("Failed to create admin user: %v", err)
			} else {
				log.Println("Default admin user created in User table successfully (username: admin, password: 123456)")
			}
		}
	}

	// 检查是否已存在测试用户
	var userCount int64
	DB.Model(&models.User{}).Where("username = ?", "testuser").Count(&userCount)
	if userCount == 0 {
		// 创建默认测试用户
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost)
		if err != nil {
			log.Printf("Failed to hash password: %v", err)
		} else {
			testUser := models.User{
				Username: "testuser",
				Email:    "testuser@example.com",
				Password: string(hashedPassword),
				Nickname: "测试用户",
				Role:     "user",
				Status:   "active",
			}

			if err := DB.Create(&testUser).Error; err != nil {
				log.Printf("Failed to create test user: %v", err)
			} else {
				log.Println("Default test user created successfully (username: testuser, password: 123456)")
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
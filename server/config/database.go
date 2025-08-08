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
	if adminCount > 0 {
		log.Println("Admin user already exists, skipping initialization")
		return
	}

	// 创建默认管理员账号
	passwordHash, err := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Failed to hash password: %v", err)
		return
	}

	admin := models.Admin{
		Username:     "admin",
		PasswordHash: string(passwordHash),
		Email:        "admin@example.com",
		Role:         "admin",
	}

	if err := DB.Create(&admin).Error; err != nil {
		log.Printf("Failed to create default admin: %v", err)
		return
	}

	log.Println("Default admin user created successfully (username: admin, password: 123456)")
}
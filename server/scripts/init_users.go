package main

import (
	"fmt"
	"log"
	"os"
	"qaminiprogram/config"
	"qaminiprogram/models"
	"time"

	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

/**
 * 用户初始化脚本
 * 初始化默认的管理员和测试用户
 * 管理员：admin/123465
 * 测试用户：testuser/123456
 */
func main() {
	// 加载环境变量
	if err := godotenv.Load(".env"); err != nil {
		log.Printf("Warning: .env file not found: %v", err)
	}
	
	// 初始化数据库连接
	config.InitDatabase()
	db := config.GetDB()
	
	log.Println("开始初始化用户数据...")
	
	// 初始化管理员
	if err := initAdmin(db); err != nil {
		log.Printf("初始化管理员失败: %v", err)
		os.Exit(1)
	}
	
	// 初始化测试用户
	if err := initTestUser(db); err != nil {
		log.Printf("初始化测试用户失败: %v", err)
		os.Exit(1)
	}
	
	log.Println("用户数据初始化完成！")
}

/**
 * 初始化管理员账号
 * 用户名：admin
 * 密码：123465
 */
func initAdmin(db *gorm.DB) error {
	username := "admin"
	password := "123465"
	
	// 生成密码哈希
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("生成管理员密码哈希失败: %v", err)
	}
	
	// 检查管理员是否已存在
	var existingAdmin models.Admin
	result := db.Where("username = ?", username).First(&existingAdmin)
	if result.Error == nil {
		// 管理员已存在，更新密码
		if err := db.Model(&existingAdmin).Update("password_hash", string(hashedPassword)).Error; err != nil {
			return fmt.Errorf("更新管理员密码失败: %v", err)
		}
		log.Printf("管理员 %s 密码已更新", username)
		return nil
	}
	
	// 创建管理员记录
	admin := models.Admin{
		Username:     username,
		PasswordHash: string(hashedPassword),
		Email:        "admin@example.com",
		Role:         "admin",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	
	if err := db.Create(&admin).Error; err != nil {
		return fmt.Errorf("创建管理员失败: %v", err)
	}
	
	log.Printf("管理员 %s 创建成功", username)
	return nil
}

/**
 * 初始化测试用户账号
 * 用户名：testuser
 * 密码：123456
 */
func initTestUser(db *gorm.DB) error {
	username := "testuser"
	password := "123456"
	
	// 生成密码哈希
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("生成用户密码哈希失败: %v", err)
	}
	
	// 检查用户是否已存在
	var existingUser models.User
	result := db.Where("username = ?", username).First(&existingUser)
	if result.Error == nil {
		// 用户已存在，更新密码
		if err := db.Model(&existingUser).Update("password", string(hashedPassword)).Error; err != nil {
			return fmt.Errorf("更新测试用户密码失败: %v", err)
		}
		log.Printf("测试用户 %s 密码已更新", username)
		return nil
	}
	
	// 创建用户记录
	user := models.User{
		Username:   username,
		Password:   string(hashedPassword),
		Nickname:   "测试用户",
		Email:      "testuser@example.com",
		Avatar:     "",
		Role:       "user",
		Status:     "active",
		IsVerified: true,
		IsGuest:    false,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
	
	if err := db.Create(&user).Error; err != nil {
		return fmt.Errorf("创建测试用户失败: %v", err)
	}
	
	log.Printf("测试用户 %s 创建成功", username)
	return nil
}
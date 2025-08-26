package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	// 加载环境变量
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found")
	}

	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	fmt.Printf("连接参数: host=%s, port=%s, user=%s, password=%s, dbname=%s\n", host, port, user, password, dbname)

	// 首先连接到MySQL服务器（不指定数据库）
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/?charset=utf8mb4&parseTime=True&loc=Local",
		user, password, host, port)

	fmt.Println("尝试连接到MySQL服务器...")
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("无法连接到MySQL服务器:", err)
	}

	fmt.Println("成功连接到MySQL服务器")

	// 创建数据库（如果不存在）
	createDBSQL := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci", dbname)
	if err := db.Exec(createDBSQL).Error; err != nil {
		log.Fatal("创建数据库失败:", err)
	}

	fmt.Printf("数据库 %s 创建成功或已存在\n", dbname)

	// 现在连接到指定的数据库
	dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user, password, host, port, dbname)

	fmt.Println("尝试连接到指定数据库...")
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("无法连接到指定数据库:", err)
	}

	fmt.Println("成功连接到指定数据库")

	// 测试创建一个简单的表
	type TestUser struct {
		ID       uint   `gorm:"primaryKey"`
		Username string `gorm:"size:50"`
		IsGuest  bool   `gorm:"default:false"`
	}

	if err := db.AutoMigrate(&TestUser{}); err != nil {
		log.Fatal("表迁移失败:", err)
	}

	fmt.Println("测试表创建成功")

	// 测试插入数据
	testUser := TestUser{
		Username: "test_guest",
		IsGuest:  true,
	}

	if err := db.Create(&testUser).Error; err != nil {
		log.Fatal("插入测试数据失败:", err)
	}

	fmt.Printf("测试数据插入成功，ID: %d\n", testUser.ID)

	// 查询测试数据
	var foundUser TestUser
	if err := db.Where("username = ?", "test_guest").First(&foundUser).Error; err != nil {
		log.Fatal("查询测试数据失败:", err)
	}

	fmt.Printf("查询到测试数据: ID=%d, Username=%s, IsGuest=%t\n", foundUser.ID, foundUser.Username, foundUser.IsGuest)

	fmt.Println("数据库连接测试完成！")
}
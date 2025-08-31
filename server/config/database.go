package config

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

/**
 * InitDatabase 初始化MySQL数据库连接
 * 只支持MySQL数据库，移除了PostgreSQL和Supabase支持
 */
func InitDatabase() {
	// 获取MySQL数据库连接参数
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	charset := os.Getenv("DB_CHARSET")

	// 设置默认值
	if host == "" {
		host = "localhost"
	}
	if port == "" {
		port = "3306"
	}
	if user == "" {
		user = "root"
	}
	if dbname == "" {
		dbname = "shuashuati"
	}
	if charset == "" {
		charset = "utf8mb4"
	}

	// 构建MySQL DSN
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=Local",
		user, password, host, port, dbname, charset)

	// 连接MySQL数据库
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Fatal("Failed to connect to MySQL database:", err)
	}

	log.Printf("MySQL database connected successfully (Host: %s:%s, Database: %s)", host, port, dbname)

	// 自动迁移数据库表结构
	// migrateDatabase() // Removed as per user's request to rely on init.sql

	// 初始化默认数据
	// initializeDefaultData() // Removed as per user's request to rely on init.sql
}

/**
 * GetDB 获取数据库实例
 * @return *gorm.DB 数据库连接实例
 */
func GetDB() *gorm.DB {
	return DB
}




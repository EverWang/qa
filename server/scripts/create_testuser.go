package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
	_ "github.com/go-sql-driver/mysql"
)

/**
 * 创建测试用户脚本
 * 在数据库中创建testuser用户
 */
func main() {
	// 加载环境变量
	if err := godotenv.Load(".env"); err != nil {
		log.Printf("Warning: .env file not found: %v", err)
	}

	// 获取数据库连接参数
	host := getEnv("DB_HOST", "localhost")
	port := getEnv("DB_PORT", "3306")
	user := getEnv("DB_USER", "root")
	password := getEnv("DB_PASSWORD", "")
	dbname := getEnv("DB_NAME", "shuashuati")
	charset := getEnv("DB_CHARSET", "utf8mb4")

	// 构建MySQL DSN
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=Local",
		user, password, host, port, dbname, charset)

	log.Printf("连接数据库: %s:%s/%s", host, port, dbname)

	// 连接MySQL数据库
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("连接数据库失败:", err)
	}
	defer db.Close()

	// 测试连接
	if err := db.Ping(); err != nil {
		log.Fatal("数据库连接测试失败:", err)
	}

	log.Println("数据库连接成功")

	// 创建测试用户
	if err := createTestUser(db); err != nil {
		log.Printf("创建测试用户失败: %v", err)
	} else {
		log.Println("测试用户创建成功")
	}
}

/**
 * 获取环境变量，如果不存在则返回默认值
 */
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

/**
 * 创建测试用户
 */
func createTestUser(db *sql.DB) error {
	username := "testuser"
	password := "123456"

	// 生成密码哈希
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("生成密码哈希失败: %v", err)
	}

	// 插入测试用户
	_, err = db.Exec(`
		INSERT INTO users (openid, username, email, password, nickname, avatar, role, status, is_verified, is_guest, created_at, updated_at) 
		VALUES (?, ?, ?, ?, ?, '', 'user', 'active', true, false, NOW(), NOW())
	`, "testuser_openid_123456", username, "testuser@example.com", string(hashedPassword), "测试用户")
	
	if err != nil {
		return fmt.Errorf("插入测试用户失败: %v", err)
	}

	log.Printf("测试用户 %s 创建成功，密码: %s", username, password)
	return nil
}
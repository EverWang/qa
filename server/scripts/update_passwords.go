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
 * 更新用户密码脚本
 * 直接使用SQL更新管理员和测试用户的密码
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

	// 更新管理员密码
	if err := updateAdminPassword(db); err != nil {
		log.Printf("更新管理员密码失败: %v", err)
	} else {
		log.Println("管理员密码更新成功")
	}

	// 更新测试用户密码
	if err := updateTestUserPassword(db); err != nil {
		log.Printf("更新测试用户密码失败: %v", err)
	} else {
		log.Println("测试用户密码更新成功")
	}

	log.Println("密码更新完成！")
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
 * 更新管理员密码
 * 用户名：admin
 * 密码：123465
 */
func updateAdminPassword(db *sql.DB) error {
	username := "admin"
	password := "123465"

	// 生成密码哈希
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("生成管理员密码哈希失败: %v", err)
	}

	// 更新管理员密码
	_, err = db.Exec("UPDATE admins SET password_hash = ? WHERE username = ?", string(hashedPassword), username)
	if err != nil {
		return fmt.Errorf("更新管理员密码失败: %v", err)
	}

	log.Printf("管理员 %s 密码已更新为: %s", username, password)
	return nil
}

/**
 * 更新测试用户密码
 * 用户名：testuser
 * 密码：123456
 */
func updateTestUserPassword(db *sql.DB) error {
	username := "testuser"
	password := "123456"

	// 生成密码哈希
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("生成测试用户密码哈希失败: %v", err)
	}

	// 更新测试用户密码
	_, err = db.Exec("UPDATE users SET password = ? WHERE username = ?", string(hashedPassword), username)
	if err != nil {
		return fmt.Errorf("更新测试用户密码失败: %v", err)
	}

	log.Printf("测试用户 %s 密码已更新为: %s", username, password)
	return nil
}
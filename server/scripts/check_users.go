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
 * 检查用户数据脚本
 * 查询数据库中的用户信息并验证密码
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

	// 检查管理员用户
	checkAdmin(db)

	// 检查测试用户
	checkTestUser(db)
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
 * 检查管理员用户
 */
func checkAdmin(db *sql.DB) {
	log.Println("\n=== 检查管理员用户 ===")
	
	// 查询管理员
	rows, err := db.Query("SELECT id, username, password_hash, email FROM admins WHERE username = ?", "admin")
	if err != nil {
		log.Printf("查询管理员失败: %v", err)
		return
	}
	defer rows.Close()

	if !rows.Next() {
		log.Println("管理员用户不存在")
		return
	}

	var id, username, passwordHash, email string
	if err := rows.Scan(&id, &username, &passwordHash, &email); err != nil {
		log.Printf("扫描管理员数据失败: %v", err)
		return
	}

	log.Printf("管理员信息: ID=%s, Username=%s, Email=%s", id, username, email)
	log.Printf("密码哈希: %s", passwordHash)

	// 验证密码
	password := "123465"
	err = bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password))
	if err != nil {
		log.Printf("密码验证失败: %v", err)
	} else {
		log.Println("密码验证成功！")
	}
}

/**
 * 检查测试用户
 */
func checkTestUser(db *sql.DB) {
	log.Println("\n=== 检查测试用户 ===")
	
	// 查询测试用户
	rows, err := db.Query("SELECT id, username, password, email, nickname FROM users WHERE username = ?", "testuser")
	if err != nil {
		log.Printf("查询测试用户失败: %v", err)
		return
	}
	defer rows.Close()

	if !rows.Next() {
		log.Println("测试用户不存在")
		return
	}

	var id, username, password, email, nickname string
	if err := rows.Scan(&id, &username, &password, &email, &nickname); err != nil {
		log.Printf("扫描测试用户数据失败: %v", err)
		return
	}

	log.Printf("测试用户信息: ID=%s, Username=%s, Email=%s, Nickname=%s", id, username, email, nickname)
	log.Printf("密码哈希: %s", password)

	// 验证密码
	userPassword := "123456"
	err = bcrypt.CompareHashAndPassword([]byte(password), []byte(userPassword))
	if err != nil {
		log.Printf("密码验证失败: %v", err)
	} else {
		log.Println("密码验证成功！")
	}
}
package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

/**
 * 修复数据库编码问题的工具程序
 * 读取SQL修复脚本并执行，确保中文数据正确存储
 */
func main() {
	// 加载环境变量
	err := godotenv.Load(".env")
	if err != nil {
		log.Printf("Warning: .env file not found: %v", err)
	}

	// 获取数据库连接参数
	host := getEnv("DB_HOST", "localhost")
	port := getEnv("DB_PORT", "3306")
	user := getEnv("DB_USER", "root")
	password := getEnv("DB_PASSWORD", "123456")
	dbname := getEnv("DB_NAME", "qaminiprogram")

	// 构建DSN，确保使用UTF8MB4编码
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local&collation=utf8mb4_unicode_ci",
		user, password, host, port, dbname)

	fmt.Println("正在连接数据库...")

	// 连接数据库
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent), // 静默模式，减少日志输出
	})
	if err != nil {
		log.Fatal("连接数据库失败:", err)
	}

	fmt.Println("数据库连接成功")

	// 读取SQL修复脚本
	sqlContent, err := ioutil.ReadFile("fix_encoding.sql")
	if err != nil {
		log.Fatal("读取SQL文件失败:", err)
	}

	fmt.Println("开始执行数据修复...")

	// 分割SQL语句并执行
	sqlStatements := strings.Split(string(sqlContent), ";")
	for i, statement := range sqlStatements {
		statement = strings.TrimSpace(statement)
		if statement == "" || strings.HasPrefix(statement, "--") {
			continue
		}

		fmt.Printf("执行第 %d 条SQL语句...\n", i+1)
		result := db.Exec(statement)
		if result.Error != nil {
			log.Printf("执行SQL语句失败: %v\nSQL: %s", result.Error, statement)
			continue
		}
	}

	fmt.Println("数据修复完成！")

	// 验证修复结果
	fmt.Println("\n验证修复结果:")
	var categories []struct {
		ID    uint   `json:"id"`
		Name  string `json:"name"`
		Level int    `json:"level"`
		Sort  int    `json:"sort"`
	}

	db.Raw("SELECT id, name, level, sort FROM categories ORDER BY level, sort").Scan(&categories)
	fmt.Println("分类列表:")
	for _, cat := range categories {
		fmt.Printf("ID: %d, 名称: %s, 级别: %d, 排序: %d\n", cat.ID, cat.Name, cat.Level, cat.Sort)
	}

	var questionCount int64
	db.Raw("SELECT COUNT(*) FROM questions").Scan(&questionCount)
	fmt.Printf("\n总题目数量: %d\n", questionCount)
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
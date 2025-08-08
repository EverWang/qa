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
)

/**
 * 修复数据库编码问题的程序
 * 使用UPDATE方式直接修改乱码分类名称
 */
func main() {
	// 加载环境变量
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("加载.env文件失败:", err)
	}

	// 构建数据库连接字符串，确保使用utf8mb4编码
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	fmt.Println("正在连接数据库...")
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("数据库连接失败:", err)
	}
	fmt.Println("数据库连接成功")

	// 读取SQL文件
	sqlContent, err := ioutil.ReadFile("fix_encoding_update.sql")
	if err != nil {
		log.Fatal("读取SQL文件失败:", err)
	}

	// 分割SQL语句
	sqlStatements := strings.Split(string(sqlContent), ";")

	fmt.Println("开始执行数据修复...")
	for i, stmt := range sqlStatements {
		stmt = strings.TrimSpace(stmt)
		if stmt == "" || strings.HasPrefix(stmt, "--") {
			continue
		}

		fmt.Printf("执行第 %d 条SQL语句...\n", i+1)
		result := db.Exec(stmt)
		if result.Error != nil {
			log.Printf("执行SQL语句失败: %v\n", result.Error)
			log.Printf("SQL: %s\n", stmt)
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

	var totalQuestions int64
	db.Raw("SELECT COUNT(*) FROM questions").Scan(&totalQuestions)
	fmt.Printf("\n总题目数量: %d\n", totalQuestions)
}
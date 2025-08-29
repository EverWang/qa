package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
	_ "github.com/go-sql-driver/mysql"
)

/**
 * 数据库初始化程序
 * 执行init.sql脚本创建数据库表结构和初始数据
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
	charset := getEnv("DB_CHARSET", "utf8mb4")

	// 构建MySQL DSN（不指定数据库，因为SQL脚本中会创建数据库）
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/?charset=%s&parseTime=True&loc=Local&multiStatements=true",
		user, password, host, port, charset)

	log.Printf("连接数据库: %s:%s", host, port)

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

	// 读取并执行SQL脚本
	if err := executeSQLFile(db, "init_simple.sql"); err != nil {
		log.Fatal("执行SQL脚本失败:", err)
	}

	log.Println("数据库初始化完成！")
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
 * 执行SQL文件
 */
func executeSQLFile(db *sql.DB, filename string) error {
	log.Printf("读取SQL文件: %s", filename)

	// 读取SQL文件
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("打开SQL文件失败: %v", err)
	}
	defer file.Close()

	// 读取文件内容
	var sqlContent strings.Builder
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		// 跳过注释行
		if strings.HasPrefix(line, "--") || line == "" {
			continue
		}
		sqlContent.WriteString(line)
		sqlContent.WriteString("\n")
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("读取SQL文件失败: %v", err)
	}

	sqlText := sqlContent.String()
	log.Printf("SQL内容长度: %d 字符", len(sqlText))

	// 分割SQL语句（按分号分割）
	statements := strings.Split(sqlText, ";")
	log.Printf("共有 %d 条SQL语句", len(statements))

	// 执行每条SQL语句
	for i, stmt := range statements {
		stmt = strings.TrimSpace(stmt)
		if stmt == "" {
			continue
		}

		log.Printf("执行第 %d 条SQL语句...", i+1)
		_, err := db.Exec(stmt)
		if err != nil {
			// 如果是表已存在的错误，忽略它
			if strings.Contains(err.Error(), "already exists") || strings.Contains(err.Error(), "Duplicate entry") {
				log.Printf("警告: %v (忽略)", err)
				continue
			}
			return fmt.Errorf("执行SQL语句失败 (第%d条): %v\nSQL: %s", i+1, err, stmt)
		}
	}

	log.Println("SQL脚本执行完成")
	return nil
}
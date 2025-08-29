package main

import (
	"log"
	"os"
	"qaminiprogram/config"
	"qaminiprogram/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// 加载环境变量
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found")
	}

	// 设置Gin模式
	mode := os.Getenv("SERVER_MODE")
	if mode == "" {
		mode = "debug"
	}
	gin.SetMode(mode)

	// 初始化数据库
	config.InitDatabase()

	// 创建Gin引擎
	r := gin.New()

	// 设置路由
	routes.SetupRoutes(r)

	// 获取端口
	port := os.Getenv("PORT")
	if port == "" {
		port = os.Getenv("SERVER_PORT")
		if port == "" {
			port = "8080"
		}
	}

	// 启动服务器
	log.Printf("Server starting on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
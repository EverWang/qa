package handler

import (
	"net/http"
	"qaminiprogram/config"
	"qaminiprogram/routes"

	"github.com/gin-gonic/gin"
)

var engine *gin.Engine

func init() {
	// 设置Gin模式为release
	gin.SetMode(gin.ReleaseMode)

	// 初始化数据库
	config.InitDatabase()

	// 创建Gin引擎
	engine = gin.New()

	// 设置路由
	routes.SetupRoutes(engine)
}

// Handler 是Vercel的入口函数
func Handler(w http.ResponseWriter, r *http.Request) {
	engine.ServeHTTP(w, r)
}
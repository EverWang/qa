package routes

import (
	"qaminiprogram/controllers"
	"qaminiprogram/middleware"

	"github.com/gin-gonic/gin"
)

// SetupRoutes 设置路由
func SetupRoutes(r *gin.Engine) {
	// 应用中间件
	r.Use(middleware.CORS())
	r.Use(middleware.Logger())
	r.Use(gin.Recovery())

	// API路由组
	api := r.Group("/api/v1")
	{
		// 公开路由（无需认证）
		public := api.Group("/")
		{
			// 用户认证
            public.POST("/auth/login", controllers.UnifiedLogin)
			public.POST("/auth/guest", controllers.GuestLogin)
			public.POST("/auth/refresh", controllers.RefreshToken)
			
			// 分类相关（公开读取）
			public.GET("/categories", controllers.GetCategories)
			public.GET("/categories/:id", controllers.GetCategoryByID)
			
			// 题目相关（公开读取）
			public.GET("/questions", controllers.GetQuestions)
			public.GET("/questions/:id", controllers.GetQuestionByID)
			public.GET("/questions/random", controllers.GetRandomQuestions)
			public.GET("/questions/category/:categoryId", controllers.GetQuestionsByCategory)
			
			// 公开统计数据
			public.GET("/statistics/overview", controllers.GetOverviewStatistics)
			

		}

		// 需要用户认证的路由
		auth := api.Group("/")
		auth.Use(middleware.JWTAuth())
		{
			// 用户相关
			auth.GET("/user/profile", controllers.GetUserProfile)
			auth.PUT("/user/profile", controllers.UpdateUserProfile)
			
			// 答题记录
			auth.POST("/answers", controllers.SubmitAnswer)
			auth.GET("/answers/history", controllers.GetAnswerHistory)
			auth.GET("/answers/statistics", controllers.GetAnswerStatistics)
			
			// 分类进度
			auth.GET("/categories/:id/progress", controllers.GetCategoryProgress)
			
			// 错题本
			auth.GET("/mistakes", controllers.GetMistakeBooks)
			auth.POST("/mistakes", controllers.AddToMistakeBook)
			auth.DELETE("/mistakes/:id", controllers.RemoveFromMistakeBook)
			auth.DELETE("/mistakes/clear", controllers.ClearMistakeBook)
			auth.PUT("/mistakes/:questionId/master", controllers.MarkMistakeAsMastered)
			auth.PUT("/mistakes/:questionId/reset", controllers.ResetMistakeStatus)
			auth.GET("/mistakes/statistics", controllers.GetMistakeBookStatistics)
		}

		// 管理员路由
		admin := api.Group("/admin")
		{
			// 管理员认证（公开路由）
		admin.POST("/login", controllers.PasswordLogin)
		}

		// 需要管理员认证的路由
		adminAuth := api.Group("/admin")
		adminAuth.Use(middleware.JWTAuth())
		adminAuth.Use(middleware.AdminAuth())
		{
			
			
			// 用户管理
			adminAuth.GET("/users", controllers.GetUsers)
			adminAuth.GET("/users/:id", controllers.GetUserByID)
			adminAuth.POST("/users", controllers.CreateUser)
			adminAuth.PUT("/users/:id", controllers.UpdateUser)
			adminAuth.PUT("/users/:id/status", controllers.UpdateUserStatus)
			adminAuth.DELETE("/users/:id", controllers.DeleteUser)
			adminAuth.DELETE("/users/batch", controllers.BatchDeleteUsers)
			
			// 分类管理
            adminAuth.GET("/categories", controllers.GetAdminCategories)
            adminAuth.POST("/categories", controllers.CreateCategory)
            adminAuth.PUT("/categories/:id", controllers.UpdateCategory)
            adminAuth.PUT("/categories/:id/status", controllers.UpdateCategoryStatus)
            adminAuth.DELETE("/categories/:id", controllers.DeleteCategory)
			
			// 题目管理
			adminAuth.GET("/questions", controllers.GetAdminQuestions)
			adminAuth.GET("/questions/:id", controllers.GetQuestionByID)
			adminAuth.POST("/questions", controllers.CreateQuestion)
			adminAuth.PUT("/questions/:id", controllers.UpdateQuestion)
			adminAuth.DELETE("/questions/:id", controllers.DeleteQuestion)
			adminAuth.DELETE("/questions/batch", controllers.BatchDeleteQuestions)
			adminAuth.POST("/questions/import", controllers.ImportQuestions)
			adminAuth.GET("/questions/export", controllers.ExportQuestions)
			
			// 数据统计
			adminAuth.GET("/statistics/overview", controllers.GetOverviewStatistics)
			adminAuth.GET("/statistics/questions", controllers.GetQuestionStatistics)
			adminAuth.GET("/statistics/users", controllers.GetUserStatistics)
			
			// 操作日志
			adminAuth.GET("/operation-logs", controllers.GetOperationLogs)
			adminAuth.GET("/logs", controllers.GetOperationLogs)
			
			// 系统设置
			adminAuth.GET("/settings/basic", controllers.GetBasicSettings)
			adminAuth.PUT("/settings/basic", controllers.UpdateBasicSettings)
			adminAuth.GET("/settings/quiz", controllers.GetQuizSettings)
			adminAuth.PUT("/settings/quiz", controllers.UpdateQuizSettings)
			
			// 系统统计
			adminAuth.GET("/statistics", controllers.GetSystemStatistics)
			adminAuth.GET("/statistics/export", controllers.ExportStatistics)
		}
	}

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
			"message": "服务运行正常",
		})
	})
}
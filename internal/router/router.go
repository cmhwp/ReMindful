package router

import (
	"os"

	"ReMindful/internal/handler"
	"ReMindful/internal/middleware"
	"ReMindful/internal/repository"
	"ReMindful/internal/service"
	"ReMindful/pkg/utils/email"

	_ "ReMindful/docs"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
)

func InitRouter(r *gin.Engine, db *gorm.DB, rdb *redis.Client, emailSender *email.EmailSender) {
	// 全局中间件
	r.Use(middleware.Cors())
	r.Use(middleware.Logger())

	jwtSecret := os.Getenv("JWT_SECRET")

	// 初始化仓库
	userRepo := repository.NewUserRepository(db)
	learningCardsRepo := repository.NewLearningCardsRepository(db)
	tagsRepo := repository.NewTagsRepository(db)
	reviewLogsRepo := repository.NewReviewLogsRepository(db)

	// 初始化服务
	userService := service.NewUserService(userRepo, rdb, emailSender)
	learningCardsService := service.NewLearningCardsService(learningCardsRepo, rdb)
	learningCardsService.SetReviewLogsRepository(reviewLogsRepo)
	tagsService := service.NewTagsService(tagsRepo)
	reviewLogsService := service.NewReviewLogsService(reviewLogsRepo)

	// 初始化处理器
	userHandler := handler.NewUserHandler(userService)
	learningCardsHandler := handler.NewLearningCardsHandler(learningCardsService)
	tagsHandler := handler.NewTagsHandler(tagsService)
	reviewLogsHandler := handler.NewReviewLogsHandler(reviewLogsService)

	// Swagger API文档
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API 路由组
	api := r.Group("/api/v1")
	{
		// 用户相关路由
		api.POST("/send-code", userHandler.SendCode)
		api.POST("/register", userHandler.Register)
		api.POST("/login", userHandler.Login)

		// 需要认证的路由
		auth := api.Group("/")
		auth.Use(middleware.JWTAuth(jwtSecret))
		{
			// 用户信息管理
			auth.GET("/user", userHandler.GetUserInfo)
			auth.PUT("/user", userHandler.UpdateUser)

			// 学习卡片相关路由
			cards := auth.Group("/learning-cards")
			{
				cards.POST("", learningCardsHandler.CreateLearningCard)       // 创建卡片
				cards.GET("", learningCardsHandler.GetLearningCards)          // 获取卡片列表
				cards.GET("/review", learningCardsHandler.GetCardsToReview)   // 获取需要复习的卡片
				cards.GET("/:id", learningCardsHandler.GetLearningCardByID)   // 获取单个卡片
				cards.PUT("/:id", learningCardsHandler.UpdateLearningCard)    // 更新卡片
				cards.DELETE("/:id", learningCardsHandler.DeleteLearningCard) // 删除卡片
				cards.POST("/:id/review", learningCardsHandler.ReviewCard)    // 复习卡片
			}

			// 标签管理路由
			tags := auth.Group("/tags")
			{
				tags.POST("", tagsHandler.CreateTag)       // 创建标签
				tags.GET("", tagsHandler.GetTags)          // 获取标签列表
				tags.GET("/:id", tagsHandler.GetTagByID)   // 获取单个标签
				tags.PUT("/:id", tagsHandler.UpdateTag)    // 更新标签
				tags.DELETE("/:id", tagsHandler.DeleteTag) // 删除标签
			}

			// 复习日志路由
			reviewLogs := auth.Group("/review-logs")
			{
				reviewLogs.GET("", reviewLogsHandler.GetReviewLogs)                // 获取复习日志
				reviewLogs.GET("/stats", reviewLogsHandler.GetReviewStats)         // 获取复习统计
				reviewLogs.GET("/progress", reviewLogsHandler.GetLearningProgress) // 获取学习进度
				reviewLogs.GET("/heatmap", reviewLogsHandler.GetReviewHeatmap)     // 获取复习热力图
			}
		}
	}
}

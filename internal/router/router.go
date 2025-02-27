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
	userService := service.NewUserService(repository.NewUserRepository(db), rdb, emailSender)
	userHandler := handler.NewUserHandler(userService)
	learningCardsService := service.NewLearningCardsService(repository.NewLearningCardsRepository(db), rdb)
	learningCardsHandler := handler.NewLearningCardsHandler(learningCardsService)
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
			auth.GET("/user", userHandler.GetUserInfo)
			auth.PUT("/user", userHandler.UpdateUser)
			// 学习卡片相关路由
			cards := auth.Group("/learning-cards")
			{
				cards.POST("", learningCardsHandler.CreateLearningCard)     // 创建卡片
				cards.GET("/:id", learningCardsHandler.GetLearningCardByID) // 获取卡片
			}
		}
	}
}

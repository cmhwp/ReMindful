package router

import (
	"os"

	"ReMindful/internal/handler"
	"ReMindful/internal/middleware"
	"ReMindful/internal/repository"
	"ReMindful/internal/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func InitRouter(r *gin.Engine, db *gorm.DB) {
	// 全局中间件
	r.Use(middleware.Cors())
	r.Use(middleware.Logger())
	jwtSecret := os.Getenv("JWT_SECRET")
	userService := service.NewUserService(repository.NewUserRepository(db))
	userHandler := handler.NewUserHandler(userService)

	// API 路由组
	api := r.Group("/api/v1")
	{
		// 用户相关路由
		api.POST("/register", userHandler.Register)
		api.POST("/login", userHandler.Login)

		// 需要认证的路由
		auth := api.Group("/")
		auth.Use(middleware.JWTAuth(jwtSecret))
		{
			auth.GET("/user", userHandler.GetUserInfo)
			auth.PUT("/user", userHandler.UpdateUser)
		}
	}
}

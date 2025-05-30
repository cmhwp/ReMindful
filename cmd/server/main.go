// @title ReMindful API
// @version 1.0
// @description ReMindful 后端API文档
// @BasePath /api/v1
// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @schemes http https
// @accept json
// @produce json

//go:generate swag init -g cmd/server/main.go -d . -o docs

package main

import (
	"fmt"
	"log"

	"ReMindful/internal/config"
	_ "ReMindful/internal/model" // 确保模型被导入
	"ReMindful/internal/router"
	"ReMindful/pkg/database"
	"ReMindful/pkg/utils/email"

	"github.com/gin-gonic/gin"
)

func main() {
	// 加载配置文件
	cfg, err := config.LoadConfig("config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 初始化MySQL连接
	db, err := database.InitMySQL(&cfg.Database)
	if err != nil {
		log.Fatalf("Failed to initialize MySQL: %v", err)
	}

	// 自动迁移数据库表结构
	if err := database.AutoMigrate(db); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// 创建数据库索引
	if err := database.CreateIndexes(db); err != nil {
		log.Printf("Warning: Failed to create indexes: %v", err)
	}

	// 初始化Redis连接
	redis, err := database.InitRedis(&cfg.Redis)
	if err != nil {
		log.Printf("Warning: Failed to initialize Redis: %v", err)
		log.Println("Continuing without Redis...")
		redis = nil
	}
	// 程序退出时关闭数据连接
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to get SQL DB: %v", err)
	}
	defer sqlDB.Close()
	defer redis.Close()

	emailSender := email.NewEmailSender(cfg.Email)
	// 初始化Gin引擎
	if cfg.Server.Mode == "production" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}
	r := gin.New()

	// 初始化路由
	router.InitRouter(r, db, redis, emailSender)

	// 启动服务器
	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	log.Printf("Starting server on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

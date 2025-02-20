//go:build ignore

package main

import (
	"fmt"
	"log"

	"github.com/cmhwp/remindful/internal/config"
	"github.com/cmhwp/remindful/internal/router"
	"github.com/cmhwp/remindful/pkg/database"
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
	// 初始化Redis连接
	redis, err := database.InitRedis(&cfg.Redis)
	if err != nil {
		log.Fatalf("Failed to initialize Redis: %v", err)
	}
	// 程序退出时关闭数据连接
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to get SQL DB: %v", err)
	}
	defer sqlDB.Close()
	defer redis.Close()

	// 初始化Gin引擎
	if cfg.Server.Mode == "production" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}
	r := gin.New()

	// 初始化路由
	router.InitRouter(r, db)

	// 启动服务器
	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	log.Printf("Starting server on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

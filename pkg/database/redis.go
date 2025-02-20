//redis数据库操作

package database

import (
	"fmt"
	"github.com/redis/go-redis/v9"
	"ReMindful/internal/config"
	"context"
)

// 初始化Redis连接
func InitRedis(cfg *config.RedisConfig) (*redis.Client, error) {
	
	db := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	_, err := db.Ping(context.Background()).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to redis: %v", err)
	}
	return db, nil
}

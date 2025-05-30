package database

import (
	"ReMindful/internal/model"
	"fmt"

	"gorm.io/gorm"
)

// AutoMigrate 自动迁移数据库表结构
func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&model.User{},
		&model.Tag{},
		&model.LearningCard{},
		&model.ReviewLog{},
	)
}

// CreateIndexes 创建必要的索引
func CreateIndexes(db *gorm.DB) error {
	// 创建索引的辅助函数
	createIndexIfNotExists := func(indexName, tableName, columns string) error {
		// 检查索引是否已存在
		var count int64
		err := db.Raw("SELECT COUNT(*) FROM information_schema.statistics WHERE table_schema = DATABASE() AND table_name = ? AND index_name = ?", tableName, indexName).Scan(&count).Error
		if err != nil {
			return err
		}

		// 如果索引不存在，则创建
		if count == 0 {
			sql := fmt.Sprintf("CREATE INDEX %s ON %s(%s)", indexName, tableName, columns)
			return db.Exec(sql).Error
		}
		return nil
	}

	// 为学习卡片创建复合索引
	if err := createIndexIfNotExists("idx_learning_cards_user_next_review", "learning_cards", "user_id, next_review"); err != nil {
		return err
	}

	// 为复习日志创建复合索引
	if err := createIndexIfNotExists("idx_review_logs_user_time", "review_logs", "user_id, review_time"); err != nil {
		return err
	}

	// 为复习日志创建卡片ID索引
	if err := createIndexIfNotExists("idx_review_logs_card_id", "review_logs", "card_id"); err != nil {
		return err
	}

	// 为标签创建用户ID索引
	if err := createIndexIfNotExists("idx_tags_user_id", "tags", "user_id"); err != nil {
		return err
	}

	return nil
}

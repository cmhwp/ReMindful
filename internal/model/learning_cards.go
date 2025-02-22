package model

import (
	"time"

	"gorm.io/gorm"
)

// LearningCard 学习卡片模型
// @Description 学习卡片信息
type LearningCard struct {
	gorm.Model
	UserID       uint        `json:"user_id" example:"1"`                             // 用户ID
	Title        string      `json:"title" example:"Git基础知识"`                         // 标题
	Content      string      `json:"content" example:"Git是分布式版本控制系统..."`              // 内容
	CardType     CardType    `json:"card_type" example:"basic"`                       // 卡片类型
	NextReview   time.Time   `json:"next_review" example:"2024-02-22T15:04:05Z07:00"` // 下次复习时间
	LastReviewAt time.Time   `json:"last_review_at"`                                  // 上次复习时间
	ReviewCount  int         `json:"review_count" example:"0"`                        // 复习次数
	Interval     int         `json:"interval" example:"1"`                            // 当前间隔天数
	EaseFactor   float64     `json:"ease_factor" example:"2.5"`                       // 简易因子
	Difficulty   float64     `json:"difficulty" example:"0.3"`                        // 难度系数
	Tags         []Tag       `json:"tags" gorm:"many2many:card_tags;"`
	ReviewLogs   []ReviewLog `json:"review_logs"`
}

// CardType 卡片类型
type CardType string

const (
	BasicCard    CardType = "basic"    // 基础卡片
	ClozeCard    CardType = "cloze"    // 填空卡片
	QuestionCard CardType = "question" // 问答卡片
)

// CreateCardRequest 创建卡片请求
type CreateCardRequest struct {
	Title    string   `json:"title" binding:"required"`
	Content  string   `json:"content" binding:"required"`
	CardType CardType `json:"card_type" binding:"required"`
	Tags     []Tag    `json:"tags"`
}

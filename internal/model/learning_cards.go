
package model

import (
	"time"

	"gorm.io/gorm"

)

type LearningCard struct {
	gorm.Model
	UserID      uint    `json:"user_id" gorm:"index"`
	Title       string  `json:"title" gorm:"size:255"`
	Content     string  `json:"content" gorm:"type:TEXT"`
	CardType    CardType `json:"card_type" gorm:"default:0"` // 卡片类型
	NextReview  time.Time `json:"next_review" gorm:"index"`     // 下次复习时间
	Interval    int       `json:"interval" gorm:"default:0"` // 当前间隔天数
	EaseFactor  float64   `json:"ease_factor" gorm:"default:2.5"`
	Difficulty  float64   `json:"difficulty" gorm:"default:0.3"` // 系统自动计算的难度系数
	Tags        []Tag     `json:"tags" gorm:"many2many:card_tags;"`
	ReviewLogs  []ReviewLog `json:"review_logs"`
  }

  // 卡片类型
type CardType int
const (
	TextCard CardType = iota // 文本卡片
	ImageCard CardType = iota // 图片卡片
)

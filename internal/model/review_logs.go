package model

import (
	"time"

	"gorm.io/gorm"
)
type ReviewLog struct {
	gorm.Model
	CardID     uint `json:"card_id" gorm:"index"`
	UserID     uint `json:"user_id" gorm:"index"`
	ReviewTime time.Time `json:"review_time"`
	Performance int    `json:"performance" gorm:"check:performance BETWEEN 1 AND 5"` // 用户自评（1-5分）
	Duration   int     `json:"duration"` // 本次复习耗时（秒）
  }
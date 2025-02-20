package model

import (
	"gorm.io/gorm"
)

type Tag struct {
	gorm.Model
	Name      string `json:"name" gorm:"uniqueIndex;size:50"`
	ColorCode string `json:"color_code" gorm:"size:7;default:'#4CAF50'"`
	UserID    uint   `json:"user_id" gorm:"index"` // 支持用户自定义标签
  }
// 用户模型
package model

import (
	"time"

	"gorm.io/gorm"
)

// 用户模型
type User struct {
	gorm.Model
	Username     string `gorm:"uniqueIndex;size:50"`
	Email        string `gorm:"uniqueIndex;size:100"`
	PhotoURL     string `gorm:"size:255"`
	PasswordHash string `gorm:"size:255"`
	WechatOpenID string `gorm:"uniqueIndex;size:100"` // 微信登录专用
	LastLoginAt  time.Time
	IsPremium    bool   `gorm:"default:false"`
	Timezone     string `gorm:"size:50;default:'UTC'"`
}

// LoginRequest 登录请求
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse 登录响应
type LoginResponse struct {
	Token    string    `json:"token"`
	ExpireAt time.Time `json:"expire_at"`
}

// RegisterRequest 注册请求
type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// RegisterResponse 注册响应
type RegisterResponse struct {
	Message string `json:"message"`
}

// UserInfo 用户信息
type UserInfoRequest struct {
	Username string `json:"username" binding:"omitempty,min=2,max=50"`
	Email    string `json:"email" binding:"omitempty,email"`
	PhotoURL string `json:"photo_url" binding:"omitempty,url"`
}

// UserInfoResponse 用户信息响应
type UserInfoResponse struct {
	Message string `json:"message"`
	User    User   `json:"user"`
}

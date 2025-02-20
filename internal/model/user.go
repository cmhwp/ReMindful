// 用户模型
package model

import (
	"time"

	"gorm.io/gorm"
)

// User 用户模型
// @Description 用户信息
type User struct {
	gorm.Model              // 基础字段(ID, CreatedAt, UpdatedAt, DeletedAt)
	Username     string     `json:"username" gorm:"uniqueIndex;size:50"`   // 用户名
	Email        string     `json:"email" gorm:"uniqueIndex;size:100"`     // 邮箱
	PhotoURL     string     `json:"photo_url" gorm:"size:255"`             // 头像URL
	PasswordHash string     `json:"-" gorm:"size:255"`                     // 密码哈希
	WechatOpenID string     `json:"-" gorm:"uniqueIndex;size:100"`         // 微信OpenID
	LastLoginAt  *time.Time `json:"last_login_at"`                         // 最后登录时间
	IsPremium    bool       `json:"is_premium" gorm:"default:false"`       // 是否是高级用户
	Timezone     string     `json:"timezone" gorm:"size:50;default:'UTC'"` // 时区
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
// @Description 获取用户信息的响应
type UserInfoResponse struct {
	Message string `json:"message"` // 响应信息
	User    User   `json:"user"`    // 用户信息
}

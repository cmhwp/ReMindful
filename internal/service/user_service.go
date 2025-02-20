package service

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"ReMindful/internal/model"
	"ReMindful/internal/repository"
	"ReMindful/pkg/jwt"
	"ReMindful/pkg/utils/email"

	"github.com/redis/go-redis/v9"

	"gorm.io/gorm"
)

type UserService struct {
	repo        *repository.UserRepository
	redis       *redis.Client
	emailSender *email.EmailSender
}

func NewUserService(repo *repository.UserRepository, redis *redis.Client, emailSender *email.EmailSender) *UserService {
	return &UserService{
		repo:        repo,
		redis:       redis,
		emailSender: emailSender,
	}
}

// SendVerificationCode 发送验证码
func (s *UserService) SendVerificationCode(email string) error {
	// 生成6位随机验证码
	code := fmt.Sprintf("%06d", rand.Intn(1000000))

	// 将验证码存入Redis，设置5分钟过期
	key := fmt.Sprintf("verification_code:%s", email)
	err := s.redis.Set(context.Background(), key, code, 5*time.Minute).Err()
	if err != nil {
		return err
	}

	// 发送验证码邮件
	return s.emailSender.SendVerificationCode(email, code)
}

// verifyCode 验证验证码
func (s *UserService) verifyCode(email, code string) error {
	key := fmt.Sprintf("verification_code:%s", email)
	storedCode, err := s.redis.Get(context.Background(), key).Result()
	if err == redis.Nil {
		return errors.New("验证码已过期")
	}
	if err != nil {
		return err
	}
	if storedCode != code {
		return errors.New("验证码错误")
	}
	// 验证成功后删除验证码
	s.redis.Del(context.Background(), key)
	return nil
}

// Register 用户注册
func (s *UserService) Register(req *model.RegisterRequest) error {
	// 验证验证码
	if err := s.verifyCode(req.Email, req.Code); err != nil {
		return err
	}

	// 检查邮箱和用户名是否已存在
	exists, err := s.repo.CheckEmailExists(req.Email, 0)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("邮箱已被注册")
	}

	exists, err = s.repo.CheckUsernameExists(req.Username, 0)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("用户名已被使用")
	}

	// 密码加密
	hashedPassword := jwt.HashPassword(req.Password)

	// 创建用户
	user := &model.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		LastLoginAt:  &time.Time{},    // 初始化最后登录时间
		Timezone:     "Asia/Shanghai", // 中国时区
	}

	// 更新最后登录时间
	now := time.Now()
	user.LastLoginAt = &now

	return s.repo.Create(user)
}

// Login 用户登录
func (s *UserService) Login(req *model.LoginRequest, jwtSecret string) (*model.LoginResponse, error) {
	user, err := s.repo.FindByEmail(req.Email)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("用户不存在")
		}
		return nil, err
	}

	// 验证密码
	if !jwt.VerifyPassword(req.Password, user.PasswordHash) {
		return nil, errors.New("密码错误")
	}

	// 更新最后登录时间
	now := time.Now()
	user.LastLoginAt = &now
	if err = s.repo.Update(user); err != nil {
		return nil, err
	}

	// 生成JWT token
	token, err := jwt.GenerateToken(user.ID, jwtSecret, 24) // token有效期24小时
	if err != nil {
		return nil, err
	}

	return &model.LoginResponse{
		Token:    token,
		ExpireAt: time.Now().Add(24 * time.Hour),
	}, nil
}

// GetUserInfo 获取用户信息
func (s *UserService) GetUserInfo(userID uint) (*model.UserInfoResponse, error) {
	user, err := s.repo.FindByID(userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("用户不存在")
		}
		return nil, err
	}

	return &model.UserInfoResponse{
		Message: "获取用户信息成功",
		User:    *user,
	}, nil
}

// UpdateUser 更新用户信息
func (s *UserService) UpdateUser(userID uint, req *model.UserInfoRequest) error {
	updates := map[string]interface{}{}

	if req.Username != "" {
		exists, err := s.repo.CheckUsernameExists(req.Username, userID)
		if err != nil {
			return err
		}
		if exists {
			return errors.New("用户名已被使用")
		}
		updates["username"] = req.Username
	}

	if req.Email != "" {
		exists, err := s.repo.CheckEmailExists(req.Email, userID)
		if err != nil {
			return err
		}
		if exists {
			return errors.New("邮箱已被使用")
		}
		updates["email"] = req.Email
	}

	if req.PhotoURL != "" {
		updates["photo_url"] = req.PhotoURL
	}

	if len(updates) == 0 {
		return nil
	}

	return s.repo.UpdateFields(userID, updates)
}

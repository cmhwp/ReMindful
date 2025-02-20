package service

import (
	"errors"
	"time"

	"ReMindful/internal/model"
	"ReMindful/internal/repository"
	"ReMindful/pkg/jwt"

	"gorm.io/gorm"
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

// Register 用户注册
func (s *UserService) Register(req *model.RegisterRequest) error {
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
	user := model.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		Timezone:     "UTC",
	}

	return s.repo.Create(&user)
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
	user.LastLoginAt = time.Now()
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

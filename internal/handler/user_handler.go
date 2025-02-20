// 用户相关处理函数
package handler

import (
	"net/http"
	"os"

	"ReMindful/internal/model"
	"ReMindful/internal/service"
	"ReMindful/pkg/utils/response"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

// Register 用户注册
func (h *UserHandler) Register(c *gin.Context) {
	var req model.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "无效的请求参数")
		return
	}

	if err := h.userService.Register(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(c, model.RegisterResponse{Message: "注册成功"})
}

// Login 用户登录
func (h *UserHandler) Login(c *gin.Context) {
	var req model.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "无效的请求参数")
		return
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	resp, err := h.userService.Login(&req, jwtSecret)
	if err != nil {
		if err.Error() == "用户不存在" || err.Error() == "密码错误" {
			response.Error(c, http.StatusUnauthorized, err.Error())
		} else {
			response.Error(c, http.StatusInternalServerError, "服务器内部错误")
		}
		return
	}

	response.Success(c, resp)
}

// GetUserInfo 获取用户信息
func (h *UserHandler) GetUserInfo(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "未授权")
		return
	}

	info, err := h.userService.GetUserInfo(userID.(uint))
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(c, info)
}

// UpdateUser 更新用户信息
func (h *UserHandler) UpdateUser(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "未授权")
		return
	}

	var req model.UserInfoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "无效的请求参数")
		return
	}

	if err := h.userService.UpdateUser(userID.(uint), &req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(c, gin.H{"message": "更新成功"})
}

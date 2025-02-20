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

// @Summary     用户注册
// @Description 新用户注册
// @Tags        用户
// @Accept      json
// @Produce     json
// @Param       request body model.RegisterRequest true "注册信息"
// @Success     200 {object} model.RegisterResponse
// @Failure     400 {object} response.Response
// @Router      /register [post]
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

// @Summary     用户登录
// @Description 用户登录获取token
// @Tags        用户
// @Accept      json
// @Produce     json
// @Param       request body model.LoginRequest true "登录信息"
// @Success     200 {object} model.LoginResponse
// @Failure     400 {object} response.Response
// @Failure     401 {object} response.Response
// @Router      /login [post]
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

// @Summary     获取用户信息
// @Description 获取当前登录用户信息
// @Tags        用户
// @Accept      json
// @Produce     json
// @Security    Bearer
// @Success     200 {object} model.UserInfoResponse
// @Failure     401 {object} response.Response
// @Router      /user [get]
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

// @Summary     更新用户信息
// @Description 更新当前登录用户信息
// @Tags        用户
// @Accept      json
// @Produce     json
// @Security    Bearer
// @Param       request body model.UserInfoRequest true "用户信息"
// @Success     200 {object} response.Response
// @Failure     400 {object} response.Response
// @Failure     401 {object} response.Response
// @Router      /user [put]
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

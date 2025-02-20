// Package response 提供统一的HTTP响应处理
package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response 统一响应结构
// @Description 所有接口统一返回的数据结构
type Response struct {
	Code    int         `json:"code" example:"200"`        // 状态码
	Message string      `json:"message" example:"success"` // 提示信息
	Data    interface{} `json:"data,omitempty"`            // 数据
}

func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    200,
		Message: "success",
		Data:    data,
	})
}

func Error(c *gin.Context, code int, message string) {
	c.JSON(code, Response{
		Code:    code,
		Message: message,
	})
}

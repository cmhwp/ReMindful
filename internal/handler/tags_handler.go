package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"ReMindful/internal/model"
	"ReMindful/internal/service"
	"ReMindful/pkg/utils/response"
)

type TagsHandler struct {
	tagsService *service.TagsService
}

func NewTagsHandler(tagsService *service.TagsService) *TagsHandler {
	return &TagsHandler{tagsService: tagsService}
}

// CreateTagRequest 创建标签请求
type CreateTagRequest struct {
	Name      string `json:"name" binding:"required,min=1,max=50"`
	ColorCode string `json:"color_code" binding:"omitempty,len=7"`
}

// @Summary 创建标签
// @Description 创建新的学习标签
// @Tags 标签
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body CreateTagRequest true "标签信息"
// @Success 200 {object} model.Tag
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response
// @Router /tags [post]
func (h *TagsHandler) CreateTag(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "未授权")
		return
	}

	var req CreateTagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "无效的请求参数")
		return
	}

	tag := &model.Tag{
		Name:      req.Name,
		ColorCode: req.ColorCode,
		UserID:    userID.(uint),
	}

	// 如果没有提供颜色代码，使用默认值
	if tag.ColorCode == "" {
		tag.ColorCode = "#4CAF50"
	}

	if err := h.tagsService.CreateTag(tag); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, tag)
}

// @Summary 获取用户标签列表
// @Description 获取当前用户的所有标签
// @Tags 标签
// @Produce json
// @Security Bearer
// @Success 200 {object} response.Response{data=[]model.Tag}
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response
// @Router /tags [get]
func (h *TagsHandler) GetTags(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "未授权")
		return
	}

	tags, err := h.tagsService.GetTagsByUserID(userID.(uint))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, tags)
}

// @Summary 根据ID获取标签
// @Description 根据ID获取标签详情
// @Tags 标签
// @Produce json
// @Security Bearer
// @Param id path int true "标签ID"
// @Success 200 {object} model.Tag
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response "未授权"
// @Failure 403 {object} response.Response "无权限"
// @Failure 404 {object} response.Response "标签不存在"
// @Failure 500 {object} response.Response
// @Router /tags/{id} [get]
func (h *TagsHandler) GetTagByID(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "未授权")
		return
	}

	id := c.Param("id")
	idUint, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的标签ID")
		return
	}

	tag, err := h.tagsService.GetTagByID(uint(idUint))
	if err != nil {
		response.Error(c, http.StatusNotFound, "标签不存在")
		return
	}

	// 检查权限
	if tag.UserID != userID.(uint) {
		response.Error(c, http.StatusForbidden, "无权限访问此标签")
		return
	}

	response.Success(c, tag)
}

// @Summary 更新标签
// @Description 更新标签信息
// @Tags 标签
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "标签ID"
// @Param request body CreateTagRequest true "更新的标签信息"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response "未授权"
// @Failure 403 {object} response.Response "无权限"
// @Failure 404 {object} response.Response "标签不存在"
// @Failure 500 {object} response.Response
// @Router /tags/{id} [put]
func (h *TagsHandler) UpdateTag(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "未授权")
		return
	}

	id := c.Param("id")
	idUint, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的标签ID")
		return
	}

	var req CreateTagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "无效的请求参数")
		return
	}

	// 获取现有标签
	tag, err := h.tagsService.GetTagByID(uint(idUint))
	if err != nil {
		response.Error(c, http.StatusNotFound, "标签不存在")
		return
	}

	// 检查权限
	if tag.UserID != userID.(uint) {
		response.Error(c, http.StatusForbidden, "无权限操作此标签")
		return
	}

	// 更新标签信息
	tag.Name = req.Name
	if req.ColorCode != "" {
		tag.ColorCode = req.ColorCode
	}

	if err := h.tagsService.UpdateTag(tag); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, gin.H{"message": "标签更新成功"})
}

// @Summary 删除标签
// @Description 删除标签
// @Tags 标签
// @Produce json
// @Security Bearer
// @Param id path int true "标签ID"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response "未授权"
// @Failure 403 {object} response.Response "无权限"
// @Failure 404 {object} response.Response "标签不存在"
// @Failure 500 {object} response.Response
// @Router /tags/{id} [delete]
func (h *TagsHandler) DeleteTag(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "未授权")
		return
	}

	id := c.Param("id")
	idUint, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的标签ID")
		return
	}

	// 获取标签检查权限
	tag, err := h.tagsService.GetTagByID(uint(idUint))
	if err != nil {
		response.Error(c, http.StatusNotFound, "标签不存在")
		return
	}

	if tag.UserID != userID.(uint) {
		response.Error(c, http.StatusForbidden, "无权限操作此标签")
		return
	}

	if err := h.tagsService.DeleteTag(uint(idUint)); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, gin.H{"message": "标签删除成功"})
}

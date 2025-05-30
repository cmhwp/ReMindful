package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"ReMindful/internal/model"
	"ReMindful/internal/service"
	"ReMindful/pkg/utils/response"
)

type LearningCardsHandler struct {
	learningCardsService *service.LearningCardsService
}

func NewLearningCardsHandler(learningCardsService *service.LearningCardsService) *LearningCardsHandler {
	return &LearningCardsHandler{learningCardsService: learningCardsService}
}

// @Summary 创建学习卡片
// @Description 创建新的学习卡片
// @Tags 学习卡片
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body model.CreateCardRequest true "学习卡片信息"
// @Success 200 {object} model.LearningCard
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response
// @Router /learning-cards [post]
func (h *LearningCardsHandler) CreateLearningCard(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "未授权")
		return
	}

	var req model.CreateCardRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "无效的请求参数")
		return
	}

	card := &model.LearningCard{
		UserID:   userID.(uint),
		Title:    req.Title,
		Content:  req.Content,
		CardType: req.CardType,
		Tags:     req.Tags,
	}

	if err := h.learningCardsService.CreateLearningCard(card); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, card)
}

// @Summary 根据ID获取学习卡片
// @Description 根据ID获取学习卡片
// @Tags 学习卡片
// @Produce json
// @Security Bearer
// @Param id path int true "卡片ID"
// @Success 200 {object} model.LearningCard
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response
// @Router /learning-cards/{id} [get]
func (h *LearningCardsHandler) GetLearningCardByID(c *gin.Context) {
	id := c.Param("id")
	idUint, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的卡片ID")
		return
	}

	card, err := h.learningCardsService.GetLearningCardByID(uint(idUint))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, card)
}

// @Summary 获取用户的学习卡片列表
// @Description 获取当前用户的所有学习卡片，支持分页
// @Tags 学习卡片
// @Produce json
// @Security Bearer
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Param tag_id query int false "标签ID过滤"
// @Param card_type query string false "卡片类型过滤" Enums(basic,cloze,question)
// @Success 200 {object} response.Response{data=object{cards=[]model.LearningCard,total=int64,page=int,page_size=int}}
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response
// @Router /learning-cards [get]
func (h *LearningCardsHandler) GetLearningCards(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "未授权")
		return
	}

	// 解析查询参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	tagIDStr := c.Query("tag_id")
	cardTypeStr := c.Query("card_type")

	var cards []*model.LearningCard
	var total int64
	var err error

	// 根据过滤条件查询
	if tagIDStr != "" {
		tagID, _ := strconv.ParseUint(tagIDStr, 10, 64)
		cards, err = h.learningCardsService.GetLearningCardsByTag(userID.(uint), uint(tagID))
		total = int64(len(cards))
	} else if cardTypeStr != "" {
		cardType := model.CardType(cardTypeStr)
		cards, err = h.learningCardsService.GetLearningCardsByCardType(userID.(uint), cardType)
		total = int64(len(cards))
	} else {
		cards, total, err = h.learningCardsService.GetLearningCardsByUserIDWithPagination(userID.(uint), page, pageSize)
	}

	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, gin.H{
		"cards":     cards,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// @Summary 更新学习卡片
// @Description 更新学习卡片信息
// @Tags 学习卡片
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "卡片ID"
// @Param request body model.CreateCardRequest true "更新的卡片信息"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response
// @Router /learning-cards/{id} [put]
func (h *LearningCardsHandler) UpdateLearningCard(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "未授权")
		return
	}

	id := c.Param("id")
	idUint, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的卡片ID")
		return
	}

	var req model.CreateCardRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "无效的请求参数")
		return
	}

	// 获取现有卡片
	card, err := h.learningCardsService.GetLearningCardByID(uint(idUint))
	if err != nil {
		response.Error(c, http.StatusNotFound, "卡片不存在")
		return
	}

	// 检查权限
	if card.UserID != userID.(uint) {
		response.Error(c, http.StatusForbidden, "无权限操作此卡片")
		return
	}

	// 更新卡片信息
	card.Title = req.Title
	card.Content = req.Content
	card.CardType = req.CardType
	card.Tags = req.Tags

	if err := h.learningCardsService.UpdateLearningCard(card); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, gin.H{"message": "卡片更新成功"})
}

// @Summary 删除学习卡片
// @Description 软删除学习卡片
// @Tags 学习卡片
// @Produce json
// @Security Bearer
// @Param id path int true "卡片ID"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response "未授权"
// @Failure 403 {object} response.Response "无权限"
// @Failure 500 {object} response.Response
// @Router /learning-cards/{id} [delete]
func (h *LearningCardsHandler) DeleteLearningCard(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "未授权")
		return
	}

	id := c.Param("id")
	idUint, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的卡片ID")
		return
	}

	// 获取卡片检查权限
	card, err := h.learningCardsService.GetLearningCardByID(uint(idUint))
	if err != nil {
		response.Error(c, http.StatusNotFound, "卡片不存在")
		return
	}

	if card.UserID != userID.(uint) {
		response.Error(c, http.StatusForbidden, "无权限操作此卡片")
		return
	}

	if err := h.learningCardsService.SoftDelete(uint(idUint)); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, gin.H{"message": "卡片删除成功"})
}

// ReviewCardRequest 复习卡片请求
type ReviewCardRequest struct {
	Quality  int  `json:"quality" binding:"required,min=0,max=5"` // 复习质量评分 0-5
	Duration int  `json:"duration" binding:"required,min=1"`      // 复习耗时（秒）
	IsHard   bool `json:"is_hard"`                                // 是否觉得困难
}

// @Summary 复习学习卡片
// @Description 提交卡片复习结果，更新复习参数
// @Tags 学习卡片
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "卡片ID"
// @Param request body ReviewCardRequest true "复习结果"
// @Success 200 {object} response.Response{data=model.LearningCard}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response "未授权"
// @Failure 403 {object} response.Response "无权限"
// @Failure 500 {object} response.Response
// @Router /learning-cards/{id}/review [post]
func (h *LearningCardsHandler) ReviewCard(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "未授权")
		return
	}

	id := c.Param("id")
	idUint, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的卡片ID")
		return
	}

	var req ReviewCardRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "无效的请求参数")
		return
	}

	// 获取卡片
	card, err := h.learningCardsService.GetLearningCardByID(uint(idUint))
	if err != nil {
		response.Error(c, http.StatusNotFound, "卡片不存在")
		return
	}

	// 检查权限
	if card.UserID != userID.(uint) {
		response.Error(c, http.StatusForbidden, "无权限操作此卡片")
		return
	}

	// 更新复习状态
	duration := time.Duration(req.Duration) * time.Second
	if err := h.learningCardsService.UpdateCardReviewStatus(card, req.Quality, duration, req.IsHard); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, card)
}

// @Summary 获取需要复习的卡片
// @Description 获取当前用户需要复习的卡片列表
// @Tags 学习卡片
// @Produce json
// @Security Bearer
// @Success 200 {object} response.Response{data=[]model.LearningCard}
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response
// @Router /learning-cards/review [get]
func (h *LearningCardsHandler) GetCardsToReview(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "未授权")
		return
	}

	cards, err := h.learningCardsService.GetCardsToReview(userID.(uint))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, cards)
}

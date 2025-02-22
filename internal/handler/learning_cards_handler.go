package handler

import (
	"net/http"
	"strconv"

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
// @Param card body model.LearningCard true "学习卡片信息"
// @Success 200 {object} model.LearningCard
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response
// @Router /api/v1/learning-cards [post]

func (h *LearningCardsHandler) CreateLearningCard(c *gin.Context) {
	var card model.LearningCard
	if err := c.ShouldBindJSON(&card); err != nil {
		response.Error(c, http.StatusBadRequest, "无效的请求参数")
		return
	}

	if err := h.learningCardsService.CreateLearningCard(&card); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, card)
}

// @Summary 根据ID获取学习卡片
// @Description 根据ID获取学习卡片
// @Tags 学习卡片
// @Produce json
// @Param id path uint true "卡片ID"
// @Success 200 {object} model.LearningCard
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/learning-cards/{id} [get]

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

// @Summary 更新学习卡片
// @Description 更新学习卡片

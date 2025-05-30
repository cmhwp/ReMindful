package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"ReMindful/internal/service"
	"ReMindful/pkg/utils/response"
)

type ReviewLogsHandler struct {
	reviewLogsService *service.ReviewLogsService
}

func NewReviewLogsHandler(reviewLogsService *service.ReviewLogsService) *ReviewLogsHandler {
	return &ReviewLogsHandler{reviewLogsService: reviewLogsService}
}

// @Summary 获取用户复习日志
// @Description 获取当前用户的复习日志，支持分页和时间范围过滤
// @Tags 复习日志
// @Produce json
// @Security Bearer
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Param start_date query string false "开始日期 (YYYY-MM-DD)"
// @Param end_date query string false "结束日期 (YYYY-MM-DD)"
// @Param card_id query int false "卡片ID过滤"
// @Success 200 {object} response.Response{data=object{logs=[]model.ReviewLog,total=int64,page=int,page_size=int}}
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response
// @Router /review-logs [get]
func (h *ReviewLogsHandler) GetReviewLogs(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "未授权")
		return
	}

	// 解析查询参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")
	cardIDStr := c.Query("card_id")

	var startDate, endDate *time.Time
	var cardID *uint

	// 解析日期参数
	if startDateStr != "" {
		if parsed, err := time.Parse("2006-01-02", startDateStr); err == nil {
			startDate = &parsed
		}
	}
	if endDateStr != "" {
		if parsed, err := time.Parse("2006-01-02", endDateStr); err == nil {
			// 设置为当天的结束时间
			endTime := parsed.Add(24*time.Hour - time.Second)
			endDate = &endTime
		}
	}

	// 解析卡片ID
	if cardIDStr != "" {
		if id, err := strconv.ParseUint(cardIDStr, 10, 64); err == nil {
			cardIDUint := uint(id)
			cardID = &cardIDUint
		}
	}

	logs, total, err := h.reviewLogsService.GetReviewLogsByUserID(userID.(uint), page, pageSize, startDate, endDate, cardID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, gin.H{
		"logs":      logs,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// @Summary 获取复习统计
// @Description 获取用户的复习统计数据
// @Tags 复习日志
// @Produce json
// @Security Bearer
// @Param period query string false "统计周期" Enums(day,week,month,year) default(week)
// @Success 200 {object} response.Response{data=object}
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response
// @Router /review-logs/stats [get]
func (h *ReviewLogsHandler) GetReviewStats(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "未授权")
		return
	}

	period := c.DefaultQuery("period", "week")

	stats, err := h.reviewLogsService.GetReviewStats(userID.(uint), period)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, stats)
}

// @Summary 获取学习进度
// @Description 获取用户的学习进度信息
// @Tags 复习日志
// @Produce json
// @Security Bearer
// @Success 200 {object} response.Response{data=object}
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response
// @Router /review-logs/progress [get]
func (h *ReviewLogsHandler) GetLearningProgress(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "未授权")
		return
	}

	progress, err := h.reviewLogsService.GetLearningProgress(userID.(uint))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, progress)
}

// @Summary 获取复习热力图数据
// @Description 获取用户复习活动的热力图数据
// @Tags 复习日志
// @Produce json
// @Security Bearer
// @Param year query int false "年份" default(2024)
// @Success 200 {object} response.Response{data=map[string]int}
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response
// @Router /review-logs/heatmap [get]
func (h *ReviewLogsHandler) GetReviewHeatmap(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "未授权")
		return
	}

	year, _ := strconv.Atoi(c.DefaultQuery("year", "2024"))

	heatmap, err := h.reviewLogsService.GetReviewHeatmap(userID.(uint), year)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, heatmap)
}

package service

import (
	"ReMindful/internal/model"
	"ReMindful/internal/repository"
	"time"
)

type ReviewLogsService struct {
	repo *repository.ReviewLogsRepository
}

func NewReviewLogsService(repo *repository.ReviewLogsRepository) *ReviewLogsService {
	return &ReviewLogsService{
		repo: repo,
	}
}

// 创建复习日志
func (s *ReviewLogsService) CreateReviewLog(log *model.ReviewLog) error {
	return s.repo.Create(log)
}

// 根据用户ID获取复习日志（支持分页和过滤）
func (s *ReviewLogsService) GetReviewLogsByUserID(userID uint, page, pageSize int, startDate, endDate *time.Time, cardID *uint) ([]*model.ReviewLog, int64, error) {
	return s.repo.FindByUserIDWithFilters(userID, page, pageSize, startDate, endDate, cardID)
}

// 获取复习统计
func (s *ReviewLogsService) GetReviewStats(userID uint, period string) (map[string]interface{}, error) {
	var startTime time.Time
	now := time.Now()

	switch period {
	case "day":
		startTime = now.AddDate(0, 0, -1)
	case "week":
		startTime = now.AddDate(0, 0, -7)
	case "month":
		startTime = now.AddDate(0, -1, 0)
	case "year":
		startTime = now.AddDate(-1, 0, 0)
	default:
		startTime = now.AddDate(0, 0, -7) // 默认一周
	}

	// 获取基础统计
	totalReviews, err := s.repo.CountReviewsByTimeRange(userID, startTime, now)
	if err != nil {
		return nil, err
	}

	avgPerformance, err := s.repo.GetAveragePerformance(userID, startTime, now)
	if err != nil {
		return nil, err
	}

	totalDuration, err := s.repo.GetTotalDuration(userID, startTime, now)
	if err != nil {
		return nil, err
	}

	// 获取每日复习数据
	dailyStats, err := s.repo.GetDailyReviewStats(userID, startTime, now)
	if err != nil {
		return nil, err
	}

	// 获取性能分布
	performanceDistribution, err := s.repo.GetPerformanceDistribution(userID, startTime, now)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"total_reviews":            totalReviews,
		"average_performance":      avgPerformance,
		"total_duration_minutes":   totalDuration / 60, // 转换为分钟
		"daily_stats":              dailyStats,
		"performance_distribution": performanceDistribution,
		"period":                   period,
		"start_date":               startTime.Format("2006-01-02"),
		"end_date":                 now.Format("2006-01-02"),
	}, nil
}

// 获取学习进度
func (s *ReviewLogsService) GetLearningProgress(userID uint) (map[string]interface{}, error) {
	// 获取总卡片数
	totalCards, err := s.repo.GetTotalCardsByUser(userID)
	if err != nil {
		return nil, err
	}

	// 获取已复习的卡片数
	reviewedCards, err := s.repo.GetReviewedCardsByUser(userID)
	if err != nil {
		return nil, err
	}

	// 获取掌握程度良好的卡片数（平均性能 >= 4）
	masteredCards, err := s.repo.GetMasteredCardsByUser(userID)
	if err != nil {
		return nil, err
	}

	// 获取需要复习的卡片数
	dueCards, err := s.repo.GetDueCardsByUser(userID)
	if err != nil {
		return nil, err
	}

	// 计算进度百分比
	var reviewProgress, masteryProgress float64
	if totalCards > 0 {
		reviewProgress = float64(reviewedCards) / float64(totalCards) * 100
		masteryProgress = float64(masteredCards) / float64(totalCards) * 100
	}

	return map[string]interface{}{
		"total_cards":      totalCards,
		"reviewed_cards":   reviewedCards,
		"mastered_cards":   masteredCards,
		"due_cards":        dueCards,
		"review_progress":  reviewProgress,
		"mastery_progress": masteryProgress,
	}, nil
}

// 获取复习热力图数据
func (s *ReviewLogsService) GetReviewHeatmap(userID uint, year int) (map[string]int, error) {
	startDate := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(year+1, 1, 1, 0, 0, 0, 0, time.UTC)

	return s.repo.GetReviewHeatmapData(userID, startDate, endDate)
}

// 获取最近的复习活动
func (s *ReviewLogsService) GetRecentActivity(userID uint, limit int) ([]*model.ReviewLog, error) {
	return s.repo.GetRecentReviews(userID, limit)
}

// 获取复习连续天数
func (s *ReviewLogsService) GetReviewStreak(userID uint) (int, error) {
	return s.repo.GetReviewStreak(userID)
}

package repository

import (
	"ReMindful/internal/model"
	"time"

	"gorm.io/gorm"
)

type ReviewLogsRepository struct {
	db *gorm.DB
}

func NewReviewLogsRepository(db *gorm.DB) *ReviewLogsRepository {
	return &ReviewLogsRepository{db: db}
}

// 创建复习日志
func (r *ReviewLogsRepository) Create(log *model.ReviewLog) error {
	return r.db.Create(log).Error
}

// 根据用户ID查找复习日志（支持分页和过滤）
func (r *ReviewLogsRepository) FindByUserIDWithFilters(userID uint, page, pageSize int, startDate, endDate *time.Time, cardID *uint) ([]*model.ReviewLog, int64, error) {
	var logs []*model.ReviewLog
	var total int64

	query := r.db.Where("user_id = ?", userID)

	// 添加时间范围过滤
	if startDate != nil {
		query = query.Where("review_time >= ?", *startDate)
	}
	if endDate != nil {
		query = query.Where("review_time <= ?", *endDate)
	}

	// 添加卡片ID过滤
	if cardID != nil {
		query = query.Where("card_id = ?", *cardID)
	}

	// 计算总数
	if err := query.Model(&model.ReviewLog{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	err := query.Preload("Card").Order("review_time DESC").Offset(offset).Limit(pageSize).Find(&logs).Error

	return logs, total, err
}

// 统计时间范围内的复习次数
func (r *ReviewLogsRepository) CountReviewsByTimeRange(userID uint, startTime, endTime time.Time) (int64, error) {
	var count int64
	err := r.db.Model(&model.ReviewLog{}).
		Where("user_id = ? AND review_time BETWEEN ? AND ?", userID, startTime, endTime).
		Count(&count).Error
	return count, err
}

// 获取平均性能
func (r *ReviewLogsRepository) GetAveragePerformance(userID uint, startTime, endTime time.Time) (float64, error) {
	var result struct {
		AvgPerformance float64
	}
	err := r.db.Model(&model.ReviewLog{}).
		Select("AVG(performance) as avg_performance").
		Where("user_id = ? AND review_time BETWEEN ? AND ?", userID, startTime, endTime).
		Scan(&result).Error
	return result.AvgPerformance, err
}

// 获取总复习时长
func (r *ReviewLogsRepository) GetTotalDuration(userID uint, startTime, endTime time.Time) (int, error) {
	var result struct {
		TotalDuration int
	}
	err := r.db.Model(&model.ReviewLog{}).
		Select("SUM(duration) as total_duration").
		Where("user_id = ? AND review_time BETWEEN ? AND ?", userID, startTime, endTime).
		Scan(&result).Error
	return result.TotalDuration, err
}

// 获取每日复习统计
func (r *ReviewLogsRepository) GetDailyReviewStats(userID uint, startTime, endTime time.Time) ([]map[string]interface{}, error) {
	var results []struct {
		Date           string  `json:"date"`
		Count          int64   `json:"count"`
		AvgPerformance float64 `json:"avg_performance"`
	}

	err := r.db.Model(&model.ReviewLog{}).
		Select("DATE(review_time) as date, COUNT(*) as count, AVG(performance) as avg_performance").
		Where("user_id = ? AND review_time BETWEEN ? AND ?", userID, startTime, endTime).
		Group("DATE(review_time)").
		Order("date").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	// 转换为 map 格式
	stats := make([]map[string]interface{}, len(results))
	for i, result := range results {
		stats[i] = map[string]interface{}{
			"date":            result.Date,
			"count":           result.Count,
			"avg_performance": result.AvgPerformance,
		}
	}

	return stats, nil
}

// 获取性能分布
func (r *ReviewLogsRepository) GetPerformanceDistribution(userID uint, startTime, endTime time.Time) (map[int]int64, error) {
	var results []struct {
		Performance int   `json:"performance"`
		Count       int64 `json:"count"`
	}

	err := r.db.Model(&model.ReviewLog{}).
		Select("performance, COUNT(*) as count").
		Where("user_id = ? AND review_time BETWEEN ? AND ?", userID, startTime, endTime).
		Group("performance").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	distribution := make(map[int]int64)
	for _, result := range results {
		distribution[result.Performance] = result.Count
	}

	return distribution, nil
}

// 获取用户总卡片数
func (r *ReviewLogsRepository) GetTotalCardsByUser(userID uint) (int64, error) {
	var count int64
	err := r.db.Model(&model.LearningCard{}).
		Where("user_id = ?", userID).
		Count(&count).Error
	return count, err
}

// 获取已复习的卡片数
func (r *ReviewLogsRepository) GetReviewedCardsByUser(userID uint) (int64, error) {
	var count int64
	err := r.db.Model(&model.LearningCard{}).
		Where("user_id = ? AND review_count > 0", userID).
		Count(&count).Error
	return count, err
}

// 获取掌握程度良好的卡片数
func (r *ReviewLogsRepository) GetMasteredCardsByUser(userID uint) (int64, error) {
	var count int64

	// 子查询：获取每个卡片的平均性能
	subQuery := r.db.Model(&model.ReviewLog{}).
		Select("card_id, AVG(performance) as avg_performance").
		Where("user_id = ?", userID).
		Group("card_id").
		Having("AVG(performance) >= ?", 4)

	err := r.db.Model(&model.LearningCard{}).
		Joins("JOIN (?) as avg_perf ON learning_cards.id = avg_perf.card_id", subQuery).
		Where("learning_cards.user_id = ?", userID).
		Count(&count).Error

	return count, err
}

// 获取需要复习的卡片数
func (r *ReviewLogsRepository) GetDueCardsByUser(userID uint) (int64, error) {
	var count int64
	now := time.Now()
	err := r.db.Model(&model.LearningCard{}).
		Where("user_id = ? AND next_review <= ?", userID, now).
		Count(&count).Error
	return count, err
}

// 获取复习热力图数据
func (r *ReviewLogsRepository) GetReviewHeatmapData(userID uint, startDate, endDate time.Time) (map[string]int, error) {
	var results []struct {
		Date  string `json:"date"`
		Count int    `json:"count"`
	}

	err := r.db.Model(&model.ReviewLog{}).
		Select("DATE(review_time) as date, COUNT(*) as count").
		Where("user_id = ? AND review_time BETWEEN ? AND ?", userID, startDate, endDate).
		Group("DATE(review_time)").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	heatmap := make(map[string]int)
	for _, result := range results {
		heatmap[result.Date] = result.Count
	}

	return heatmap, nil
}

// 获取最近的复习记录
func (r *ReviewLogsRepository) GetRecentReviews(userID uint, limit int) ([]*model.ReviewLog, error) {
	var logs []*model.ReviewLog
	err := r.db.Where("user_id = ?", userID).
		Preload("Card").
		Order("review_time DESC").
		Limit(limit).
		Find(&logs).Error
	return logs, err
}

// 获取复习连续天数
func (r *ReviewLogsRepository) GetReviewStreak(userID uint) (int, error) {
	var dates []string
	err := r.db.Model(&model.ReviewLog{}).
		Select("DISTINCT DATE(review_time) as date").
		Where("user_id = ?", userID).
		Order("date DESC").
		Pluck("date", &dates).Error

	if err != nil {
		return 0, err
	}

	if len(dates) == 0 {
		return 0, nil
	}

	// 计算连续天数
	streak := 0
	today := time.Now().Format("2006-01-02")
	yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")

	// 检查今天或昨天是否有复习记录
	if len(dates) > 0 && (dates[0] == today || dates[0] == yesterday) {
		streak = 1

		for i := 1; i < len(dates); i++ {
			expectedDate := time.Now().AddDate(0, 0, -i).Format("2006-01-02")
			if dates[i] == expectedDate {
				streak++
			} else {
				break
			}
		}
	}

	return streak, nil
}

package service

import (
	"ReMindful/internal/model"
	"ReMindful/internal/repository"
	"ReMindful/pkg/algorithm"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type LearningCardsService struct {
	repo           *repository.LearningCardsRepository
	reviewLogsRepo *repository.ReviewLogsRepository
	redis          *redis.Client
}

func NewLearningCardsService(repo *repository.LearningCardsRepository, redis *redis.Client) *LearningCardsService {
	return &LearningCardsService{
		repo:  repo,
		redis: redis,
	}
}

// 设置复习日志仓库
func (s *LearningCardsService) SetReviewLogsRepository(reviewLogsRepo *repository.ReviewLogsRepository) {
	s.reviewLogsRepo = reviewLogsRepo
}

// 缓存相关的常量
const (
	cardCacheKeyPrefix  = "learning_card:"
	cardCacheExpiration = 24 * time.Hour
)

// 创建学习卡片
func (s *LearningCardsService) CreateLearningCard(card *model.LearningCard) error {
	// 校验参数
	if card.UserID == 0 {
		return errors.New("用户ID不能为0")
	}

	// 设置初始复习时间
	now := time.Now()
	card.LastReviewAt = now
	card.NextReview = now.Add(24 * time.Hour) // 默认24小时后复习
	card.ReviewCount = 0
	card.EaseFactor = 2.5 // 默认简易因子
	card.Difficulty = 0.3 // 默认难度系数

	// 创建学习卡片
	return s.repo.Create(card)
}

// 根据ID查找学习卡片
func (s *LearningCardsService) GetLearningCardByID(id uint) (*model.LearningCard, error) {
	// 先从缓存获取
	if s.redis != nil {
		card, err := s.getCardFromCache(id)
		if err == nil {
			return card, nil
		}
	}

	// 缓存未命中，从数据库获取
	card, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	// 写入缓存
	if s.redis != nil {
		s.setCardToCache(card)
	}

	return card, nil
}

// 从缓存获取卡片
func (s *LearningCardsService) getCardFromCache(id uint) (*model.LearningCard, error) {
	key := fmt.Sprintf("%s%d", cardCacheKeyPrefix, id)
	data, err := s.redis.Get(context.Background(), key).Bytes()
	if err != nil {
		return nil, err
	}

	var card model.LearningCard
	if err := json.Unmarshal(data, &card); err != nil {
		return nil, err
	}
	return &card, nil
}

// 将卡片写入缓存
func (s *LearningCardsService) setCardToCache(card *model.LearningCard) error {
	data, err := json.Marshal(card)
	if err != nil {
		return err
	}

	key := fmt.Sprintf("%s%d", cardCacheKeyPrefix, card.ID)
	return s.redis.Set(context.Background(), key, data, cardCacheExpiration).Err()
}

// 获取需要复习的卡片
func (s *LearningCardsService) GetCardsToReview(userID uint) ([]*model.LearningCard, error) {
	now := time.Now()
	return s.repo.FindByReviewTimeRange(userID, now.Add(-24*time.Hour), now)
}

// 更新卡片复习状态
func (s *LearningCardsService) UpdateCardReviewStatus(card *model.LearningCard, quality int, duration time.Duration, isHard bool) error {
	// 计算复习质量
	if quality < 0 {
		quality = algorithm.GetReviewQuality(duration, quality > 2, isHard)
	}

	// 应用 SM-2 算法
	params := algorithm.CalculateNextReview(&algorithm.SM2Parameters{
		ReviewCount: card.ReviewCount,
		Difficulty:  card.Difficulty,
		LastReview:  card.LastReviewAt,
		NextReview:  card.NextReview,
	}, quality)

	// 更新卡片
	card.ReviewCount = params.ReviewCount
	card.Difficulty = params.Difficulty
	card.LastReviewAt = params.LastReview
	card.NextReview = params.NextReview

	// 创建复习日志
	if s.reviewLogsRepo != nil {
		reviewLog := &model.ReviewLog{
			CardID:      card.ID,
			UserID:      card.UserID,
			ReviewTime:  time.Now(),
			Performance: quality,
			Duration:    int(duration.Seconds()),
		}

		// 保存复习日志（忽略错误，不影响主流程）
		s.reviewLogsRepo.Create(reviewLog)
	}

	// 更新数据库
	if err := s.repo.UpdateLearningCard(card); err != nil {
		return err
	}

	// 更新缓存
	if s.redis != nil {
		s.setCardToCache(card)
	}

	return nil
}

// 根据用户ID查找学习卡片
func (s *LearningCardsService) GetLearningCardsByUserID(userID uint) ([]*model.LearningCard, error) {
	return s.repo.FindByUserID(userID)
}

// 更新学习卡片
func (s *LearningCardsService) UpdateLearningCard(card *model.LearningCard) error {
	// 清除缓存
	if s.redis != nil {
		key := fmt.Sprintf("%s%d", cardCacheKeyPrefix, card.ID)
		s.redis.Del(context.Background(), key)
	}

	return s.repo.UpdateLearningCard(card)
}

// 根据标签过滤查询
func (s *LearningCardsService) GetLearningCardsByTag(userID uint, tagID uint) ([]*model.LearningCard, error) {
	return s.repo.FindByTag(userID, tagID)
}

// 根据难度范围查询
func (s *LearningCardsService) GetLearningCardsByDifficultyRange(userID uint, min, max float64) ([]*model.LearningCard, error) {
	return s.repo.FindByDifficultyRange(userID, min, max)
}

// 根据卡片类型查询
func (s *LearningCardsService) GetLearningCardsByCardType(userID uint, cardType model.CardType) ([]*model.LearningCard, error) {
	return s.repo.FindByCardType(userID, cardType)
}

// 分页查询
func (s *LearningCardsService) GetLearningCardsByUserIDWithPagination(userID uint, page, pageSize int) ([]*model.LearningCard, int64, error) {
	return s.repo.FindByUserIDWithPagination(userID, page, pageSize)
}

// 根据复习时间范围查询
func (s *LearningCardsService) GetLearningCardsByReviewTimeRange(userID uint, start, end time.Time) ([]*model.LearningCard, error) {
	return s.repo.FindByReviewTimeRange(userID, start, end)
}

// 批量更新复习参数
func (s *LearningCardsService) UpdateReviewParameters(cards []*model.LearningCard) error {
	return s.repo.UpdateReviewParameters(cards)
}

// 软删除
func (s *LearningCardsService) SoftDelete(id uint) error {
	// 清除缓存
	if s.redis != nil {
		key := fmt.Sprintf("%s%d", cardCacheKeyPrefix, id)
		s.redis.Del(context.Background(), key)
	}

	return s.repo.SoftDelete(id)
}

package repository

import (
	"ReMindful/internal/model"
	"time"
	"gorm.io/gorm"
)

type LearningCardsRepository struct {
	db *gorm.DB
}

func NewLearningCardsRepository(db *gorm.DB) *LearningCardsRepository {
	return &LearningCardsRepository{db: db}
}

// Create 创建学习卡片
func (r *LearningCardsRepository) Create(learningCard *model.LearningCard) error {
	return r.db.Create(learningCard).Error
}

// FindByID 通过ID查找学习卡片
func (r *LearningCardsRepository) FindByID(id uint) (*model.LearningCard, error) {
	var learningCard model.LearningCard
	err := r.db.Where("id = ?", id).First(&learningCard).Error
	if err != nil {
		return nil, err
	}
	return &learningCard, nil
}

// FindByUserID 通过用户ID查找学习卡片	
func (r *LearningCardsRepository) FindByUserID(userID uint) ([]*model.LearningCard, error) {
	var learningCards []*model.LearningCard
	err := r.db.Where("user_id = ?", userID).Find(&learningCards).Error
	if err != nil {
		return nil, err
	}
	return learningCards, nil
}

// 更新学习卡片
func (r *LearningCardsRepository) UpdateLearningCard(learningCard *model.LearningCard) error {
	return r.db.Model(&model.LearningCard{}).Where("id = ?", learningCard.ID).Updates(learningCard).Error
}


// 页查询方法（带预加载关联）
func (r *LearningCardsRepository) FindByUserIDWithPagination(userID uint, page, pageSize int) ([]*model.LearningCard, int64, error) {
	var learningCards []*model.LearningCard
	var total int64

	query := r.db.Model(&model.LearningCard{}).
		Where("user_id = ?", userID).
		Preload("Tags").
		Preload("ReviewLogs")

	// 计算总条数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	if err := query.Offset((page - 1) * pageSize).
	Limit(pageSize).
	Order("next_review DESC").
	Find(&learningCards).Error; err != nil {
		return nil, 0, err
	}

	return learningCards, total, nil
}

// 根据复习时间范围查询（用于间隔重复算法）
func (r *LearningCardsRepository) FindByReviewTimeRange(userID uint, start, end time.Time) ([]*model.LearningCard, error) {
	var cards []*model.LearningCard
	err := r.db.Model(&model.LearningCard{}).
	Where("user_id = ?", userID).
	Where("next_review BETWEEN ? AND ?", start, end).
	Preload("ReviewLogs").	
	Find(&cards).Error
	if err != nil {
		return nil, err
	}
	return cards, nil	
}

//批量更新复习参数（间隔重复的核心操作）
func (r *LearningCardsRepository) UpdateReviewParameters(cards []*model.LearningCard) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
        for _, card := range cards {
            if err := tx.Model(card).
                Updates(map[string]interface{}{
                    "next_review": card.NextReview,
                    "interval":    card.Interval,
                    "ease_factor": card.EaseFactor,
                    "difficulty":  card.Difficulty,
                }).Error; err != nil {
                return err
            }
        }
        return nil
    })
}
//软删除
func (r *LearningCardsRepository) SoftDelete(id uint) error {
	return r.db.Model(&model.LearningCard{}).Where("id = ?", id).Update("deleted_at", time.Now()).Error
}

// 根据标签过滤查询
func (r *LearningCardsRepository) FindByTag(userID uint, tagID uint) ([]*model.LearningCard, error) {
	var cards []*model.LearningCard
	err := r.db.Joins("JOIN card_tags ON card_tags.learning_card_id = learning_cards.id").
	Where("learning_cards.user_id = ? AND card_tags.tag_id = ?", userID, tagID).
	Preload("Tags").
	Find(&cards).Error
	if err != nil {
		return nil, err
	}
	return cards, nil
}

// 根据难度范围查询
func (r *LearningCardsRepository) FindByDifficultyRange(userID uint, min, max float64) ([]*model.LearningCard, error) {
	var cards []*model.LearningCard
	err := r.db.Model(&model.LearningCard{}).
	Where("user_id = ? AND difficulty BETWEEN ? AND ?", userID, min, max).	
	Preload("Tags").
	Preload("ReviewLogs").
	Find(&cards).Error
	if err != nil {
		return nil, err
	}
	return cards, nil
}

// 根据卡片类型查询
func (r *LearningCardsRepository) FindByCardType(userID uint, cardType model.CardType) ([]*model.LearningCard, error) {
	var cards []*model.LearningCard
	err := r.db.Model(&model.LearningCard{}).
	Where("user_id = ? AND card_type = ?", userID, cardType).
	Preload("Tags").
	Preload("ReviewLogs").
	Find(&cards).Error
	if err != nil {
		return nil, err
	}
	return cards, nil
}	

package repository

import (
	"ReMindful/internal/model"

	"gorm.io/gorm"
)

type TagsRepository struct {
	db *gorm.DB
}

func NewTagsRepository(db *gorm.DB) *TagsRepository {
	return &TagsRepository{db: db}
}

// 创建标签
func (r *TagsRepository) Create(tag *model.Tag) error {
	return r.db.Create(tag).Error
}

// 根据ID查找标签
func (r *TagsRepository) FindByID(id uint) (*model.Tag, error) {
	var tag model.Tag
	err := r.db.First(&tag, id).Error
	if err != nil {
		return nil, err
	}
	return &tag, nil
}

// 根据用户ID查找标签列表
func (r *TagsRepository) FindByUserID(userID uint) ([]*model.Tag, error) {
	var tags []*model.Tag
	err := r.db.Where("user_id = ?", userID).Find(&tags).Error
	return tags, err
}

// 根据名称和用户ID查找标签
func (r *TagsRepository) FindByNameAndUserID(name string, userID uint) (*model.Tag, error) {
	var tag model.Tag
	err := r.db.Where("name = ? AND user_id = ?", name, userID).First(&tag).Error
	if err != nil {
		return nil, err
	}
	return &tag, nil
}

// 更新标签
func (r *TagsRepository) Update(tag *model.Tag) error {
	return r.db.Save(tag).Error
}

// 删除标签
func (r *TagsRepository) Delete(id uint) error {
	return r.db.Delete(&model.Tag{}, id).Error
}

// 根据名称搜索标签
func (r *TagsRepository) SearchByName(userID uint, name string) ([]*model.Tag, error) {
	var tags []*model.Tag
	err := r.db.Where("user_id = ? AND name LIKE ?", userID, "%"+name+"%").Find(&tags).Error
	return tags, err
}

// 获取标签使用统计
func (r *TagsRepository) GetTagUsageStats(userID uint) (map[uint]int64, error) {
	var results []struct {
		TagID uint  `json:"tag_id"`
		Count int64 `json:"count"`
	}

	err := r.db.Table("card_tags").
		Select("tag_id, COUNT(*) as count").
		Joins("JOIN tags ON tags.id = card_tags.tag_id").
		Where("tags.user_id = ?", userID).
		Group("tag_id").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	stats := make(map[uint]int64)
	for _, result := range results {
		stats[result.TagID] = result.Count
	}

	return stats, nil
}

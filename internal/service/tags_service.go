package service

import (
	"ReMindful/internal/model"
	"ReMindful/internal/repository"
	"errors"
)

type TagsService struct {
	repo *repository.TagsRepository
}

func NewTagsService(repo *repository.TagsRepository) *TagsService {
	return &TagsService{
		repo: repo,
	}
}

// 创建标签
func (s *TagsService) CreateTag(tag *model.Tag) error {
	// 校验参数
	if tag.UserID == 0 {
		return errors.New("用户ID不能为0")
	}
	if tag.Name == "" {
		return errors.New("标签名称不能为空")
	}

	// 检查同一用户下是否已存在相同名称的标签
	existingTag, err := s.repo.FindByNameAndUserID(tag.Name, tag.UserID)
	if err == nil && existingTag != nil {
		return errors.New("标签名称已存在")
	}

	return s.repo.Create(tag)
}

// 根据ID获取标签
func (s *TagsService) GetTagByID(id uint) (*model.Tag, error) {
	return s.repo.FindByID(id)
}

// 根据用户ID获取标签列表
func (s *TagsService) GetTagsByUserID(userID uint) ([]*model.Tag, error) {
	return s.repo.FindByUserID(userID)
}

// 更新标签
func (s *TagsService) UpdateTag(tag *model.Tag) error {
	if tag.Name == "" {
		return errors.New("标签名称不能为空")
	}

	// 检查同一用户下是否已存在相同名称的标签（排除当前标签）
	existingTag, err := s.repo.FindByNameAndUserID(tag.Name, tag.UserID)
	if err == nil && existingTag != nil && existingTag.ID != tag.ID {
		return errors.New("标签名称已存在")
	}

	return s.repo.Update(tag)
}

// 删除标签
func (s *TagsService) DeleteTag(id uint) error {
	return s.repo.Delete(id)
}

// 根据名称搜索标签
func (s *TagsService) SearchTagsByName(userID uint, name string) ([]*model.Tag, error) {
	return s.repo.SearchByName(userID, name)
}

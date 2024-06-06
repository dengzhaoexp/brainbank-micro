package dao

import (
	"chat/repositry/model"
	"context"
	"gorm.io/gorm"
)

type ResourceDao struct {
	*gorm.DB
}

func NewResourceDao(ctx context.Context) *ResourceDao {
	return &ResourceDao{
		_db,
	}
}

func (rd *ResourceDao) CreateConversation(cs model.Conversation) (err error) {
	err = rd.DB.Model(&model.Conversation{}).Create(&cs).Error
	return
}

func (rd *ResourceDao) GetConversationsByUserId(userId string) ([]*model.Conversation, error) {
	var cs []*model.Conversation
	err := rd.DB.Where("created_by = ?", userId).Find(&cs).Error
	return cs, err
}

func (rd *ResourceDao) GetConversationsByUserIdWithOffset(userId string, offset, limit uint32) ([]*model.Conversation, error) {
	var cs []*model.Conversation
	err := rd.DB.Model(&model.Conversation{}).Where("created_by = ?", userId).
		Offset(int(offset)).
		Limit(int(limit)).
		Find(&cs).Error
	if err != nil {
		return nil, err
	}
	return cs, nil
}

func (rd *ResourceDao) GetConversationByID(id string) (c *model.Conversation, err error) {
	err = rd.DB.Model(&model.Conversation{}).
		Where("conversation_id = ?", id).
		First(&c).Error
	return
}

func (rd *ResourceDao) UpdateConversation(c model.Conversation) (err error) {
	err = rd.DB.Model(&c).Where("id = ? AND conversation_id = ?", c.ID, c.ConversationId).Updates(&c).Error
	return
}

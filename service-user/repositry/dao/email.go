package dao

import (
	"context"
	"gorm.io/gorm"
	"user/repositry/model"
)

type MailDao struct {
	*gorm.DB
}

func NewMailDao(ctx context.Context) *UserDao {
	return &UserDao{_db}
}

func NewMailDaoByDb(db *gorm.DB) *MailDao {
	return &MailDao{
		db,
	}
}

func (m *MailDao) GetResource(id string) (d *model.Mail, err error) {
	err = m.DB.Model(&model.Mail{}).Where("id = ?", id).Find(&d).Error
	if err != nil {
		return nil, err
	}
	return d, nil
}

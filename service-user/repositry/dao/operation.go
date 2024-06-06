package dao

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"user/repositry/model"
)

type UserDao struct {
	*gorm.DB
}

func NewUserDao(ctx context.Context) *UserDao {
	return &UserDao{_db.WithContext(ctx)}
}

func NewUSerDaoByDb(db *gorm.DB) *UserDao {
	return &UserDao{
		db,
	}
}

func (ud *UserDao) IsEmailExists(email string) (bool, error) {
	var user model.User
	err := ud.DB.Model(&model.User{}).Where("email_address = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 该用户还未注册
			return false, nil
		}
		// 查询过程中发生其他错误
		return false, err
	}

	// 检查 user 是否为空结构体
	if user == (model.User{}) {
		// 该用户还未注册
		return false, nil
	}

	// 该用户已经注册
	return true, nil
}

func (ud *UserDao) AddUser(user *model.User) error {
	err := ud.DB.Model(&model.User{}).Create(user).Error
	return err
}

func (ud *UserDao) UpdateUser(user *model.User) error {
	err := ud.DB.Model(&model.User{}).Where("id = ?", user.ID).Updates(user).Error
	return err
}

func (ud *UserDao) GetUserByEmail(email string) (user *model.User, err error) {
	err = ud.DB.Model(&model.User{}).Where("email_address = ?", email).First(&user).Error
	return
}

func (ud *UserDao) GetUserByUserID(userID string) (user *model.User, err error) {
	err = ud.DB.Model(&model.User{}).Where("user_id = ?", userID).First(&user).Error
	return
}

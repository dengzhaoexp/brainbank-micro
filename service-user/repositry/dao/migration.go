package dao

import (
	"user/repositry/model"
)

func migration() error {
	err := _db.Set("gorm:table_options", "charset=utf8mb4").
		AutoMigrate(&model.User{}, &model.Mail{})
	if err != nil {
		return err
	}
	return nil
}

package dao

import "chat/repositry/model"

func migration() {
	err := _db.Set("gorm:table_options", "charset=utf8mb4").
		AutoMigrate(
			&model.Conversation{},
		)
	if err != nil {
		panic(err)
	}
}

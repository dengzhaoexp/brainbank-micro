package dao

import (
	"file/repositry/model"
)

func migration() error {
	err := _db.Set("gorm:table_options", "charset=utf8mb4").
		AutoMigrate(
			&model.Neuron{},
			&model.NeuronConfig{},
			&model.LLMParam{},
			&model.Permission{},
			&model.Document{},
		)
	if err != nil {
		return err
	}

	return nil
}

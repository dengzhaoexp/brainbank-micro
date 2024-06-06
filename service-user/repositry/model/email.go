package model

type Mail struct {
	ID      string `gorm:"unique"`
	Setting string
}

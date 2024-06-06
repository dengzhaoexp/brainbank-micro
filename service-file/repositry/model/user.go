package model

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	gorm.Model
	UserId                  int64  `gorm:"unique"`
	UserName                string `gorm:"unique"`
	Nickname                string
	EmailAddress            string
	AccountPassword         string
	LastLoginTime           *time.Time
	Identity                string // The user's role or identity within the system (e.g., admin, regular user).
	Status                  string // The account status (e.g., active, suspended, banned).
	Avatar                  string
	UsedSpace               int64
	TotalSpace              int64
	CreateTime              int64
	PhoneNumber             uint32
	PreferredLanguage       string // The language preference for the user.
	NotificationPreferences string // User-specific notification settings.
	UserGroups              string // A list of groups the user belongs to (useful for permission management).
}

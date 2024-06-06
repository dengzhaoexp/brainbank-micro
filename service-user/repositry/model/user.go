package model

import (
	"crypto/md5"
	"fmt"
	"gorm.io/gorm"
	"reflect"
	"time"
)

type User struct {
	gorm.Model
	UserId                  string `gorm:"unique"`
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

func (u *User) SetPassword(password string) (err error) {
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(password))
	cipherStr := md5Ctx.Sum(nil)
	u.AccountPassword = fmt.Sprintf("%x", cipherStr)
	return nil
}

func (u *User) CheckPassword(password string) bool {
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(password))
	cipherStr := md5Ctx.Sum(nil)
	md5Str := fmt.Sprintf("%x", cipherStr)
	return reflect.DeepEqual(md5Str, u.AccountPassword)
}

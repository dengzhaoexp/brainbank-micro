package dao

import (
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"strings"
	"time"
	"user/config"
	userLogger "user/pkg/utils/logger"
)

var _db *gorm.DB

func InitMysql() {
	dbConfig := config.Config.Mysql
	dsn := strings.Join([]string{dbConfig.UserName, ":", dbConfig.Password, "@tcp(", dbConfig.DbHost, ":", dbConfig.DbPort,
		")/", dbConfig.DbName, "?charset=", dbConfig.Charset, "&parseTime=True"}, "")
	userLogger.LogrusObj.Info("Load mysql configuration successfully.")

	var ormLogger = logger.Default
	if gin.Mode() == "debug" {
		ormLogger = logger.Default.LogMode(logger.Info)
	}

	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       dsn,
		DefaultStringSize:         256,
		DisableDatetimePrecision:  true,
		DontSupportRenameIndex:    true,
		DontSupportRenameColumn:   true,
		SkipInitializeWithVersion: false,
	}), &gorm.Config{
		Logger: ormLogger,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		userLogger.LogrusObj.Error("Error opening mysql connection:", err)
		panic(err)
	}
	userLogger.LogrusObj.Info("Open mysql successfully.")

	sqlDb, _ := db.DB()
	sqlDb.SetMaxIdleConns(20)  // 设置最大连接池
	sqlDb.SetMaxOpenConns(100) // 设置最大打开数
	sqlDb.SetConnMaxLifetime(time.Second * 30)
	_db = db
	userLogger.LogrusObj.Info("The basic setup of mysql was successful.")

	if err = migration(); err != nil {
		userLogger.LogrusObj.Error("Error occurs in model migration to mysql database:", err)
		panic(err)
	}
	userLogger.LogrusObj.Info("File model migration to mysql database successful")
}

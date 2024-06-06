package dao

import (
	"file/config"
	fileLogger "file/pkg/utils/logger"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"strings"
	"time"
)

var _db *gorm.DB

func InitMysql() {
	dbConfig := config.Config.Mysql
	dsn := strings.Join([]string{dbConfig.UserName, ":", dbConfig.Password, "@tcp(", dbConfig.DbHost, ":", dbConfig.DbPort,
		")/", dbConfig.DbName, "?charset=", dbConfig.Charset, "&parseTime=True"}, "")
	fileLogger.LogrusObj.Info("Load mysql configuration successfully.")

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
		fileLogger.LogrusObj.Error("Error opening mysql connection:", err)
		panic(err)
	}
	fileLogger.LogrusObj.Info("Open mysql successfully.")

	sqlDb, _ := db.DB()
	sqlDb.SetMaxIdleConns(20)  // 设置最大连接池
	sqlDb.SetMaxOpenConns(100) // 设置最大打开数
	sqlDb.SetConnMaxLifetime(time.Second * 30)
	_db = db
	fileLogger.LogrusObj.Info("The basic setup of mysql was successful.")

	if err = migration(); err != nil {
		fileLogger.LogrusObj.Error("Error occurs in model migration to mysql database:", err)
		panic(err)
	}
	fileLogger.LogrusObj.Info("File model migration to mysql database successful")
}

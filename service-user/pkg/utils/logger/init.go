package logger

import (
	"github.com/sirupsen/logrus"
	"os"
	"path"
	"time"
)

var LogrusObj *logrus.Logger

func InitLog() {
	if LogrusObj != nil {
		src, _ := setOutputFile()
		LogrusObj.Out = src
		return
	}
	// 日志对象实例化
	src, _ := setOutputFile()
	logger := logrus.New()
	logger.Out = src
	// 设置日志级别
	logger.SetLevel(logrus.DebugLevel)
	logger.SetReportCaller(true)
	logger.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})
	LogrusObj = logger
}

func setOutputFile() (*os.File, error) {
	now := time.Now()
	logFilePath := ""
	if dir, err := os.Getwd(); err == nil {
		logFilePath = dir + "/logs/"
	}
	_, err := os.Stat(logFilePath)
	if os.IsNotExist(err) {
		if err := os.Mkdir(logFilePath, 0777); err != nil {
			return nil, err
		}
	}
	logFileName := now.Format("2006-01-02") + ".log"
	fileName := path.Join(logFilePath, logFileName)
	_, err = os.Stat(fileName)
	if os.IsNotExist(err) {
		if _, err := os.Create(fileName); err != nil {
			return nil, err
		}
	}
	return os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
}

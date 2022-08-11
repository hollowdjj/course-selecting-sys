package logger

import (
	"github.com/hollowdjj/course-selecting-sys/conf"
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

var logger *logrus.Logger

//Get an instance of logger
func GetInstance() *logrus.Logger {
	return logger
}

//Init logger
func InitLogger(path string) {
	logger = logrus.New()
	//show functions which print logs
	logger.SetReportCaller(true)

	//setup log format
	logger.SetFormatter(&logrus.TextFormatter{
		DisableColors:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})

	//setup rotate file
	logConf := conf.GetLogger()
	logger.SetOutput(&lumberjack.Logger{
		Filename:   path,
		MaxSize:    logConf.MaxSize,
		MaxBackups: logConf.MaxBackups,
		LocalTime:  true,
		MaxAge:     logConf.MaxAge,
	})
}

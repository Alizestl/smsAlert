// logger/logger.go
package logger

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

var Log *logrus.Logger

func InitLogger() {
	Log = logrus.New()
	logFile := &lumberjack.Logger{
		Filename:   "/var/log/webhook/alerts.log", // 需要给予sudo权限
		MaxSize:    10,                            // megabytes
		MaxBackups: 7,
		MaxAge:     7, // days
		Compress:   true,
	}
	Log.SetOutput(logFile)

	if _, err := logFile.Write([]byte("")); err != nil {
		panic(fmt.Sprintf("Failed to initialize logger: %v", err))
	}

	Log.Info("Logger initialized successfully")
}

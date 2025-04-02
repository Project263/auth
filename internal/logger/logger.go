package logger

import (
	"fmt"
	"runtime"

	"github.com/sirupsen/logrus"
)

func InitLogger(logLevel, mode string) {
	_, err := logrus.ParseLevel(logLevel)
	if err != nil {
		logrus.Fatalf("Ошибка при парсинге уровня логирования: %v", err)
	}
	// logrus.SetLevel(level)
	logrus.SetReportCaller(true)
	if mode == "prod" {
		logrus.SetFormatter(&logrus.JSONFormatter{
			CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) {
				return fmt.Sprintf("%s:%d ", frame.Function, frame.Line), ""
			},
		})
	} else {
		logrus.SetFormatter(&logrus.TextFormatter{
			FullTimestamp: true,
			CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) {
				return fmt.Sprintf("%s:%d ", frame.Function, frame.Line), ""
			},
		})
	}
}

package logger

import (
	"fmt"
	"runtime"

	"github.com/sirupsen/logrus"
)

func InitLogger(logLevel string) {
	level, err := logrus.ParseLevel(logLevel)
	if err != nil {
		logrus.Fatalf("Ошибка при парсинге уровня логирования: %v", err)
	}
	logrus.SetLevel(level)
	logrus.SetReportCaller(true)
	logrus.SetFormatter(&logrus.JSONFormatter{
		CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) {
			return fmt.Sprintf("%s:%d", frame.Function, frame.Line), ""
		},
	})

}

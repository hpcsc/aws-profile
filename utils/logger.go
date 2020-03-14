package utils

import (
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

type LogrusLogger struct {
	logger *logrus.Logger
}

func NewLogrusLogger() Logger {
	lLogger := &logrus.Logger{
		Out: &lumberjack.Logger{
			Filename: ExpandHomeDirectory("~/.aws-profile/log"),
			MaxSize:  10,
		},
		Formatter: &logrus.TextFormatter{
			FullTimestamp:          true,
			DisableLevelTruncation: true,
		},
		Level: logrus.InfoLevel,
	}

	return &LogrusLogger{
		logger: lLogger,
	}
}

func (l *LogrusLogger) Debugf(format string, args ...interface{}) {
	l.logger.Debugf(format, args...)
}

func (l *LogrusLogger) Infof(format string, args ...interface{}) {
	l.logger.Infof(format, args...)
}

func (l *LogrusLogger) Warnf(format string, args ...interface{}) {
	l.logger.Warnf(format, args...)
}

func (l *LogrusLogger) Errorf(format string, args ...interface{}) {
	l.logger.Errorf(format, args...)
}

func (l *LogrusLogger) Fatalf(format string, args ...interface{}) {
	l.logger.Fatalf(format, args...)
}

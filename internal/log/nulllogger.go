package log

type NullLogger struct {
}

func (l *NullLogger) Debugf(format string, args ...interface{}) {
}

func (l *NullLogger) Infof(format string, args ...interface{}) {
}

func (l *NullLogger) Warnf(format string, args ...interface{}) {
}

func (l *NullLogger) Errorf(format string, args ...interface{}) {
}

func (l *NullLogger) Fatalf(format string, args ...interface{}) {
}

package logger

type Logger interface {
	Infow(msg string, keysAndValues ...interface{})
	Infof(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Debugf(format string, args ...interface{})
	Debug(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
	Info(args ...interface{})
	Fatal(args ...interface{})
	Sync() error
}

package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// returnLogLevel reads the system environment "LOG_LEVEL" and pass it to the logger constructors
func returnLogLevel() zap.AtomicLevel {
	switch logLevel := os.Getenv("LOG_LEVEL"); logLevel {
	case "debug", "DEBUG":
		return zap.NewAtomicLevelAt(zapcore.DebugLevel)
	case "info", "INFO":
		return zap.NewAtomicLevelAt(zapcore.InfoLevel)
	case "warn", "WARN":
		return zap.NewAtomicLevelAt(zapcore.WarnLevel)
	case "error", "ERROR":
		return zap.NewAtomicLevelAt(zapcore.ErrorLevel)
	case "dpanic", "DPANIC":
		return zap.NewAtomicLevelAt(zapcore.DPanicLevel)
	case "panic", "PANIC":
		return zap.NewAtomicLevelAt(zapcore.PanicLevel)
	case "fatal", "FATAL":
		return zap.NewAtomicLevelAt(zapcore.FatalLevel)
	default:
		return zap.NewAtomicLevelAt(zapcore.InfoLevel)
	}
}

// New creates a Logger object with following fields: "level", "timestamp", "caller", "msg".
// If the logger prints a Fatal log the "stacktrace" field will be added.
func New(fields ...interface{}) Logger {
	var err error
	cfg := zap.Config{
		Level:       returnLogLevel(),
		Development: false,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding:         "json",
		EncoderConfig:    zap.NewProductionEncoderConfig(),
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
	}

	log, err := cfg.Build()
	if err != nil {
		exampleLog := zap.NewExample().Sugar()
		exampleLog.Error("could not initialize logger")
		return exampleLog
	}
	sugar := log.Sugar().With(fields...)
	return sugar
}

// NewZap is used for when you need a zap.Logger object. eg.: grpc_zap.UnaryServerInterceptor(zap.Logger)
func NewZap(fields map[string]interface{}) *zap.Logger {
	cfg := zap.Config{
		Level:             returnLogLevel(),
		Development:       false,
		DisableCaller:     false,
		DisableStacktrace: false,
		Sampling:          nil,
		Encoding:          "json",
		EncoderConfig:     zapcore.EncoderConfig{},
		OutputPaths:       []string{"stdout"},
		ErrorOutputPaths:  []string{"stdout"},
		InitialFields:     fields,
	}

	log, err := cfg.Build()
	if err != nil {
		exampleLog := zap.NewExample()
		exampleLog.Error("could not initialize logger")
		return exampleLog
	}
	return log
}

// WithFields adds additional fields to the logger
// To remove the fields, instantiate a new logger object with `New()`
func WithFields(log Logger, fields ...interface{}) Logger {
	if l, ok := log.(*zap.SugaredLogger); ok {
		return l.With(fields...)
	}

	log.Error("incompatible logger type")
	return log
}

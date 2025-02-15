package logger

import (
	"fmt"
	"time"

	"github.com/test/library-app/internal/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// log is a global variable which holds the logger instance created once during the app starts
var log *zap.Logger

func Log() *zap.Logger {
	if log != nil {
		return log
	}
	tmpLog, _ := zap.NewDevelopment()
	return tmpLog
}

// InitLogger initializes the logger
func InitLogger() error {
	cfg := zap.NewDevelopmentConfig()
	cfg.EncoderConfig.StacktraceKey = "stack"
	cfg.EncoderConfig.EncodeTime = func(t time.Time, pae zapcore.PrimitiveArrayEncoder) {
		pae.AppendString(t.Format(config.LogConfig.Format))
	}
	cfg.Level = zap.NewAtomicLevel()
	cfg.Encoding = config.LogConfig.Encoding

	logger, err := cfg.Build()
	if err != nil {
		log.With(zap.Error(err)).Warn("Settings not applied to the logger")
		return err
	}
	if log != nil {
		log.Sync()
	}
	log = logger
	return nil
}

func Debugf(format string, vals ...interface{}) {
	Log().Debug(fmt.Sprintf(format, vals...))
}

func Infof(format string, vals ...interface{}) {
	Log().Info(fmt.Sprintf(format, vals...))
}

func Errorf(format string, vals ...interface{}) {
	Log().Error(fmt.Sprintf(format, vals...))
}

func Warnf(format string, vals ...interface{}) {
	Log().Warn(fmt.Sprintf(format, vals...))
}

func Panicf(format string, vals ...interface{}) {
	Log().Panic(fmt.Sprintf(format, vals...))
}

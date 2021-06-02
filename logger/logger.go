package logger

import (
	"os"

	"github.com/hashicorp/go-hclog"
)

type Logger struct {
	logger hclog.Logger
}

type NewLoggerOpts struct {
	Name  string
	Level hclog.Level
}

func NewLogger(opts ...*NewLoggerOpts) *Logger {
	opt := &NewLoggerOpts{
		Level: hclog.Info,
	}

	if len(opts) > 0 {
		opt = opts[0]
	}

	hclogger := hclog.New(&hclog.LoggerOptions{
		Name:  opt.Name,
		Level: opt.Level,
	})

	return &Logger{
		logger: hclogger,
	}
}

func (l *Logger) Info(message string, args ...interface{}) {
	l.logger.Info(message, args...)
}

func (l *Logger) Warn(message string, args ...interface{}) {
	l.logger.Warn(message, args...)
}

func (l *Logger) Error(message string, args ...interface{}) {
	l.logger.Error(message, args...)
}

func (l *Logger) Debug(message string, args ...interface{}) {
	l.logger.Debug(message, args...)
}

func (l *Logger) Fatal(message string, args ...interface{}) {
	l.logger.Log(hclog.LevelFromString("FATAL"), message, args...)

	l.exit()
}

func (l *Logger) exit() {
	os.Exit(1)
}

func LevelInfo() hclog.Level {
	return hclog.Info
}

func LevelWarn() hclog.Level {
	return hclog.Warn
}

func LevelError() hclog.Level {
	return hclog.Error
}

func LevelDebug() hclog.Level {
	return hclog.Debug
}

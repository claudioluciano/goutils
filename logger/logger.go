package logger

import (
	"github.com/sirupsen/logrus"
)

type Logger struct {
	internalLogger *logrus.Logger
}

type NewLoggerOpts struct {
	LogLevel logrus.Level
}

func NewLogger(opts *NewLoggerOpts) *Logger {
	if opts == nil {
		opts = &NewLoggerOpts{
			LogLevel: logrus.InfoLevel,
		}
	}

	baseLogger := logrus.New()

	logger := &Logger{internalLogger: baseLogger}

	logger.internalLogger.Level = opts.LogLevel
	logger.internalLogger.Formatter = &logrus.JSONFormatter{}

	return logger
}

func (l *Logger) Info(message string) {
	l.internalLogger.Info(message)
}

func (l *Logger) InfoWithFields(message error, fields map[string]interface{}) {
	l.internalLogger.WithFields(fields).Info(message)
}

func (l *Logger) Warn(message string) {
	l.internalLogger.Warn(message)
}

func (l *Logger) WarnWithFields(message error, fields map[string]interface{}) {
	l.internalLogger.WithFields(fields).Warn(message)
}

func (l *Logger) Error(message string) {
	l.internalLogger.Error(message)
}

func (l *Logger) ErrorWithError(message string, err error) {
	l.internalLogger.WithError(err).Error(message)
}

func (l *Logger) Debug(message string) {
	l.internalLogger.Debug(message)
}

func (l *Logger) DebugWithFields(message string, fields map[string]interface{}) {
	l.internalLogger.WithFields(fields).Debug(message)
}

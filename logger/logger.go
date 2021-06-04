package logger

import (
	"os"

	"github.com/hashicorp/go-hclog"
)

type Client struct {
	logger hclog.Logger
}

type NewLoggerOpts struct {
	Name  string
	Level hclog.Level
}

func NewLogger(opts ...*NewLoggerOpts) *Client {
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

	return &Client{
		logger: hclogger,
	}
}

func (l *Client) Info(message string, args ...interface{}) {
	l.logger.Info(message, args...)
}

func (l *Client) Warn(message string, args ...interface{}) {
	l.logger.Warn(message, args...)
}

func (l *Client) Error(message string, args ...interface{}) {
	l.logger.Error(message, args...)
}

func (l *Client) Debug(message string, args ...interface{}) {
	l.logger.Debug(message, args...)
}

func (l *Client) Fatal(message string, args ...interface{}) {
	l.logger.Log(hclog.LevelFromString("FATAL"), message, args...)

	l.exit()
}

func (l *Client) exit() {
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

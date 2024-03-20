package lumberjack

import (
	"net/url"

	lumberjack2 "gopkg.in/natefinch/lumberjack.v2"
)

type Logger struct {
	*lumberjack2.Logger
}

func NewLogger(config *Config) *Logger {
	temp := newDefaultConfig()
	if config == nil {
		return temp
	}

	temp.Filename = config.Filename
	temp.Compress = config.Compress
	temp.LocalTime = config.Localtime
	if config.Maxsize > 0 {
		temp.MaxSize = int(config.Maxsize)
	}
	if config.Maxage > 0 {
		temp.MaxAge = int(config.Maxage)
	}
	if config.Maxbackups > 0 {
		temp.MaxBackups = int(config.Maxbackups)
	}
	return temp
}

func (l *Logger) Sync() error {
	return nil
}

func NewLoggerWithURL(config *Config, u *url.URL) *Logger {
	l := NewLogger(config)
	if u == nil {
		return l
	}

	if l.Filename == "" {
		if u.Opaque != "" {
			l.Filename = u.Opaque
		} else {
			l.Filename = u.Path
		}
	}
	return l
}

func newDefaultConfig() *Logger {
	return &Logger{Logger: &lumberjack2.Logger{
		MaxSize:    1024,
		MaxAge:     7,
		MaxBackups: 3,
	}}
}

package zap

import (
	"fmt"

	"github.com/go-kratos/kratos/v2/log"
	"go.uber.org/zap"
)

var _ log.Logger = (*Logger)(nil)

type Logger struct {
	log *zap.Logger
}

func NewLogger(opts ...Option) (*Logger, error) {
	_options := options{
		zapConfig: zap.NewProductionConfig(),
		zapOptions: []zap.Option{
			zap.AddCallerSkip(3),
		},
	}

	for _, o := range opts {
		o(&_options)
	}

	zlog, err := _options.zapConfig.Build(_options.zapOptions...)
	if err != nil {
		return nil, err
	}

	l := &Logger{zlog}

	return l, nil
}

func (l *Logger) Log(level log.Level, keyvals ...interface{}) error {
	if len(keyvals) == 0 || len(keyvals)%2 != 0 {
		l.log.Warn(fmt.Sprint("Keyvalues must appear in pairs: ", keyvals))
		return nil
	}

	var data []zap.Field
	for i := 0; i < len(keyvals); i += 2 {
		data = append(data, zap.Any(fmt.Sprint(keyvals[i]), keyvals[i+1]))
	}

	switch level {
	case log.LevelDebug:
		l.log.Debug("", data...)
	case log.LevelInfo:
		l.log.Info("", data...)
	case log.LevelWarn:
		l.log.Warn("", data...)
	case log.LevelError:
		l.log.Error("", data...)
	case log.LevelFatal:
		l.log.Fatal("", data...)
	}
	return nil
}

func (l *Logger) Sync() error {
	return l.log.Sync()
}

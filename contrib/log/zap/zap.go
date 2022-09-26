package zap

import (
	"fmt"

	"github.com/go-kratos/kratos/v2/log"
	"go.uber.org/zap"
)

var _ log.Logger = (*Logger)(nil)

// Option is zap logger option.
type Option func(*options)

type options struct {
	isSugar bool
}

// WithLogSugar with sugared logger.
func WithLogSugar() Option {
	return func(o *options) {
		o.isSugar = true
	}
}

type Logger struct {
	log     *zap.Logger
	isSugar bool
}

func NewLogger(zlog *zap.Logger, opts ...Option) *Logger {
	opt := &options{isSugar: false}
	for _, o := range opts {
		o(opt)
	}

	return &Logger{
		log:     zlog,
		isSugar: opt.isSugar,
	}
}

func (l *Logger) Log(level log.Level, keyvals ...interface{}) error {
	if len(keyvals) == 0 || len(keyvals)%2 != 0 {
		if l.isSugar {
			if len(keyvals) == 0 {
				return nil
			}
			if len(keyvals)%2 != 0 {
				keyvals = append(keyvals, "")
			}

		} else {
			l.log.Warn(fmt.Sprint("Keyvalues must appear in pairs: ", keyvals))
			return nil
		}
	}

	var (
		data []zap.Field
		msg  string
	)
	for i := 0; i < len(keyvals); i += 2 {
		if l.isSugar {
			key, ok := keyvals[i].(string)
			if !ok {
				continue
			}
			if key == log.DefaultMessageKey {
				msg, _ = keyvals[i+1].(string)
			}

		} else {
			data = append(data, zap.Any(fmt.Sprint(keyvals[i]), keyvals[i+1]))
		}
	}

	if l.isSugar {
		switch level {
		case log.LevelDebug:
			l.log.Sugar().Debug(msg)
		case log.LevelInfo:
			l.log.Sugar().Info(msg)
		case log.LevelWarn:
			l.log.Sugar().Warn(msg)
		case log.LevelError:
			l.log.Sugar().Error(msg)
		case log.LevelFatal:
			l.log.Sugar().Fatal(msg)
		}

	} else {
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
	}

	return nil
}

func (l *Logger) Sync() error {
	return l.log.Sync()
}

func (l *Logger) Close() error {
	return l.Sync()
}

package zap

import (
	"fmt"
	"net/url"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var _ log.Logger = (*Logger)(nil)

type Logger struct {
	log    *zap.Logger
	msgKey string
}

type Option func(*Logger)

// WithMessageKey with message key.
func WithMessageKey(key string) Option {
	return func(l *Logger) {
		l.msgKey = key
	}
}

func NewLogger(zlog *zap.Logger) *Logger {
	return &Logger{
		log:    zlog,
		msgKey: log.DefaultMessageKey,
	}
}

func (l *Logger) Log(level log.Level, keyvals ...interface{}) error {
	var (
		msg    = ""
		keylen = len(keyvals)
	)
	if keylen == 0 || keylen%2 != 0 {
		l.log.Warn(fmt.Sprint("Keyvalues must appear in pairs: ", keyvals))
		return nil
	}

	data := make([]zap.Field, 0, (keylen/2)+1)
	for i := 0; i < keylen; i += 2 {
		if keyvals[i].(string) == l.msgKey {
			msg, _ = keyvals[i+1].(string)
			continue
		}
		data = append(data, zap.Any(fmt.Sprint(keyvals[i]), keyvals[i+1]))
	}

	switch level {
	case log.LevelDebug:
		l.log.Debug(msg, data...)
	case log.LevelInfo:
		l.log.Info(msg, data...)
	case log.LevelWarn:
		l.log.Warn(msg, data...)
	case log.LevelError:
		l.log.Error(msg, data...)
	case log.LevelFatal:
		l.log.Fatal(msg, data...)
	}
	return nil
}

func (l *Logger) Sync() error {
	return l.log.Sync()
}

func (l *Logger) Close() error {
	return l.Sync()
}

////////////////////////////////////////////////////////////////

type Sink = zap.Sink

func RegisterSink(scheme string, factory func(*url.URL) (Sink, error)) error {
	return zap.RegisterSink(scheme, factory)
}

func NewZapLogger(config *Config, opts ...zap.Option) (*Logger, error) {
	zapConfig := newDefaultConfig()

	if config != nil {
		zapConfig.Development = config.Development
		zapConfig.DisableCaller = config.DisableCaller
		zapConfig.DisableStacktrace = config.DisableStacktrace
		if config.Encoding != "" {
			zapConfig.Encoding = config.Encoding
		}
		if len(config.OutputPaths) != 0 {
			zapConfig.OutputPaths = config.OutputPaths
		}
		if len(config.ErrorOutputPaths) != 0 {
			zapConfig.ErrorOutputPaths = config.ErrorOutputPaths
		}
		for k, v := range config.InitialFields {
			if zapConfig.InitialFields == nil {
				zapConfig.InitialFields = make(map[string]any)
			}
			zapConfig.InitialFields[k] = v
		}
		if config.Level != "" {
			var err error
			zapConfig.Level, err = zap.ParseAtomicLevel(config.Level)
			if err != nil {
				return nil, errors.Wrap(err, "zap config level invalid")
			}
		}
		if config.EncoderConfig != nil {
			zapConfig.EncoderConfig.SkipLineEnding = config.EncoderConfig.SkipLineEnding
			zapConfig.EncoderConfig.MessageKey = genKey(zapConfig.EncoderConfig.MessageKey, config.EncoderConfig.MessageKey)
			zapConfig.EncoderConfig.LevelKey = genKey(zapConfig.EncoderConfig.LevelKey, config.EncoderConfig.LevelKey)
			zapConfig.EncoderConfig.TimeKey = genKey(zapConfig.EncoderConfig.TimeKey, config.EncoderConfig.TimeKey)
			zapConfig.EncoderConfig.NameKey = genKey(zapConfig.EncoderConfig.NameKey, config.EncoderConfig.NameKey)
			zapConfig.EncoderConfig.CallerKey = genKey(zapConfig.EncoderConfig.CallerKey, config.EncoderConfig.CallerKey)
			zapConfig.EncoderConfig.FunctionKey = genKey(zapConfig.EncoderConfig.FunctionKey, config.EncoderConfig.FunctionKey)
			zapConfig.EncoderConfig.StacktraceKey = genKey(zapConfig.EncoderConfig.StacktraceKey, config.EncoderConfig.StacktraceKey)
			zapConfig.EncoderConfig.LevelKey = genKey(zapConfig.EncoderConfig.LevelKey, config.EncoderConfig.LevelKey)
			zapConfig.EncoderConfig.LevelKey = genKey(zapConfig.EncoderConfig.LevelKey, config.EncoderConfig.LevelKey)
			zapConfig.EncoderConfig.LevelKey = genKey(zapConfig.EncoderConfig.LevelKey, config.EncoderConfig.LevelKey)
			zapConfig.EncoderConfig.LevelKey = genKey(zapConfig.EncoderConfig.LevelKey, config.EncoderConfig.LevelKey)
			zapConfig.EncoderConfig.LevelKey = genKey(zapConfig.EncoderConfig.LevelKey, config.EncoderConfig.LevelKey)
			if config.EncoderConfig.LineEnding != "" {
				zapConfig.EncoderConfig.LineEnding = config.EncoderConfig.LineEnding
			}
			if config.EncoderConfig.ConsoleSeparator != "" {
				zapConfig.EncoderConfig.ConsoleSeparator = config.EncoderConfig.ConsoleSeparator
			}
		}
	}

	l, err := zapConfig.Build(opts...)
	if err != nil {
		return nil, errors.Wrap(err, "zap config build failed")
	}
	kl := NewLogger(l)
	if zapConfig.EncoderConfig.MessageKey != "" {
		kl.msgKey = zapConfig.EncoderConfig.MessageKey
	}
	return kl, nil
}

func newDefaultConfig() zap.Config {
	zapConfig := zap.NewProductionConfig()
	zapConfig.Sampling = nil                                          // 暂不支持采样
	zapConfig.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder   // 默认使用ISO8601时间编码器
	zapConfig.EncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder // 使用大写字母记录日志级别
	return zapConfig
}

func genKey(old, new string) string {
	if new == "" {
		return old
	}
	if new == "-" {
		return ""
	}
	return new
}

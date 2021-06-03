package log

import (
	"fmt"
	"os"

	"github.com/go-kratos/kratos/v2/log"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var _ log.Logger = (*ZapLogger)(nil)

type ZapLogger struct {
	log  *zap.Logger
	Sync func() error
}

// NewZapLogger return ZapLogger
func NewZapLogger(encoder zapcore.EncoderConfig, level zap.AtomicLevel, opts ...zap.Option) *ZapLogger {
	writeSyncer := getLogWriter()
	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoder),
		zapcore.NewMultiWriteSyncer(
			zapcore.AddSync(os.Stdout),
			writeSyncer,
		), level)
	zapLogger := zap.New(core, opts...)
	return &ZapLogger{log: zapLogger, Sync: zapLogger.Sync}
}

// Log Implementation of logger interface
func (l *ZapLogger) Log(level log.Level, keyvals ...interface{}) error {
	fmt.Println(len(keyvals))
	if len(keyvals) == 0 {
		return nil
	}
	// Use zap.Sugar() when keyvals is odd
	if len(keyvals)%2 != 0 {
		switch level {
		case 0:
			l.log.Sugar().Debug(keyvals)
		case 1:
			l.log.Sugar().Info(keyvals)
		case 2:
			l.log.Sugar().Warn(keyvals)
		case 3:
			l.log.Sugar().Error(keyvals)
		}
		return nil
	}

	// Zap.Field is used when keyvals pairs appear
	var data []zap.Field
	for i := 0; i < len(keyvals); i += 2 {
		data = append(data, zap.Any(fmt.Sprint(keyvals[i]), fmt.Sprint(keyvals[i+1])))
	}
	switch level {
	case 0:
		l.log.Debug("", data...)
	case 1:
		l.log.Info("", data...)
	case 2:
		l.log.Warn("", data...)
	case 3:
		l.log.Error("", data...)
	}
	return nil
}

// getLogWriter writing logs to rolling files
func getLogWriter() zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   "./test.log",
		MaxSize:    10,
		MaxBackups: 5,
		MaxAge:     30,
		Compress:   false,
	}
	return zapcore.AddSync(lumberJackLogger)
}
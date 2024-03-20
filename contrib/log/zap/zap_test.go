package zap

import (
	"net/url"
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/stretchr/testify/assert"
)

type testWriteSyncer struct {
	output []string
}

func (x *testWriteSyncer) Write(p []byte) (n int, err error) {
	x.output = append(x.output, string(p))
	return len(p), nil
}

func (x *testWriteSyncer) Sync() error {
	return nil
}

func (x *testWriteSyncer) Close() error {
	return nil
}

func TestLogger(t *testing.T) {
	syncer := &testWriteSyncer{}
	encoderCfg := zapcore.EncoderConfig{
		MessageKey:     "msg",
		LevelKey:       "level",
		NameKey:        "logger",
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
	}
	core := zapcore.NewCore(zapcore.NewJSONEncoder(encoderCfg), syncer, zap.DebugLevel)
	zlogger := zap.New(core).WithOptions()
	logger := NewLogger(zlogger)

	defer func() { _ = logger.Close() }()

	zlog := log.NewHelper(logger)

	zlog.Debugw("log", "debug")
	zlog.Infow("log", "info")
	zlog.Warnw("log", "warn")
	zlog.Errorw("log", "error")
	zlog.Errorw("log", "error", "except warn")

	except := []string{
		"{\"level\":\"debug\",\"msg\":\"\",\"log\":\"debug\"}\n",
		"{\"level\":\"info\",\"msg\":\"\",\"log\":\"info\"}\n",
		"{\"level\":\"warn\",\"msg\":\"\",\"log\":\"warn\"}\n",
		"{\"level\":\"error\",\"msg\":\"\",\"log\":\"error\"}\n",
		"{\"level\":\"warn\",\"msg\":\"Keyvalues must appear in pairs: [log error except warn]\"}\n",
	}
	for i, s := range except {
		if s != syncer.output[i] {
			t.Logf("except=%s, got=%s", s, syncer.output[i])
			t.Fail()
		}
	}
}

func TestZapLogger(t *testing.T) {
	cfg := &Config{
		Level:         "debug",
		OutputPaths:   []string{"test:"}, // 指定test sink输出日志
		InitialFields: map[string]string{"m": "a"},
		EncoderConfig: &EncoderConfig{
			LevelKey:      "level",
			NameKey:       "logger",
			TimeKey:       "-", // 关闭时间输出
			StacktraceKey: "-", // 关闭stacktrace输出
		},
	}
	// 注册test sink
	syncer := &testWriteSyncer{}
	err := RegisterSink("test", func(_ *url.URL) (Sink, error) {
		return syncer, nil
	})
	assert.NoError(t, err)
	logger, err := NewZapLogger(cfg, zap.WithCaller(false))
	assert.NoError(t, err)
	defer func() { _ = logger.Close() }()

	zlog := log.NewHelper(logger)

	zlog.Debugw("msg", "debug")
	zlog.Infow("msg", "info")
	zlog.Warnw("msg", "warn")
	zlog.Errorw("msg", "error")
	zlog.Errorw("msg", "error", "except warn")

	except := []string{
		"{\"level\":\"DEBUG\",\"msg\":\"debug\",\"m\":\"a\"}\n",
		"{\"level\":\"INFO\",\"msg\":\"info\",\"m\":\"a\"}\n",
		"{\"level\":\"WARN\",\"msg\":\"warn\",\"m\":\"a\"}\n",
		"{\"level\":\"ERROR\",\"msg\":\"error\",\"m\":\"a\"}\n",
		"{\"level\":\"WARN\",\"msg\":\"Keyvalues must appear in pairs: [msg error except warn]\",\"m\":\"a\"}\n",
	}
	for i, s := range except {
		if s != syncer.output[i] {
			t.Logf("except=%s, got=%s", s, syncer.output[i])
			t.Fail()
		}
	}
}

func TestZapLoggerMsgKey(t *testing.T) {
	cfg := &Config{
		Level:       "debug",
		OutputPaths: []string{"test2:"}, // 指定test sink输出日志
		EncoderConfig: &EncoderConfig{
			MessageKey:    "log",
			LevelKey:      "level",
			NameKey:       "logger",
			TimeKey:       "-", // 关闭时间输出
			StacktraceKey: "-", // 关闭stacktrace输出
		},
	}
	// 注册test sink
	syncer := &testWriteSyncer{}
	err := RegisterSink("test2", func(_ *url.URL) (Sink, error) {
		return syncer, nil
	})
	assert.NoError(t, err)
	logger, err := NewZapLogger(cfg, zap.WithCaller(false))
	assert.NoError(t, err)
	defer func() { _ = logger.Close() }()

	zlog := log.NewHelper(logger)

	zlog.Debugw("log", "debug")
	zlog.Infow("log", "info")
	zlog.Warnw("log", "warn")
	zlog.Errorw("log", "error")
	zlog.Errorw("log", "error", "m", "a")

	except := []string{
		"{\"level\":\"DEBUG\",\"log\":\"debug\"}\n",
		"{\"level\":\"INFO\",\"log\":\"info\"}\n",
		"{\"level\":\"WARN\",\"log\":\"warn\"}\n",
		"{\"level\":\"ERROR\",\"log\":\"error\"}\n",
		"{\"level\":\"ERROR\",\"log\":\"error\",\"m\":\"a\"}\n",
	}
	for i, s := range except {
		if s != syncer.output[i] {
			t.Logf("except=%s, got=%s", s, syncer.output[i])
			t.Fail()
		}
	}
}

package zap

import (
	"testing"

	"github.com/go-kratos/kratos/v2/log"
)

func TestLogger(t *testing.T) {
	logger, err := NewLogger(
		WithFormat("json"),
	)
	if err != nil {
		t.Error(err)
	}
	defer func() { _ = logger.Sync() }()

	zlog := log.NewHelper(logger)

	zlog.Debugw("log", "debug")
	zlog.Infow("log", "info")
	zlog.Warnw("log", "warn")
	zlog.Errorw("log", "error")
}

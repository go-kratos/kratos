package zap

import (
	"testing"

	"go.uber.org/zap"

	"github.com/SeeMusic/kratos/v2/log"
)

func TestLogger(t *testing.T) {
	logger := NewLogger(zap.NewExample())

	defer func() { _ = logger.Sync() }()

	zlog := log.NewHelper(logger)

	zlog.Debugw("log", "debug")
	zlog.Infow("log", "info")
	zlog.Warnw("log", "warn")
	zlog.Errorw("log", "error")
}

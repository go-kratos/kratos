package aliyun

import (
	"testing"

	"github.com/go-kratos/kratos/v2/log"
)

func TestAliyunLogger(t *testing.T) {
	tmpLogger := NewAliyunLog(
		WithEndpoint(""),
		WithAccessKey(""),
		WithAccessSecret(""),
		WithLogstore(""),
		WithProject(""),
	)

	if aliLogger, ok := tmpLogger.(AliyunLog); ok {
		aliLogger.GetProducer().Start()
		defer aliLogger.GetProducer().Close(6000)

		_ = aliLogger.Log(log.LevelDebug, "key", "kratos")
	}
}

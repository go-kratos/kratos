package tencent

import (
	"testing"

	"github.com/go-kratos/kratos/v2/log"
)

func TestWithEndpoint(t *testing.T) {
	opts := new(options)
	endpoint := "eee"
	funcEndpoint := WithEndpoint(endpoint)
	funcEndpoint(opts)
	if opts.endpoint != "eee" {
		t.Errorf("WithEndpoint() = %s, want %s", opts.endpoint, endpoint)
	}
}

func TestWithTopicId(t *testing.T) {
	opts := new(options)
	topicID := "ee"
	funcTopicID := WithTopicID(topicID)
	funcTopicID(opts)
	if opts.topicID != "ee" {
		t.Errorf("WithTopicId() = %s, want %s", opts.endpoint, topicID)
	}
}

func TestWithAccessKey(t *testing.T) {
	opts := new(options)
	accessKey := "ee"
	funcAccessKey := WithAccessKey(accessKey)
	funcAccessKey(opts)
	if opts.accessKey != "ee" {
		t.Errorf("WithAccessKey() = %s, want %s", opts.endpoint, accessKey)
	}
}

func TestWithAccessSecret(t *testing.T) {
	opts := new(options)
	accessSecret := "ee"
	funcAccessSecret := WithAccessSecret(accessSecret)
	funcAccessSecret(opts)
	if opts.accessSecret != "ee" {
		t.Errorf("WithAccessSecret() = %s, want %s", opts.accessSecret, accessSecret)
	}
}

func TestTestLogger(t *testing.T) {
	topicID := "aa"
	logger, err := NewLogger(
		WithTopicID(topicID),
		WithEndpoint("ap-shanghai.cls.tencentcs.com"),
		WithAccessKey("a"),
		WithAccessSecret("b"),
	)
	if err != nil {
		t.Error(err)
		return
	}
	defer logger.Close()
	logger.GetProducer()
	flog := log.NewHelper(logger)
	flog.Debug("log", "test")
	flog.Info("log", "test")
	flog.Warn("log", "test")
	flog.Error("log", "test")
}

func TestLog(t *testing.T) {
	topicID := "foo"
	logger, err := NewLogger(
		WithTopicID(topicID),
		WithEndpoint("ap-shanghai.cls.tencentcs.com"),
		WithAccessKey("a"),
		WithAccessSecret("b"),
	)
	if err != nil {
		t.Error(err)
		return
	}
	defer logger.Close()
	err = logger.Log(log.LevelDebug, 0, int8(1), int16(2), int32(3))
	if err != nil {
		t.Errorf("Log() returns error:%v", err)
	}
	err = logger.Log(log.LevelDebug, uint(0), uint8(1), uint16(2), uint32(3))
	if err != nil {
		t.Errorf("Log() returns error:%v", err)
	}
	err = logger.Log(log.LevelDebug, uint(0), uint8(1), uint16(2), uint32(3))
	if err != nil {
		t.Errorf("Log() returns error:%v", err)
	}
	err = logger.Log(log.LevelDebug, int64(0), uint64(1), float32(2), float64(3))
	if err != nil {
		t.Errorf("Log() returns error:%v", err)
	}
	err = logger.Log(log.LevelDebug, []byte{0, 1, 2, 3}, "foo")
	if err != nil {
		t.Errorf("Log() returns error:%v", err)
	}
}

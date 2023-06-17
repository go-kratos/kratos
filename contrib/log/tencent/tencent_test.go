package tencent

import (
	"math"
	"reflect"
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

func TestNewString(t *testing.T) {
	ptr := newString("")
	if kind := reflect.TypeOf(ptr).Kind(); kind != reflect.Ptr {
		t.Errorf("want type: %v, got type: %v", reflect.Ptr, kind)
	}
}

func TestToString(t *testing.T) {
	tests := []struct {
		name string
		in   interface{}
		out  string
	}{
		{"float64", 6.66, "6.66"},
		{"max float64", math.MaxFloat64, "179769313486231570000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"}, //nolint:lll
		{"float32", float32(6.66), "6.66"},
		{"max float32", float32(math.MaxFloat32), "340282350000000000000000000000000000000"},
		{"int", math.MaxInt64, "9223372036854775807"},
		{"uint", uint(math.MaxUint64), "18446744073709551615"},
		{"int8", int8(math.MaxInt8), "127"},
		{"uint8", uint8(math.MaxUint8), "255"},
		{"int16", int16(math.MaxInt16), "32767"},
		{"uint16", uint16(math.MaxUint16), "65535"},
		{"int32", int32(math.MaxInt32), "2147483647"},
		{"uint32", uint32(math.MaxUint32), "4294967295"},
		{"int64", int64(math.MaxInt64), "9223372036854775807"},
		{"uint64", uint64(math.MaxUint64), "18446744073709551615"},
		{"string", "abc", "abc"},
		{"bool", false, "false"},
		{"[]byte", []byte("abc"), "abc"},
		{"struct", struct{ Name string }{}, `{"Name":""}`},
		{"nil", nil, ""},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			out := toString(test.in)
			if test.out != out {
				t.Fatalf("want: %s, got: %s", test.out, out)
			}
		})
	}
}

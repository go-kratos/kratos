package aliyun

import (
	"math"
	"reflect"
	"testing"

	"github.com/go-kratos/kratos/v2/log"
)

func TestWithEndpoint(t *testing.T) {
	opts := new(options)
	endpoint := "foo"
	funcEndpoint := WithEndpoint(endpoint)
	funcEndpoint(opts)
	if opts.endpoint != endpoint {
		t.Errorf("WithEndpoint() = %s, want %s", opts.endpoint, endpoint)
	}
}

func TestWithProject(t *testing.T) {
	opts := new(options)
	project := "foo"
	funcProject := WithProject(project)
	funcProject(opts)
	if opts.project != project {
		t.Errorf("WithProject() = %s, want %s", opts.project, project)
	}
}

func TestWithLogstore(t *testing.T) {
	opts := new(options)
	logstore := "foo"
	funcLogstore := WithLogstore(logstore)
	funcLogstore(opts)
	if opts.logstore != logstore {
		t.Errorf("WithLogatore() = %s, want %s", opts.logstore, logstore)
	}
}

func TestWithAccessKey(t *testing.T) {
	opts := new(options)
	accessKey := "foo"
	funcAccessKey := WithAccessKey(accessKey)
	funcAccessKey(opts)
	if opts.accessKey != accessKey {
		t.Errorf("WithAccessKey() = %s, want %s", opts.accessKey, accessKey)
	}
}

func TestWithAccessSecret(t *testing.T) {
	opts := new(options)
	accessSecret := "foo"
	funcAccessSecret := WithAccessSecret(accessSecret)
	funcAccessSecret(opts)
	if opts.accessSecret != accessSecret {
		t.Errorf("WithAccessSecret() = %s, want %s", opts.accessSecret, accessSecret)
	}
}

func TestLogger(_ *testing.T) {
	project := "foo"
	logger := NewAliyunLog(WithProject(project))
	defer logger.Close()
	logger.GetProducer()
	flog := log.NewHelper(logger)
	flog.Debug("log", "test")
	flog.Info("log", "test")
	flog.Warn("log", "test")
	flog.Error("log", "test")
}

func TestLog(t *testing.T) {
	project := "foo"
	logger := NewAliyunLog(WithProject(project))
	defer logger.Close()
	err := logger.Log(log.LevelDebug, 0, int8(1), int16(2), int32(3))
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
	err = logger.Log(log.LevelDebug, true, 0)
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

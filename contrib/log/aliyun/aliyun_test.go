package aliyun

import (
	"context"
	"log/slog"
	"math"
	"reflect"
	"testing"
	"time"

	klog "github.com/go-kratos/kratos/v3/log"
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

func TestWithLogOptions(t *testing.T) {
	opts := new(options)
	fn := WithLogOptions(klog.WithAttrs(slog.String("service.name", "helloworld")))
	fn(opts)
	if len(opts.logOptions) != 1 {
		t.Fatalf("len(logOptions) = %d, want 1", len(opts.logOptions))
	}
}

func TestLogger(t *testing.T) {
	project := "foo"
	handler, err := NewHandler(WithProject(project))
	if err != nil {
		t.Error(err)
		return
	}
	defer handler.Close()
	handler.GetProducer()
	logger := slog.New(handler)
	logger.Debug("logtest")
	logger.Info("logtest")
	logger.Warn("logtest")
	logger.Error("logtest")
}

func TestLog(t *testing.T) {
	project := "foo"
	handler, err := NewHandler(WithProject(project))
	if err != nil {
		t.Error(err)
		return
	}
	defer handler.Close()
	record := slog.NewRecord(time.Now(), slog.LevelDebug, "message", 0)
	record.AddAttrs(
		slog.Any("int8", int8(1)),
		slog.Any("int16", int16(2)),
		slog.Any("int32", int32(3)),
		slog.Any("uint8", uint8(1)),
		slog.Any("uint16", uint16(2)),
		slog.Any("uint32", uint32(3)),
		slog.Any("int64", int64(0)),
		slog.Any("uint64", uint64(1)),
		slog.Any("float32", float32(2)),
		slog.Any("float64", float64(3)),
		slog.Any("bytes", []byte{0, 1, 2, 3}),
		slog.Any("bool", true),
	)
	err = handler.Handle(context.Background(), record)
	if err != nil {
		t.Errorf("Handle() returns error:%v", err)
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
		in   any
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

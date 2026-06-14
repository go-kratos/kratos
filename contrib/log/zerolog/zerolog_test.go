package zerolog

import (
	"errors"
	"testing"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/rs/zerolog"
)

type testWriter struct {
	output string
}

func (w *testWriter) Write(p []byte) (n int, err error) {
	w.output = string(p)
	return len(p), nil
}

func TestLogger(t *testing.T) {
	out := &testWriter{}
	ll := zerolog.New(out)
	logger := NewLogger(ll)

	zlog := log.NewHelper(logger)

	testcases := []struct {
		fn      func(keyvals ...interface{})
		keyvals []interface{}
		expect  string
	}{
		{
			fn:      zlog.Debugw,
			keyvals: []interface{}{"id", int64(2234222344), "int32", int32(-7463882), "int16", int16(8384), "int8", int8(22), "int", int(-328833)},
			expect:  `{"level":"debug","id":2234222344,"int32":-7463882,"int16":8384,"int8":22,"int":-328833}` + "\n",
		},
		{
			fn:      zlog.Debugw,
			keyvals: []interface{}{"uint8", uint8(35), "uint16", uint16(2134), "uint32", uint32(32345), "uint64", uint64(88739222), "uint", uint(328833)},
			expect:  `{"level":"debug","uint8":35,"uint16":2134,"uint32":32345,"uint64":88739222,"uint":328833}` + "\n",
		},
		{
			fn:      zlog.Debugw,
			keyvals: []interface{}{"float32", float32(35.0000001), "float64", float64(213444863.3329988888999)},
			expect:  `{"level":"debug","float32":35,"float64":213444863.3329989}` + "\n",
		},
		{
			fn:      zlog.Infow,
			keyvals: []interface{}{"ts", time.Date(2023, 6, 1, 12, 13, 14, 0, time.UTC), "dur", 34 * time.Second},
			expect:  `{"level":"info","ts":"2023-06-01T12:13:14Z","dur":34000}` + "\n",
		},
		{
			fn:      zlog.Infow,
			keyvals: []interface{}{"obj", map[string]interface{}{"f1": "\x1ffoo\x1fbar\x1fbazemoji \u2764\ufe0f!", "f2": 342342333, "f3": false}},
			expect:  `{"level":"info","obj":{"f1":"\u001ffoo\u001fbar\u001fbazemoji ❤️!","f2":342342333,"f3":false}}` + "\n",
		},
		{
			fn:      zlog.Warnw,
			keyvals: []interface{}{"msg", "id not exist", "id", uint64(384383343242)},
			expect:  `{"level":"warn","msg":"id not exist","id":384383343242}` + "\n",
		},
		{
			fn:      zlog.Errorw,
			keyvals: []interface{}{"err", errors.New("io error")},
			expect:  `{"level":"error","err":"io error"}` + "\n",
		},
		{
			fn:      zlog.Infow,
			keyvals: []interface{}{"msg", "test marshal", "except warn"},
			expect:  `{"level":"warn","msg":"test marshal","except warn":"KEYVALS UNPAIRED"}` + "\n",
		},
	}

	for _, s := range testcases {
		s.fn(s.keyvals...)
		if s.expect != out.output {
			t.Fatalf("except=%s, got=%s", s.expect, out.output)
		}
	}
}

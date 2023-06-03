package zerolog

import (
	"errors"
	"testing"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/rs/zerolog"
)

type testWriter struct {
	output []string
}

func (x *testWriter) Write(p []byte) (n int, err error) {
	x.output = append(x.output, string(p))
	return len(p), nil
}

func TestLogger(t *testing.T) {
	out := &testWriter{output: make([]string, 0, 16)}
	ll := zerolog.New(out)
	logger := NewLogger(ll)

	zlog := log.NewHelper(logger)

	zlog.Debugw("id", int64(2234222344))
	zlog.Infow("obj", map[string]interface{}{"f1": "\x1ffoo\x1fbar\x1fbazemoji \u2764\ufe0f!", "f2": 342342333, "f3": false})
	zlog.Warnw("msg", "id not exist", "id", uint64(384383343242))
	zlog.Errorw("err", errors.New("io error"))
	zlog.Infow("msg", "test marshal", "except warn")

	except := []string{
		`{"level":"debug","id":2234222344}` + "\n",
		`{"level":"info","obj":{"f1":"\u001ffoo\u001fbar\u001fbazemoji ❤️!","f2":342342333,"f3":false}}` + "\n",
		`{"level":"warn","msg":"id not exist","id":384383343242}` + "\n",
		`{"level":"error","err":"io error"}` + "\n",
		`{"level":"warn","msg":"test marshal","except warn":"KEYVALS UNPAIRED"}` + "\n",
	}
	for i, s := range except {
		if s != out.output[i] {
			t.Fatalf("except=%s, got=%s", s, out.output[i])
		}
	}
}

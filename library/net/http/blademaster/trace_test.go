package blademaster

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"go-common/library/net/trace"
)

func TestTrace(t *testing.T) {
	wait := make(chan bool, 1)
	eng := New()
	eng.Use(Trace())
	eng.GET("/test-trace", func(c *Context) {
		if _, ok := trace.FromContext(c.Context); !ok {
			t.Errorf("expect get trace from context")
		}
		c.Writer.Write([]byte("pong"))
		wait <- true
	})
	go eng.Run("127.0.0.1:28080")
	time.Sleep(time.Second)
	http.Get("http://127.0.0.1:28080/test-trace")
	<-wait
}

func TestTraceClient(t *testing.T) {
	wait := make(chan bool, 1)
	trace.Init(nil)
	root := trace.New("test-title")
	ctx := trace.NewContext(context.Background(), root)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		keys := []string{trace.KeyTraceID, trace.KeyTraceSpanID, trace.KeyTraceParentID, trace.KeyTraceID}
		for _, key := range keys {
			if r.Header.Get(key) == "" {
				t.Errorf("empty key: %s", key)
			}
		}
		wait <- true
	}))
	defer srv.Close()
	req, err := http.NewRequest(http.MethodGet, srv.URL, nil)
	if err != nil {
		t.Fatal(err)
	}
	req = req.WithContext(ctx)
	client := &http.Client{Transport: &TraceTransport{}}
	if _, err = client.Do(req); err != nil {
		t.Fatal(err)
	}
	<-wait
}

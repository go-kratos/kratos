package http

import (
	http "net/http"
	"sync"
	"testing"
	"time"

	bm "go-common/library/net/http/blademaster"
	xtime "go-common/library/time"
)

var once sync.Once

func startServer() {
	c := &ServerConfig{
		Addr:         "localhost:18080",
		Timeout:      xtime.Duration(time.Second),
		ReadTimeout:  xtime.Duration(time.Second),
		WriteTimeout: xtime.Duration(time.Second),
	}
	engine := bm.Default()
	engine.GET("/test", func(ctx *bm.Context) {
		ctx.JSON("", nil)
	})
	Serve(engine, c)
}

func TestServer2(t *testing.T) {
	once.Do(startServer)
	resp, err := http.Get("http://localhost:18080/test")
	if err != nil {
		t.Errorf("HTTPServ: get error(%v)", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("http.Get get error code:%d", resp.StatusCode)
	}
	resp.Body.Close()
}

func BenchmarkServer2(b *testing.B) {
	once.Do(startServer)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			resp, err := http.Get("http://localhost:18080/test")
			if err != nil {
				b.Errorf("HTTPServ: get error(%v)", err)
				continue
			}
			if resp.StatusCode != http.StatusOK {
				b.Errorf("HTTPServ: get error status code:%d", resp.StatusCode)
			}
			resp.Body.Close()
		}
	})
}

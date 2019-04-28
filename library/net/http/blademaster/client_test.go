package blademaster

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
	"testing"
	"time"

	"go-common/library/ecode"
	"go-common/library/net/http/blademaster/tests"
	"go-common/library/net/netutil/breaker"
	xtime "go-common/library/time"
)

func TestClient(t *testing.T) {
	c := &ServerConfig{
		Addr:         "localhost:8081",
		Timeout:      xtime.Duration(time.Second),
		ReadTimeout:  xtime.Duration(time.Second),
		WriteTimeout: xtime.Duration(time.Second),
	}
	engine := Default()
	engine.GET("/mytest", func(ctx *Context) {
		time.Sleep(time.Millisecond * 500)
		ctx.JSON("", nil)
	})
	engine.GET("/mytest1", func(ctx *Context) {
		time.Sleep(time.Millisecond * 500)
		ctx.JSON("", nil)
	})
	engine.SetConfig(c)
	engine.Start()

	client := NewClient(
		&ClientConfig{
			App: &App{
				Key:    "53e2fa226f5ad348",
				Secret: "3cf6bd1b0ff671021da5f424fea4b04a",
			},
			Dial:      xtime.Duration(time.Second),
			Timeout:   xtime.Duration(time.Second),
			KeepAlive: xtime.Duration(time.Second),
			Breaker: &breaker.Config{
				Window:  10 * xtime.Duration(time.Second),
				Sleep:   50 * xtime.Duration(time.Millisecond),
				Bucket:  10,
				Ratio:   0.5,
				Request: 100,
			},
		})
	var res struct {
		Code int `json:"code"`
	}
	// test Get
	if err := client.Get(context.Background(), "http://api.bilibili.com/x/server/now", "", nil, &res); err != nil {
		t.Errorf("HTTPClient: expected no error but got %v", err)
	}
	if res.Code != 0 {
		t.Errorf("HTTPClient: expected code=0 but got %d", res.Code)
	}
	// test Post
	if err := client.Post(context.Background(), "http://api.bilibili.com/x/server/now", "", nil, &res); err != nil {
		t.Errorf("HTTPClient: expected no error but got %v", err)
	}
	if res.Code != -405 {
		t.Errorf("HTTPClient: expected code=-405 but got %d", res.Code)
	}
	// test DialTimeout 172.168.1.1 can't connect.
	client.SetConfig(&ClientConfig{Dial: xtime.Duration(time.Second * 5)})
	if err := client.Post(context.Background(), "http://172.168.1.1/x/server/now", "", nil, &res); err == nil {
		t.Errorf("HTTPClient: expected error but got %v", err)
	}
	// test server and timeout.
	client.SetConfig(&ClientConfig{KeepAlive: xtime.Duration(time.Second * 20), Timeout: xtime.Duration(time.Millisecond * 400)})
	if err := client.Get(context.Background(), "http://localhost:8081/mytest", "", nil, &res); err == nil {
		t.Errorf("HTTPClient: expected error timeout for request")
	}
	client.SetConfig(&ClientConfig{Timeout: xtime.Duration(time.Second),
		URL: map[string]*ClientConfig{"http://localhost:8081/mytest1": {Timeout: xtime.Duration(time.Millisecond * 300)}}})
	if err := client.Get(context.Background(), "http://localhost:8081/mytest", "", nil, &res); err != nil {
		t.Errorf("HTTPClient: expected no error but got %v", err)
	}
	if err := client.Get(context.Background(), "http://localhost:8081/mytest1", "", nil, &res); err == nil {
		t.Errorf("HTTPClient: expected error timeout for path")
	}
	client.SetConfig(&ClientConfig{
		Host: map[string]*ClientConfig{"api.bilibili.com": {Timeout: xtime.Duration(time.Millisecond * 300)}},
	})
	if err := client.Get(context.Background(), "http://api.bilibili.com/x/server/now", "", nil, &res); err != nil {
		t.Errorf("HTTPClient: expected no error but got %v", err)
	}
	client.SetConfig(&ClientConfig{
		Host: map[string]*ClientConfig{"api.bilibili.com": {Timeout: xtime.Duration(time.Millisecond * 1)}},
	})
	if err := client.Get(context.Background(), "http://api.bilibili.com/x/server/now", "", nil, &res); err == nil {
		t.Errorf("HTTPClient: expected error timeout but got %v", err)
	}
	client.SetConfig(&ClientConfig{KeepAlive: xtime.Duration(time.Second * 70)})
}

func TestDo(t *testing.T) {
	var (
		aid    = 5463320
		uri    = "http://api.bilibili.com/x/server/now"
		req    *http.Request
		client *Client
		err    error
	)
	client = NewClient(
		&ClientConfig{
			App: &App{
				Key:    "53e2fa226f5ad348",
				Secret: "3cf6bd1b0ff671021da5f424fea4b04a",
			},
			Dial:      xtime.Duration(time.Second),
			Timeout:   xtime.Duration(time.Second),
			KeepAlive: xtime.Duration(time.Second),
			Breaker: &breaker.Config{
				Window:  10 * xtime.Duration(time.Second),
				Sleep:   50 * xtime.Duration(time.Millisecond),
				Bucket:  10,
				Ratio:   0.5,
				Request: 100,
			},
		})
	params := url.Values{}
	params.Set("aid", strconv.Itoa(aid))
	if req, err = client.NewRequest("GET", uri, "", params); err != nil {
		t.Errorf("client.NewRequest: get error(%v)", err)
	}
	var res struct {
		Code int `json:"code"`
	}
	if err = client.Do(context.TODO(), req, &res); err != nil {
		t.Errorf("Do: client.Do get error(%v) url: %s", err, realURL(req))
	}
}

func BenchmarkDo(b *testing.B) {
	once.Do(startServer)
	cf := &ClientConfig{
		App: &App{
			Key:    "53e2fa226f5ad348",
			Secret: "3cf6bd1b0ff671021da5f424fea4b04a",
		},
		Dial:      xtime.Duration(time.Second),
		Timeout:   xtime.Duration(time.Second),
		KeepAlive: xtime.Duration(time.Second),
		Breaker: &breaker.Config{
			Window:  1 * xtime.Duration(time.Second),
			Sleep:   5 * xtime.Duration(time.Millisecond),
			Bucket:  1,
			Ratio:   0.5,
			Request: 10,
		},
		URL: map[string]*ClientConfig{
			"http://api.bilibili.com/x/server/now":  {Timeout: xtime.Duration(time.Second)},
			"http://api.bilibili.com/x/server/nowx": {Timeout: xtime.Duration(time.Second)},
		},
	}
	client := NewClient(cf)
	uri := "http://api.bilibili.com/x/server/now"
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			// client.SetConfig(cf)
			req, err := client.NewRequest("GET", uri, "", nil)
			if err != nil {
				b.Errorf("newRequest: get error(%v)", err)
				continue
			}
			var res struct {
				Code int `json:"code"`
			}
			if err = client.Do(context.TODO(), req, &res); err != nil {
				b.Errorf("Do: client.Do get error(%v) url: %s", err, realURL(req))
			}
		}
	})
	uri = "http://api.bilibili.com/x/server/nowx" // NOTE: for breaker
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			// client.SetConfig(cf)
			req, err := client.NewRequest("GET", uri, "", nil)
			if err != nil {
				b.Errorf("newRequest: get error(%v)", err)
				continue
			}
			var res struct {
				Code int `json:"code"`
			}
			if err = client.Do(context.TODO(), req, &res); err != nil {
				if ecode.ServiceUnavailable.Equal(err) {
					b.Logf("Do: client.Do get error(%v) url: %s", err, realURL(req))
				}
			}
		}
	})
}

func TestRESTfulClient(t *testing.T) {
	c := &ServerConfig{
		Addr:         "localhost:8082",
		Timeout:      xtime.Duration(time.Second),
		ReadTimeout:  xtime.Duration(time.Second),
		WriteTimeout: xtime.Duration(time.Second),
	}
	engine := Default()
	engine.GET("/mytest/1", func(ctx *Context) {
		time.Sleep(time.Millisecond * 500)
		ctx.JSON("", nil)
	})
	engine.GET("/mytest/2/1", func(ctx *Context) {
		time.Sleep(time.Millisecond * 500)
		// ctx.AbortWithStatus(http.StatusInternalServerError)
		ctx.JSON(nil, ecode.ServerErr)
	})
	engine.SetConfig(c)
	engine.Start()

	client := NewClient(
		&ClientConfig{
			App: &App{
				Key:    "53e2fa226f5ad348",
				Secret: "3cf6bd1b0ff671021da5f424fea4b04a",
			},
			Dial:      xtime.Duration(time.Second),
			Timeout:   xtime.Duration(time.Second),
			KeepAlive: xtime.Duration(time.Second),
			Breaker: &breaker.Config{
				Window:  10 * xtime.Duration(time.Second),
				Sleep:   50 * xtime.Duration(time.Millisecond),
				Bucket:  10,
				Ratio:   0.5,
				Request: 100,
			},
		})

	var res struct {
		Code int `json:"code"`
	}
	if err := client.RESTfulGet(context.Background(), "http://localhost:8082/mytest/%d", "", nil, &res, 1); err != nil {
		t.Errorf("HTTPClient: expected error RESTfulGet err: %v", err)
	}
	if res.Code != 0 {
		t.Errorf("HTTPClient: expected code=0 but got %d", res.Code)
	}

	if err := client.RESTfulGet(context.Background(), "http://localhost:8082/mytest/%d/%d", "", nil, &res, 2, 1); err != nil {
		t.Errorf("HTTPClient: expected error RESTfulGet err: %v", err)
	}
	if res.Code != -500 {
		t.Errorf("HTTPClient: expected code=-500 but got %d", res.Code)
	}
}

func TestRaw(t *testing.T) {
	var (
		aid    = 5463320
		uri    = "http://api.bilibili.com/x/server/now"
		req    *http.Request
		client *Client
		err    error
	)
	client = NewClient(
		&ClientConfig{
			App: &App{
				Key:    "53e2fa226f5ad348",
				Secret: "3cf6bd1b0ff671021da5f424fea4b04a",
			},
			Dial:      xtime.Duration(time.Second),
			Timeout:   xtime.Duration(time.Second),
			KeepAlive: xtime.Duration(time.Second),
			Breaker: &breaker.Config{
				Window:  10 * xtime.Duration(time.Second),
				Sleep:   50 * xtime.Duration(time.Millisecond),
				Bucket:  10,
				Ratio:   0.5,
				Request: 100,
			},
		})
	params := url.Values{}
	params.Set("aid", strconv.Itoa(aid))
	if req, err = client.NewRequest("GET", uri, "", params); err != nil {
		t.Errorf("client.NewRequest: get error(%v)", err)
	}
	var (
		bs []byte
	)
	if bs, err = client.Raw(context.TODO(), req); err != nil {
		t.Errorf("Do: client.Do get error(%v) url: %s", err, realURL(req))
	}
	t.Log(string(bs))
}

func TestJSON(t *testing.T) {
	var (
		aid    = 5463320
		uri    = "http://api.bilibili.com/x/server/now"
		req    *http.Request
		client *Client
		err    error
	)
	client = NewClient(
		&ClientConfig{
			App: &App{
				Key:    "53e2fa226f5ad348",
				Secret: "3cf6bd1b0ff671021da5f424fea4b04a",
			},
			Dial:      xtime.Duration(time.Second),
			Timeout:   xtime.Duration(time.Second),
			KeepAlive: xtime.Duration(time.Second),
			Breaker: &breaker.Config{
				Window:  10 * xtime.Duration(time.Second),
				Sleep:   50 * xtime.Duration(time.Millisecond),
				Bucket:  10,
				Ratio:   0.5,
				Request: 100,
			},
		})
	params := url.Values{}
	params.Set("aid", strconv.Itoa(aid))
	if req, err = client.NewRequest("GET", uri, "", params); err != nil {
		t.Errorf("client.NewRequest: get error(%v)", err)
	}
	var res struct {
		Code int `json:"code"`
	}
	if err = client.Do(context.TODO(), req, &res); err != nil {
		t.Errorf("Do: client.Do get error(%v) url: %s", err, realURL(req))
	}
}

func TestPB(t *testing.T) {
	var (
		uri    = "http://172.18.33.143:13500/playurl/batch"
		req    *http.Request
		client *Client
		err    error
	)
	client = NewClient(
		&ClientConfig{
			App: &App{
				Key:    "53e2fa226f5ad348",
				Secret: "3cf6bd1b0ff671021da5f424fea4b04a",
			},
			Dial:      xtime.Duration(time.Second),
			Timeout:   xtime.Duration(time.Second),
			KeepAlive: xtime.Duration(time.Second),
			Breaker: &breaker.Config{
				Window:  10 * xtime.Duration(time.Second),
				Sleep:   50 * xtime.Duration(time.Millisecond),
				Bucket:  10,
				Ratio:   0.5,
				Request: 100,
			},
		})
	params := url.Values{}
	params.Set("cid", "10108859,10108860")
	params.Set("uip", "222.73.196.18")
	params.Set("qn", "16")
	params.Set("platform", "html5")
	params.Set("layout", "pb")
	if req, err = client.NewRequest("GET", uri, "", params); err != nil {
		t.Errorf("client.NewRequest: get error(%v)", err)
	}
	var res = new(tests.BvcResponseMsg)
	if err = client.PB(context.TODO(), req, res); err != nil {
		t.Errorf("Do: client.Do get error(%v) url: %s", err, realURL(req))
	}
	t.Log(res)
}

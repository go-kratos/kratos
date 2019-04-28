package rate

import (
	"context"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	bm "go-common/library/net/http/blademaster"
)

func TestLimiterUrl(t *testing.T) {
	l := New(&Config{
		URLs: map[string]*Limit{"/limit/test": &Limit{Limit: 1, Burst: 2}},
	})
	if !l.Allow("testApp", "/limit/test") {
		t.Logf("request should pass,but blocked")
		t.FailNow()
	}
	if !l.Allow("testApp", "/limit/test") {
		t.Logf("request should pass,but blocked")
		t.FailNow()
	}
	if l.Allow("testApp", "/limit/test") {
		t.Logf("request should block,but passed")
		t.FailNow()
	}
}

func TestLimiterApp(t *testing.T) {
	l := New(&Config{
		Apps: map[string]*Limit{"testApp": &Limit{Limit: 1, Burst: 2}},
	})
	if !l.Allow("testApp", "/limit/test") {
		t.Logf("request should pass,but blocked")
		t.FailNow()
	}
	if !l.Allow("testApp", "/limit/test") {
		t.Logf("request should pass,but blocked")
		t.FailNow()
	}
	if l.Allow("testApp", "/limit/test") {
		t.Logf("request should block,but passed")
		t.FailNow()
	}
}

func TestLimiterUrlApp(t *testing.T) {
	l := New(&Config{
		Apps: map[string]*Limit{"testApp": &Limit{Limit: 2, Burst: 1}},
		URLs: map[string]*Limit{"/limit/test": &Limit{Limit: 2, Burst: 1}},
	})
	if !l.Allow("testApp", "/limit/test") {
		t.Logf("request should pass,but blocked")
		t.FailNow()
	}
	if l.Allow("testApp", "/limit/test") {
		t.Logf("request should block,but passed")
		t.FailNow()
	}
	l.Reload(&Config{
		Apps: map[string]*Limit{"testApp": &Limit{Limit: 1, Burst: 2}},
		URLs: map[string]*Limit{"/limit/test": &Limit{Limit: 1, Burst: 2}},
	})
	if !l.Allow("testApp", "/limit/test") {
		t.Logf("request should pass,but blocked")
		t.FailNow()
	}
	if !l.Allow("testApp", "/limit/test") {
		t.Logf("request should pass,but blocked")
		t.FailNow()
	}
	if l.Allow("testApp", "/limit/test") {
		t.Logf("request should block,but passed")
		t.FailNow()
	}
}

func TestLimiterHandler(t *testing.T) {
	l := New(&Config{
		Apps: map[string]*Limit{"testApp": &Limit{Limit: 1, Burst: 1}},
		URLs: map[string]*Limit{"/limit/test": &Limit{Limit: 2, Burst: 4}},
	})
	engine := bm.New()
	engine.Use(l.Handler())
	engine.GET("/limit/test", func(c *bm.Context) {
		c.String(200, "pass")
	})
	go engine.Run(":18080")
	defer func() {
		engine.Server().Shutdown(context.TODO())
	}()

	time.Sleep(time.Millisecond * 20)
	code, content, err := httpGet("http://127.0.0.1:18080/limit/test?appkey=testApp")
	if err != nil {
		t.Logf("http get failed,err:=%v", err)
		t.FailNow()
	}
	if code != 200 || string(content) != "pass" {
		t.Logf("request should pass by limiter,but blocked")
		t.FailNow()
	}

	_, content, err = httpGet("http://127.0.0.1:18080/limit/test?appkey=testApp")
	if err != nil {
		t.Logf("http get failed,err:=%v", err)
		t.FailNow()
	}
	if string(content) == "pass" {
		t.Logf("request should block by limiter,but passed")
		t.FailNow()
	}
}

func httpGet(url string) (code int, content []byte, err error) {
	resp, err := http.Get(url)
	if err != nil {

		return
	}
	defer resp.Body.Close()

	content, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	code = resp.StatusCode
	return
}

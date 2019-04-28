package permit

import (
	"context"
	"net/http"
	"net/url"
	"sync"
	"testing"
	"time"

	"go-common/library/cache/memcache"
	"go-common/library/container/pool"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/netutil/breaker"
	xtime "go-common/library/time"
)

var (
	once sync.Once
)

type Response struct {
	Code int    `json:"code"`
	Data string `json:"data"`
}

func init() {
	log.Init(nil)
}

func client() *bm.Client {
	return bm.NewClient(&bm.ClientConfig{
		App: &bm.App{
			Key:    "test",
			Secret: "test",
		},
		Dial:      xtime.Duration(time.Second),
		Timeout:   xtime.Duration(time.Second),
		KeepAlive: xtime.Duration(time.Second * 10),
		Breaker: &breaker.Config{
			Window:  xtime.Duration(time.Second),
			Sleep:   xtime.Duration(time.Millisecond * 100),
			Bucket:  10,
			Ratio:   0.5,
			Request: 100,
		},
	})
}

func getPermit() *Permit {
	return New(&Config{
		DsHTTPClient: &bm.ClientConfig{
			App: &bm.App{
				Key:    "manager-go",
				Secret: "949bbb2dd3178252638c2407578bc7ad",
			},
			Dial:      xtime.Duration(time.Second),
			Timeout:   xtime.Duration(time.Second),
			KeepAlive: xtime.Duration(time.Second * 10),
			Breaker: &breaker.Config{
				Window:  xtime.Duration(time.Second),
				Sleep:   xtime.Duration(time.Millisecond * 100),
				Bucket:  10,
				Ratio:   0.5,
				Request: 100,
			},
		},
		MaHTTPClient: &bm.ClientConfig{
			App: &bm.App{
				Key:    "f6433799dbd88751",
				Secret: "36f8ddb1806207fe07013ab6a77a3935",
			},
			Dial:      xtime.Duration(time.Second),
			Timeout:   xtime.Duration(time.Second),
			KeepAlive: xtime.Duration(time.Second * 10),
			Breaker: &breaker.Config{
				Window:  xtime.Duration(time.Second),
				Sleep:   xtime.Duration(time.Millisecond * 100),
				Bucket:  10,
				Ratio:   0.5,
				Request: 100,
			},
		},
		Session: &SessionConfig{
			SessionIDLength: 32,
			CookieLifeTime:  1800,
			CookieName:      "mng-go",
			Domain:          ".bilibili.co",
			Memcache: &memcache.Config{
				Config: &pool.Config{
					Active:      10,
					Idle:        5,
					IdleTimeout: xtime.Duration(time.Second * 80),
				},
				Name:         "go-business/auth",
				Proto:        "tcp",
				Addr:         "172.16.33.54:11211",
				DialTimeout:  xtime.Duration(time.Millisecond * 1000),
				ReadTimeout:  xtime.Duration(time.Millisecond * 1000),
				WriteTimeout: xtime.Duration(time.Millisecond * 1000),
			},
		},
		ManagerHost:     "http://uat-manager.bilibili.co",
		DashboardHost:   "http://dashboard-mng.bilibili.co",
		DashboardCaller: "manager-go",
	})
}

func engine() *bm.Engine {
	e := bm.NewServer(nil)
	a := getPermit()

	e.GET("/login", a.Verify(), func(c *bm.Context) {
		c.JSON("pass", nil)
	})
	e.GET("/tag/del", a.Permit("TAG_DEL"), func(c *bm.Context) {
		c.JSON("pass", nil)
	})
	e.GET("/tag/admin", a.Permit("TAG_ADMIN"), func(c *bm.Context) {
		c.JSON("pass", nil)
	})
	return e
}

func setSession(uid int64, username string) (string, error) {
	a := getPermit()
	sv := a.sm.newSession(context.TODO())
	sv.Set("username", username)
	mcConn := a.sm.mc.Get(context.TODO())
	defer mcConn.Close()
	key := sv.Sid
	item := &memcache.Item{
		Key:        key,
		Object:     sv,
		Flags:      memcache.FlagJSON,
		Expiration: int32(a.sm.c.CookieLifeTime),
	}
	if err := mcConn.Set(item); err != nil {
		return "", err
	}
	return key, nil
}

func startEngine(t *testing.T) func() {
	return func() {
		e := engine()
		err := e.Run(":18080")
		if err != nil {
			t.Fatalf("failed to run server!%v", err)
		}
	}
}

func TestLoginSuccess(t *testing.T) {
	go once.Do(startEngine(t))
	time.Sleep(time.Millisecond * 100)

	sid, err := setSession(2233, "caoguoliang")
	if err != nil {
		t.Fatalf("faild to set session !err:=%v", err)
	}
	query := url.Values{}
	query.Set("test", "test")
	cli := client()
	req, err := cli.NewRequest("GET", "http://127.0.0.1:18080/login", "", query)
	if err != nil {
		t.Fatalf("Failed to build request: %v", err)
	}
	req.AddCookie(&http.Cookie{
		Name:  "mng-go",
		Value: sid,
	})
	req.AddCookie(&http.Cookie{
		Name:  "username",
		Value: "caoguoliang",
	})
	req.AddCookie(&http.Cookie{
		Name:  "_AJSESSIONID",
		Value: "87fa8450e93511e79ed8522233007f8a",
	})
	res := Response{}
	if err := cli.Do(context.TODO(), req, &res); err != nil {
		t.Fatalf("Failed to send request: %v", err)
	}
	if res.Code != 0 || res.Data != "pass" {
		t.Fatalf("Unexpected response code(%d) data(%v)", res.Code, res.Data)
	}
}

func TestLoginFail(t *testing.T) {
	go once.Do(startEngine(t))
	time.Sleep(time.Millisecond * 100)

	query := url.Values{}
	query.Set("test", "test")
	cli := client()
	req, err := cli.NewRequest("GET", "http://127.0.0.1:18080/login", "", query)
	if err != nil {
		t.Fatalf("Failed to build request: %v", err)
	}
	req.AddCookie(&http.Cookie{
		Name:  "mng-go",
		Value: "fakesess",
	})
	req.AddCookie(&http.Cookie{
		Name:  "username",
		Value: "caoguoliang",
	})
	req.AddCookie(&http.Cookie{
		Name:  "_AJSESSIONID",
		Value: "testsess",
	})
	res := Response{}
	if err := cli.Do(context.TODO(), req, &res); err != nil {
		t.Fatalf("Failed to send request: %v", err)
	}
	if res.Code != ecode.Unauthorized.Code() {
		t.Fatalf("This request should be forbidden: code(%d) data(%v)", res.Code, res.Data)
	}
}

func TestVerifySuccess(t *testing.T) {
	go once.Do(startEngine(t))
	time.Sleep(time.Millisecond * 100)

	sid, err := setSession(2233, "caoguoliang")
	if err != nil {
		t.Fatalf("faild to set session !err:=%v", err)
	}
	query := url.Values{}
	query.Set("test", "test")
	cli := client()
	req, err := cli.NewRequest("GET", "http://127.0.0.1:18080/tag/del", "", query)
	if err != nil {
		t.Fatalf("Failed to build request: %v", err)
	}
	req.AddCookie(&http.Cookie{
		Name:  "mng-go",
		Value: sid,
	})
	req.AddCookie(&http.Cookie{
		Name:  "username",
		Value: "caoguoliang",
	})
	req.AddCookie(&http.Cookie{
		Name:  "_AJSESSIONID",
		Value: "87fa8450e93511e79ed8522233007f8a",
	})
	res := Response{}
	if err := cli.Do(context.TODO(), req, &res); err != nil {
		t.Fatalf("Failed to send request: %v", err)
	}
	if res.Code != 0 || res.Data != "pass" {
		t.Fatalf("Unexpected response code(%d) data(%v)", res.Code, res.Data)
	}
}

func TestVerifyFail(t *testing.T) {
	go once.Do(startEngine(t))
	time.Sleep(time.Millisecond * 100)

	sid, err := setSession(2233, "caoguoliang")
	if err != nil {
		t.Fatalf("faild to set session !err:=%v", err)
	}
	query := url.Values{}
	query.Set("test", "test")
	cli := client()
	req, err := cli.NewRequest("GET", "http://127.0.0.1:18080/tag/admin", "", query)
	if err != nil {
		t.Fatalf("Failed to build request: %v", err)
	}
	req.AddCookie(&http.Cookie{
		Name:  "mng-go",
		Value: sid,
	})
	req.AddCookie(&http.Cookie{
		Name:  "username",
		Value: "caoguoliang",
	})
	req.AddCookie(&http.Cookie{
		Name:  "_AJSESSIONID",
		Value: "87fa8450e93511e79ed8522233007f8a",
	})
	res := Response{}
	if err := cli.Do(context.TODO(), req, &res); err != nil {
		t.Fatalf("Failed to send request: %v", err)
	}
	if res.Code != ecode.AccessDenied.Code() {
		t.Fatalf("This request should be forbidden: code(%d) data(%v)", res.Code, res.Data)
	}
}

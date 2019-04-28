package permit_test

import (
	"fmt"
	"time"

	"go-common/library/cache/memcache"
	"go-common/library/container/pool"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/permit"
	"go-common/library/net/metadata"
	"go-common/library/net/netutil/breaker"
	xtime "go-common/library/time"
)

// This example create a permit middleware instance and attach to several path,
// it will validate request by specified policy and put extra information into context. e.g., `uid`.
// It provides additional handler functions to provide the identification for your business handler.
func Example() {
	a := permit.New(&permit.Config{
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
		Session: &permit.SessionConfig{
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
				Name:         "go-business/permit",
				Proto:        "tcp",
				Addr:         "172.16.33.54:11211",
				DialTimeout:  xtime.Duration(time.Millisecond * 1000),
				ReadTimeout:  xtime.Duration(time.Millisecond * 1000),
				WriteTimeout: xtime.Duration(time.Millisecond * 1000),
			},
		},
		ManagerHost:     "http://uat-manager.bilibili.co",
		DashboardHost:   "http://uat-dashboard-mng.bilibili.co",
		DashboardCaller: "manager-go",
	})

	p := permit.New2(nil)

	e := bm.NewServer(nil)

	//Check whether the user has logged in
	e.GET("/login", a.Verify(), func(c *bm.Context) {
		c.JSON("pass", nil)
	})
	//Check whether the user has logged in,and check th user has the access permisson to the specifed path
	e.GET("/tag/del", a.Permit("TAG_DEL"), func(c *bm.Context) {
		uid := metadata.Int64(c, metadata.Uid)
		username := metadata.String(c, metadata.Username)
		c.JSON(fmt.Sprintf("pass uid(%d) username(%s)", uid, username), nil)
	})
	e.GET("/check/login", p.Verify2(), func(c *bm.Context) {
		c.JSON("pass", nil)
	})
	e.POST("/tag/del", p.Permit2("TAG_DEL"), func(c *bm.Context) {
		uid := metadata.Int64(c, metadata.Uid)
		username := metadata.String(c, metadata.Username)
		c.JSON(fmt.Sprintf("pass uid(%d) username(%s)", uid, username), nil)
	})

	e.Run(":18080")
}

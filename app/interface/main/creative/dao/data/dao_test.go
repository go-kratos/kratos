package data

import (
	"context"
	"flag"
	"go-common/app/interface/main/creative/conf"
	"go-common/library/cache/memcache"
	"go-common/library/log"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/bouk/monkey"
	"github.com/smartystreets/goconvey/convey"
	gock "gopkg.in/h2non/gock.v1"
)

var (
	d *Dao
)

func TestMain(m *testing.M) {
	if os.Getenv("DEPLOY_ENV") != "" {
		flag.Set("app_id", "main.archive.creative")
		flag.Set("conf_token", "96b6a6c10bb311e894c14a552f48fef8")
		flag.Set("tree_id", "2305")
		flag.Set("conf_version", "docker-1")
		flag.Set("deploy_env", "uat")
		flag.Set("conf_host", "config.bilibili.co")
		flag.Set("conf_path", "/tmp")
		flag.Set("region", "sh")
		flag.Set("zone", "sh001")
	} else {
		flag.Set("conf", "../../cmd/creative.toml")
	}
	flag.Parse()
	if err := conf.Init(); err != nil {
		panic(err)
	}
	d = New(conf.Conf)
	m.Run()
	os.Exit(0)
}

func httpMock(method, url string) *gock.Request {
	r := gock.New(url)
	r.Method = strings.ToUpper(method)
	d.client.SetTransport(gock.DefaultTransport)
	d.statClient.SetTransport(gock.DefaultTransport)
	return r
}
func TestStat(t *testing.T) {
	var (
		c   = context.TODO()
		mid = int64(1)
		ip  = "127.0.0.1"
		dt  = "iamdt"
	)
	convey.Convey("TestStat1", t, func(ctx convey.C) {
		connGuard := monkey.PatchInstanceMethod(reflect.TypeOf(d.mc), "Get", func(_ *memcache.Pool, _ context.Context) memcache.Conn {
			return memcache.MockWith(memcache.ErrNotFound)
		})
		defer connGuard.Unpatch()
		_, err := d.Stat(c, ip, mid)
		ctx.Convey("Then err should be nil.stat should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldNotBeNil)
		})
	})
	convey.Convey("TestStat2", t, func(ctx convey.C) {
		connGuard := monkey.PatchInstanceMethod(reflect.TypeOf(d.mc), "Get", func(_ *memcache.Pool, _ context.Context) memcache.Conn {
			return memcache.MockWith(memcache.ErrNotFound)
		})
		defer connGuard.Unpatch()
		_, err := d.UpStat(c, mid, dt)
		ctx.Convey("Then err should be nil.stat should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldNotBeNil)
		})
	})
	convey.Convey("Ping", t, func(ctx convey.C) {
		connGuard := monkey.PatchInstanceMethod(reflect.TypeOf(d.mc), "Get", func(_ *memcache.Pool, _ context.Context) memcache.Conn {
			return memcache.MockWith(memcache.ErrNotFound)
		})
		defer connGuard.Unpatch()
		err := d.Ping(c)
		ctx.Convey("Then err should be nil.stat should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldNotBeNil)
		})
	})
	convey.Convey("Close", t, func(ctx convey.C) {
		connGuard := monkey.PatchInstanceMethod(reflect.TypeOf(d.mc), "Close", func(_ *memcache.Pool) error {
			return nil
		})
		defer connGuard.Unpatch()
		err := d.Close()
		ctx.Convey("Then err should be nil.stat should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
	convey.Convey("AddCache", t, func(ctx convey.C) {
		d.AddCache(func() {
			log.Error("d.AddCache")
		})
	})
}

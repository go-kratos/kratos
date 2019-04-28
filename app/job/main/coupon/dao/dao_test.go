package dao

import (
	"context"
	"flag"
	"os"
	"strings"
	"testing"

	"github.com/smartystreets/goconvey/convey"

	"go-common/app/job/main/coupon/conf"

	gock "gopkg.in/h2non/gock.v1"
)

var (
	d *Dao
)

func TestMain(m *testing.M) {
	if os.Getenv("DEPLOY_ENV") != "" {
		flag.Set("app_id", "main.account.coupon-job")
		flag.Set("conf_token", "ddd3dd8ed06d0ca10b6c8f122581d035")
		flag.Set("tree_id", "23029")
		flag.Set("conf_version", "docker-1")
		flag.Set("deploy_env", "uat")
		flag.Set("conf_host", "config.bilibili.co")
		flag.Set("conf_path", "/tmp")
		flag.Set("region", "sh")
		flag.Set("zone", "sh001")
	} else {
		flag.Set("conf", "../cmd/coupon-job.toml")
	}
	flag.Parse()
	if err := conf.Init(); err != nil {
		panic(err)
	}
	d = New(conf.Conf)
	d.client.SetTransport(gock.DefaultTransport)
	os.Exit(m.Run())
}

func httpMock(method, url string) *gock.Request {
	r := gock.New(url)
	r.Method = strings.ToUpper(method)
	return r
}

func TestDaoPing(t *testing.T) {
	convey.Convey("TestDaoPing", t, func(convCtx convey.C) {
		err := d.Ping(context.Background())
		convCtx.Convey("Then err should be nil.", func(convCtx convey.C) {
			convCtx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaopingMC(t *testing.T) {
	convey.Convey("TestDaopingMC", t, func(convCtx convey.C) {
		err := d.pingMC(context.Background())
		convCtx.Convey("Then err should be nil.", func(convCtx convey.C) {
			convCtx.So(err, convey.ShouldBeNil)
		})
	})
}

package monitor

import (
	"context"
	"flag"
	"go-common/app/service/main/up/conf"
	"net/url"
	"os"
	"testing"

	"github.com/smartystreets/goconvey/convey"
	"go-common/app/service/main/up/dao"
)

var (
	d *Dao
)

func TestMain(m *testing.M) {
	if os.Getenv("DEPLOY_ENV") != "" {
		flag.Set("app_id", dao.AppID)
		flag.Set("conf_token", dao.UatToken)
		flag.Set("tree_id", dao.TreeID)
		flag.Set("conf_version", "docker-1")
		flag.Set("deploy_env", "uat")
		flag.Set("conf_host", "config.bilibili.co")
		flag.Set("conf_path", "/tmp")
		flag.Set("region", "sh")
		flag.Set("zone", "sh001")
	} else {
		flag.Set("conf", "../../cmd/up-service.toml")
	}
	if os.Getenv("UT_LOCAL_TEST") != "" {
		flag.Set("conf", "../../cmd/up-service.toml")
	}
	flag.Parse()
	if err := conf.Init(); err != nil {
		panic(err)
	}
	d = New(conf.Conf)
	m.Run()
	os.Exit(0)
}

func TestMonitorSend(t *testing.T) {
	var (
		c        = context.TODO()
		username = ""
		msg      = ""
	)
	convey.Convey("Send", t, func(ctx convey.C) {
		err := d.Send(c, username, msg)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestMonitorgetSign(t *testing.T) {
	var (
		params url.Values
	)
	convey.Convey("getSign", t, func(ctx convey.C) {
		sign := d.getSign(params)
		ctx.Convey("Then sign should not be nil.", func(ctx convey.C) {
			ctx.So(sign, convey.ShouldNotBeNil)
		})
	})
}

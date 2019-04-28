package dao

import (
	"context"
	"flag"
	"os"
	"testing"

	"go-common/app/service/main/identify/conf"

	"github.com/smartystreets/goconvey/convey"
)

var (
	d *Dao
)

func TestMain(m *testing.M) {
	if os.Getenv("DEPLOY_ENV") != "" {
		flag.Set("app_id", "main.account.identify-service")
		flag.Set("conf_appid", "identify-service")
		flag.Set("conf_token", "ed882ffb3c5f9ec3f4a28691a2ac3d84")
		flag.Set("tree_id", "11172")
		flag.Set("conf_version", "docker-1")
		flag.Set("deploy_env", "uat")
		flag.Set("conf_host", "config.bilibili.co")
		flag.Set("conf_path", "/tmp")
		flag.Set("region", "sh")
		flag.Set("zone", "sh001")
	} else {
		flag.Set("app_id", "main.account.identify-service")
		flag.Set("conf_appid", "identify-service")
		flag.Set("conf_token", "15fd302391c22077da13643e11d8bab6")
		flag.Set("tree_id", "11172")
		flag.Set("conf_version", "server-1")
		flag.Set("deploy_env", "dev")
		flag.Set("conf_env", "dev")
		flag.Set("conf_host", "config.bilibili.co")
		flag.Set("conf_path", "/tmp")
		flag.Set("region", "sh")
		flag.Set("zone", "sh001")
	}
	flag.Parse()
	if err := conf.Init(); err != nil {
		panic(err)
	}
	d = New(conf.Conf)
	m.Run()
	os.Exit(0)
}

func TestDao_Ping(t *testing.T) {
	convey.Convey("dao ping", t, func(ctx convey.C) {
		err := d.Ping(context.Background())
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDao_pingMC(t *testing.T) {
	convey.Convey("dao ping", t, func(ctx convey.C) {
		err := d.pingMC(context.Background())
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

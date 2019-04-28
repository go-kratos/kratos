package archive

import (
	"context"
	"flag"
	"os"
	"testing"

	"go-common/app/interface/main/app-feed/conf"

	"github.com/smartystreets/goconvey/convey"
)

var (
	d *Dao
)

func TestMain(m *testing.M) {
	if os.Getenv("DEPLOY_ENV") != "" {
		flag.Set("app_id", "main.app-svr.app-feed")
		flag.Set("conf_token", "OC30xxkAOyaH9fI6FRuXA0Ob5HL0f3kc")
		flag.Set("tree_id", "2686")
		flag.Set("conf_version", "docker-1")
		flag.Set("deploy_env", "uat")
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

func TestPing(t *testing.T) {
	var (
		c = context.TODO()
	)
	convey.Convey("Ping", t, func(ctx convey.C) {
		d.PingMC(c)
	})
}

func TestArchives(t *testing.T) {
	var (
		c    = context.TODO()
		aids = []int64{1}
	)
	convey.Convey("Ping", t, func(ctx convey.C) {
		_, err := d.Archives(c, aids)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestArchivesWithPlayer(t *testing.T) {
	var (
		c         = context.TODO()
		aids      = []int64{1}
		qn        = 32
		platform  = "android"
		fnver     = 0
		fnval     = 0
		forceHost = 0
	)
	convey.Convey("Ping", t, func(ctx convey.C) {
		_, err := d.ArchivesWithPlayer(c, aids, qn, platform, fnver, fnval, forceHost, 0)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

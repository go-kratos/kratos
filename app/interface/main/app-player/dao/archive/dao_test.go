package archive

import (
	"context"
	"flag"
	"fmt"
	"os"
	"testing"

	"go-common/app/interface/main/app-player/conf"

	"github.com/smartystreets/goconvey/convey"
)

var (
	d *Dao
)

func TestMain(m *testing.M) {
	if os.Getenv("DEPLOY_ENV") != "" {
		flag.Set("app_id", "main.app-svr.app-player")
		flag.Set("conf_token", "e477d98a7c5689623eca4f32f6af735c")
		flag.Set("tree_id", "52581")
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
		c = context.Background()
	)
	convey.Convey("Ping", t, func(ctx convey.C) {
		err := d.Ping(c)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestArchiveCache(t *testing.T) {
	var (
		c   = context.Background()
		aid = int64(10110670)
	)
	convey.Convey("ArchiveCache", t, func(ctx convey.C) {
		arc, err := d.ArchiveCache(c, aid)
		fmt.Printf("%#v", arc)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestViews(t *testing.T) {
	var (
		c    = context.Background()
		aids = []int64{10110670}
	)
	convey.Convey("Views", t, func(ctx convey.C) {
		_, err := d.Views(c, aids)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

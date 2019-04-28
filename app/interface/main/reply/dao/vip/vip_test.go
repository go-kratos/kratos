package vip

import (
	"context"
	"flag"
	"go-common/app/interface/main/reply/conf"
	"os"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

var (
	d *Dao
)

func TestMain(m *testing.M) {
	if os.Getenv("DEPLOY_ENV") != "" {
		flag.Set("app_id", "")
		flag.Set("conf_token", "")
		flag.Set("tree_id", "")
		flag.Set("conf_version", "docker-1")
		flag.Set("deploy_env", "uat")
		flag.Set("conf_host", "config.bilibili.co")
		flag.Set("conf_path", "/tmp")
		flag.Set("region", "sh")
		flag.Set("zone", "sh001")
	} else {
		flag.Set("conf", "../../cmd/reply-test.toml")
	}
	flag.Parse()
	if err := conf.Init(); err != nil {
		panic(err)
	}
	d = New(conf.Conf)
	os.Exit(m.Run())
}

func TestVipNew(t *testing.T) {
	convey.Convey("New", t, func(ctx convey.C) {
		var (
			c = conf.Conf
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := New(c)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestVipEmoji(t *testing.T) {
	convey.Convey("Emoji", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			emjs, emjM, err := d.Emoji(c)
			ctx.Convey("Then err should be nil.emjs,emjM should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(emjM, convey.ShouldNotBeNil)
				ctx.So(emjs, convey.ShouldNotBeNil)
			})
		})
	})
}

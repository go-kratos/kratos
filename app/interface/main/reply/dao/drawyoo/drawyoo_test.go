package drawyoo

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

func TestDrawyooNew(t *testing.T) {
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

func TestDrawyooInfo(t *testing.T) {
	convey.Convey("Info", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			hid = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			info, err := d.Info(c, hid)
			ctx.Convey("Then err should be nil.info should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldNotBeNil)
				ctx.So(info, convey.ShouldBeNil)
			})
		})
	})
}

func TestDrawyooInfos(t *testing.T) {
	convey.Convey("Infos", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			hids = []int64{}
			d    = New(conf.Conf)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			info, err := d.Infos(c, hids)
			ctx.Convey("Then err should be nil.info should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldNotBeNil)
				ctx.So(info, convey.ShouldBeNil)
			})
		})
	})
}

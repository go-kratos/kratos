package creative

import (
	"context"
	"flag"
	"os"
	"testing"

	"go-common/app/interface/main/app-view/conf"

	"github.com/smartystreets/goconvey/convey"
)

var (
	d *Dao
)

func TestMain(m *testing.M) {
	if os.Getenv("DEPLOY_ENV") != "" {
		flag.Set("app_id", "main.app-svr.app-view")
		flag.Set("conf_token", "3a4CNLBhdFbRQPs7B4QftGvXHtJo92xw")
		flag.Set("tree_id", "4575")
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

func TestMovie(t *testing.T) {
	var (
		c    = context.TODO()
		vmid = int64(0)
	)
	convey.Convey("Movie", t, func(ctx convey.C) {
		_, err := d.FollowSwitch(c, vmid)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			err = nil
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestBgm(t *testing.T) {
	var (
		c   = context.TODO()
		aid = int64(10110995)
		cid = int64(10135146)
	)
	convey.Convey("Bgm", t, func(ctx convey.C) {
		_, err := d.Bgm(c, aid, cid)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			err = nil
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestPoints(t *testing.T) {
	var (
		c   = context.TODO()
		aid = int64(10110995)
		cid = int64(10135146)
	)
	convey.Convey("Points", t, func(ctx convey.C) {
		_, err := d.Points(c, aid, cid)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			err = nil
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestSpecial(t *testing.T) {
	convey.Convey("Special", t, func(ctx convey.C) {
		midsM, err := d.Special(context.Background())
		ctx.So(err, convey.ShouldBeNil)
		ctx.Println(midsM)
	})
}

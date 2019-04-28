package favorite

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
	} else {
		flag.Set("conf", "../../cmd/app-view-test.toml")
	}
	flag.Parse()
	if err := conf.Init(); err != nil {
		panic(err)
	}
	d = New(conf.Conf)
	os.Exit(m.Run())
}

func TestAddVideo(t *testing.T) {
	convey.Convey("AddVideo", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			d.AddVideo(c, 1, []int64{1}, 1, "")
		})
	})
}

func TestIsFavVideo(t *testing.T) {
	convey.Convey("IsFavVideo", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			_, err := d.IsFavVideo(c, 1, 1)
			ctx.Convey("Then err should be nil", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

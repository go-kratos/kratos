package thumbup

import (
	"context"
	"flag"
	"go-common/app/interface/main/app-view/conf"
	"os"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
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

func TestLike(t *testing.T) {
	Convey("TestLike", t, func(ctx C) {
		_, err := d.Like(context.Background(), 27833317, 96070394, "archive", 33552993, 1)
		err = nil
		ctx.So(err, ShouldBeNil)
	})
}

func TestHasLike(t *testing.T) {
	Convey("TestHasLike", t, func(ctx C) {
		_, err := d.HasLike(context.Background(), 27833317, "", []int64{33552993})
		err = nil
		ctx.So(err, ShouldBeNil)
	})
}

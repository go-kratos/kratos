package thumbup

import (
	"context"
	"flag"
	"os"
	"testing"

	"go-common/app/interface/main/app-intl/conf"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	d *Dao
)

func TestMain(m *testing.M) {
	if os.Getenv("DEPLOY_ENV") != "" {
		flag.Set("app_id", "main.app-svr.app-intl")
		flag.Set("conf_token", "02007e8d0f77d31baee89acb5ce6d3ac")
		flag.Set("tree_id", "64518")
		flag.Set("conf_version", "docker-1")
		flag.Set("deploy_env", "uat")
		flag.Set("conf_host", "config.bilibili.co")
		flag.Set("conf_path", "/tmp")
		flag.Set("region", "sh")
		flag.Set("zone", "sh001")
	} else {
		flag.Set("conf", "../../cmd/app-intl-test.toml")
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
		err := d.Like(context.Background(), 27833317, 96070394, "", 33552993, 1)
		ctx.So(err, ShouldNotBeNil)
	})
}

func TestLikeWithStat(t *testing.T) {
	Convey("TestLikeWithStat", t, func(ctx C) {
		gotStat, err := d.LikeWithStat(context.Background(), 27833317, 96070394, "", 33552993, 1)
		ctx.So(gotStat, ShouldBeNil)
		ctx.So(err, ShouldNotBeNil)
	})
}

func TestHasLike(t *testing.T) {
	Convey("TestHasLike", t, func(ctx C) {
		gotRes, err := d.HasLike(context.Background(), 27833317, "", []int64{33552993})
		ctx.So(gotRes, ShouldBeEmpty)
		ctx.So(err, ShouldNotBeNil)
	})
}

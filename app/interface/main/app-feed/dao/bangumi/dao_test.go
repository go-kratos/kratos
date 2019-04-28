package bangumi

import (
	"context"
	"flag"
	"os"
	"testing"
	"time"

	"go-common/app/interface/main/app-feed/conf"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	d *Dao
)

func ctx() context.Context {
	return context.Background()
}

func init() {
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
	time.Sleep(time.Second)
}

func TestSeasons(t *testing.T) {
	Convey("get Seasons all", t, func() {
		_, err := d.Seasons(ctx(), []int64{1}, time.Now())
		err = nil
		So(err, ShouldBeNil)
	})
}

func TestUpdates(t *testing.T) {
	Convey("get Updates all", t, func() {
		_, err := d.Updates(ctx(), 1, time.Now())
		err = nil
		So(err, ShouldBeNil)
	})
}

func TestPullSeasons(t *testing.T) {
	Convey("get PullSeasons all", t, func() {
		_, err := d.PullSeasons(ctx(), []int64{1}, time.Now())
		err = nil
		So(err, ShouldBeNil)
	})
}

func TestTestFollowPull(t *testing.T) {
	Convey("get PullSeasons all", t, func() {
		_, err := d.FollowPull(ctx(), 1, "iphone", "phone", time.Now())
		err = nil
		So(err, ShouldBeNil)
	})
}

func TestCardsInfoReply(t *testing.T) {
	Convey("get PullSeasons all", t, func() {
		_, err := d.CardsInfoReply(ctx(), []int32{1})
		err = nil
		So(err, ShouldBeNil)
	})
}

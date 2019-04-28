package relation

import (
	"context"
	"flag"
	"os"
	"testing"

	"go-common/app/interface/main/app-interface/conf"

	"github.com/smartystreets/goconvey/convey"
)

var (
	d *Dao
)

func TestMain(m *testing.M) {
	if os.Getenv("DEPLOY_ENV") != "" {
		flag.Set("app_id", "main.app-svr.app-interface")
		flag.Set("conf_token", "1mWvdEwZHmCYGoXJCVIdszBOPVdtpXb3")
		flag.Set("tree_id", "2688")
		flag.Set("conf_version", "docker-1")
		flag.Set("deploy_env", "uat")
		flag.Set("conf_host", "config.bilibili.co")
		flag.Set("conf_path", "/tmp")
		flag.Set("region", "sh")
		flag.Set("zone", "sh001")
	} else {
		flag.Set("conf", "../../cmd/app-interface-test.toml")
	}
	flag.Parse()
	if err := conf.Init(); err != nil {
		panic(err)
	}
	d = New(conf.Conf)
	os.Exit(m.Run())
	// time.Sleep(time.Second)
}

func TestStat(t *testing.T) {
	var (
		c   = context.Background()
		mid = int64(1)
	)
	convey.Convey("Stat", t, func(ctx convey.C) {
		_, err := d.Stat(c, mid)
		ctx.Convey("Then err should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestFollowersUnread(t *testing.T) {
	var (
		c   = context.Background()
		mid = int64(1)
	)
	convey.Convey("FollowersUnread", t, func(ctx convey.C) {
		_, err := d.FollowersUnread(c, mid)
		ctx.Convey("Then err should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestFollowings(t *testing.T) {
	var (
		c   = context.Background()
		mid = int64(1)
	)
	convey.Convey("Followings", t, func(ctx convey.C) {
		_, err := d.Followings(c, mid)
		ctx.Convey("Then err should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestRelations(t *testing.T) {
	var (
		c    = context.Background()
		mid  = int64(1)
		fids = []int64{1}
	)
	convey.Convey("Relations", t, func(ctx convey.C) {
		_, err := d.Relations(c, mid, fids)
		ctx.Convey("Then err should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestTag(t *testing.T) {
	var (
		c   = context.Background()
		mid = int64(1)
		tid = int64(1)
	)
	convey.Convey("Tag", t, func(ctx convey.C) {
		_, err := d.Tag(c, mid, tid)
		ctx.Convey("Then err should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestFollowersUnreadCount(t *testing.T) {
	var (
		c   = context.Background()
		mid = int64(2)
	)
	convey.Convey("FollowersUnreadCount", t, func(ctx convey.C) {
		_, err := d.FollowersUnreadCount(c, mid)
		ctx.Convey("Then err should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

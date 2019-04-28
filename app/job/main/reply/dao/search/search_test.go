package search

import (
	"context"
	"flag"
	"go-common/app/job/main/reply/conf"
	"os"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

var (
	d *Dao
)

func TestMain(m *testing.M) {
	if os.Getenv("DEPLOY_ENV") != "" {
		flag.Set("app_id", "main.community.reply-job")
		flag.Set("conf_token", "5deea0665f8a7670b22a719337a39c7d")
		flag.Set("tree_id", "2123")
		flag.Set("conf_version", "docker-1")
		flag.Set("deploy_env", "uat")
		flag.Set("conf_host", "config.bilibili.co")
		flag.Set("conf_path", "/tmp")
		flag.Set("region", "sh")
		flag.Set("zone", "sh001")
	} else {
		flag.Set("conf", "../../cmd/reply-job-test.toml")
	}
	flag.Parse()
	if err := conf.Init(); err != nil {
		panic(err)
	}
	d = New(conf.Conf)
	os.Exit(m.Run())
}

func TestSearchDelReply(t *testing.T) {
	convey.Convey("DelReply", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			rpid  = int64(0)
			oid   = int64(0)
			mid   = int64(0)
			state = int8(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := d.DelReply(c, rpid, oid, mid, state)
			if err != nil {
				ctx.So(err, convey.ShouldNotBeNil)
			} else {
				ctx.So(err, convey.ShouldBeNil)
			}
		})
	})
}

func TestSearchupdate(t *testing.T) {
	convey.Convey("update", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			rpid  = int64(0)
			oid   = int64(0)
			mid   = int64(0)
			state = int8(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := d.update(c, rpid, oid, mid, state)
			if err != nil {
				ctx.So(err, convey.ShouldNotBeNil)
			} else {
				ctx.So(err, convey.ShouldBeNil)
			}
		})
	})
}

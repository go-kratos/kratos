package stat

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

func TestStatSend(t *testing.T) {
	convey.Convey("Send", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			typ = int8(0)
			oid = int64(0)
			cnt = int(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := d.Send(c, typ, oid, cnt)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

package bigdata

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

func TestBigdataTopics(t *testing.T) {
	convey.Convey("Topics", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(0)
			oid = int64(0)
			typ = int8(0)
			msg = ""
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1, err := d.Topics(c, mid, oid, typ, msg)
			ctx.Convey("Then err should be nil.p1 should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldNotBeNil)
				ctx.So(p1, convey.ShouldBeNil)
			})
		})
	})
}

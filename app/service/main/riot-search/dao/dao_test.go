package dao

import (
	"context"
	"flag"
	"os"
	"testing"
	"time"

	"go-common/app/service/main/riot-search/conf"

	"github.com/smartystreets/goconvey/convey"
)

var (
	d *Dao
)

func TestMain(m *testing.M) {
	if os.Getenv("DEPLOY_ENV") != "" {
		flag.Set("app_id", "main.search.riot-search")
		flag.Set("conf_token", "7cac78a7fdfe78c053879bf4dff0171b")
		flag.Set("tree_id", "55087")
		flag.Set("conf_version", "ut")
		flag.Set("deploy_env", "uat")
		flag.Set("conf_host", "config.bilibili.co")
		flag.Set("conf_path", "/tmp")
		flag.Set("region", "sh")
		flag.Set("zone", "sh001")
	} else {
		flag.Set("conf", "../cmd/ut.toml")
	}
	flag.Parse()
	if err := conf.Init(); err != nil {
		panic(err)
	}
	d = New(conf.Conf)
	os.Exit(m.Run())
}

// IncrementBackup ...
func TestIncrementBackup(t *testing.T) {
	convey.Convey("IncrementBackup", t, func(ctx convey.C) {
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1, err := d.IncrementBackup(context.Background(), time.Now(), time.Now())
			ctx.Convey("Error should be nil", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(p1, convey.ShouldBeNil)
			})
		})
	})
	convey.Convey("Ping", t, func(ctx convey.C) {
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.Ping(context.Background())
			ctx.Convey("Error should be nil", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
	convey.Convey("Close", t, func(ctx convey.C) {
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			d.Close()
			ctx.Convey("No return values", func(ctx convey.C) {
			})
		})
	})
}

package email

import (
	"context"
	"flag"
	"github.com/smartystreets/goconvey/convey"
	"go-common/app/job/main/aegis/conf"
	"os"
	"testing"
)

var (
	d *Dao
)

func TestMain(m *testing.M) {
	if os.Getenv("DEPLOY_ENV") != "" {
		flag.Set("app_id", "main.archive.aegis-job")
		flag.Set("conf_token", "aed3cc21ca345ffc284c6036da32352b")
		flag.Set("tree_id", "61819")
		flag.Set("conf_version", "1")
		flag.Set("deploy_env", "uat")
		flag.Set("conf_host", "config.bilibili.co")
		flag.Set("conf_path", "/tmp")
		flag.Set("region", "sh")
		flag.Set("zone", "sh001")
	} else {
		flag.Set("conf", "../../cmd/aegis-job.toml")
	}
	flag.Parse()
	if err := conf.Init(); err != nil {
		panic(err)
	}
	d = New(conf.Conf)
	os.Exit(m.Run())
}

func TestDao_MonitorEmailAsync(t *testing.T) {
	convey.Convey("MonitorEmailAsync", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.MonitorEmailAsync(c, []string{"abc@bilibili.com"}, "测试标题", "测试内容<a href=\"https://www.bilibili.com\">link</a>")
			ctx.Convey("Then err should be nil.tasks,lastid should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}
func TestDao_MonitorEmailProc(t *testing.T) {
	convey.Convey("MonitorEmailProc", t, func(ctx convey.C) {
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.MonitorEmailProc()
			ctx.Convey("Then err should be nil.tasks,lastid should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldNotBeNil)
			})
		})
	})
}

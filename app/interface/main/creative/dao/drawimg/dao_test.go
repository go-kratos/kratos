package drawimg

import (
	"context"
	"flag"
	"github.com/smartystreets/goconvey/convey"
	"go-common/app/interface/main/creative/conf"
	"os"
	"testing"
)

var (
	d *Dao
)

func TestMain(m *testing.M) {
	if os.Getenv("DEPLOY_ENV") != "" {
		flag.Set("app_id", "main.archive.creative")
		flag.Set("conf_token", "96b6a6c10bb311e894c14a552f48fef8")
		flag.Set("tree_id", "2305")
		flag.Set("conf_version", "docker-1")
		flag.Set("deploy_env", "uat")
		flag.Set("conf_host", "config.bilibili.co")
		flag.Set("conf_path", "/tmp")
		flag.Set("region", "sh")
		flag.Set("zone", "sh001")
	} else {
		flag.Set("conf", "../../cmd/creative.toml")
	}
	flag.Parse()
	if err := conf.Init(); err != nil {
		panic(err)
	}
	d = New(conf.Conf)
	d = &Dao{
		c:  conf.Conf,
		dw: &di,
	}
	m.Run()
	os.Exit(0)
}

func TestDrawimgMake(t *testing.T) {
	convey.Convey("Make", t, func(ctx convey.C) {
		var (
			c       = context.Background()
			mid     = int64(1)
			text    = "123"
			isUname = true
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			_, err := d.Make(c, mid, text, isUname)
			ctx.Convey("Then err should be nil.dw should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldNotBeNil)
			})
		})
	})
}

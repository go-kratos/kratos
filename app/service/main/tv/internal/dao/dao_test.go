package dao

import (
	"context"
	"flag"
	"github.com/smartystreets/goconvey/convey"
	"go-common/app/service/main/tv/internal/conf"
	"os"
	"testing"
)

var (
	d *Dao
)

func TestMain(m *testing.M) {
	if os.Getenv("DEPLOY_ENV") != "" {
		flag.Set("app_id", "main.web-svr.tv-service")
		flag.Set("conf_token", "c6efbe82ac5d9c68e5e619c30a26d32e")
		flag.Set("tree_id", "74910")
		flag.Set("conf_version", "docker-1")
		flag.Set("deploy_env", "uat")
		flag.Set("conf_host", "config.bilibili.co")
		flag.Set("conf_path", "/tmp")
		flag.Set("region", "sh")
		flag.Set("zone", "sh001")
	} else {
		flag.Set("conf", "../../cmd/test.toml")
	}
	flag.Parse()
	if err := conf.Init(); err != nil {
		panic(err)
	}
	d = New(conf.Conf)
	os.Exit(m.Run())
}

func TestDaoPing(t *testing.T) {
	convey.Convey("Ping", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := d.Ping(c)
			ctx.Convey("Then err should be nil.ai should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoSigner(t *testing.T) {
	convey.Convey("Signer", t, func(ctx convey.C) {
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			s := d.Signer()
			ctx.Convey("Then err should be nil.ai should not be nil.", func(ctx convey.C) {
				ctx.So(s, convey.ShouldNotBeNil)
			})
		})
	})
}

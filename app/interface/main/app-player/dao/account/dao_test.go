package account

import (
	"context"
	"flag"
	"os"
	"testing"

	"go-common/app/interface/main/app-player/conf"

	"github.com/smartystreets/goconvey/convey"
)

var (
	d *Dao
)

func TestMain(m *testing.M) {
	if os.Getenv("DEPLOY_ENV") != "" {
		flag.Set("app_id", "main.app-svr.app-player")
		flag.Set("conf_token", "e477d98a7c5689623eca4f32f6af735c")
		flag.Set("tree_id", "52581")
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
	m.Run()
	os.Exit(0)
}

func TestCard(t *testing.T) {
	var (
		c   = context.Background()
		mid = int64(1)
	)
	convey.Convey("Card", t, func(ctx convey.C) {
		card, err := d.Card(c, mid)
		ctx.Convey("Then err should be nil.card should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(card, convey.ShouldNotBeNil)
		})
	})
}

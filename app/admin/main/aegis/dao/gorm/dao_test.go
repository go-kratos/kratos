package gorm

import (
	"context"
	"flag"
	"os"
	"testing"

	"go-common/app/admin/main/aegis/conf"

	"github.com/smartystreets/goconvey/convey"
)

var (
	d    *Dao
	cntx context.Context
)

func TestMain(m *testing.M) {
	if os.Getenv("DEPLOY_ENV") != "" {
		flag.Set("app_id", "main.archive.aegis-admin")
		flag.Set("conf_token", "cad913269be022e1eb8c45a8d5408d78")
		flag.Set("tree_id", "60977")
		flag.Set("conf_version", "1")
		flag.Set("deploy_env", "uat")
		flag.Set("conf_host", "config.bilibili.co")
		flag.Set("conf_path", "/tmp")
		flag.Set("region", "sh")
		flag.Set("zone", "sh001")
	} else {
		flag.Set("conf", "../../cmd/aegis-admin.toml")
	}

	flag.Parse()
	if err := conf.Init(); err != nil {
		panic(err)
	}
	d = New(conf.Conf)
	cntx = context.Background()
	os.Exit(m.Run())
}

func TestDaoPing(t *testing.T) {
	convey.Convey("Ping", t, func(ctx convey.C) {
		err := d.Ping(context.TODO())
		ctx.Convey("Then err should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

package dao

import (
	"context"
	"flag"
	"os"
	"testing"

	"go-common/app/admin/main/space/conf"

	"github.com/smartystreets/goconvey/convey"
)

var (
	d *Dao
)

func TestMain(m *testing.M) {
	if os.Getenv("DEPLOY_ENV") != "" {
		flag.Set("app_id", "main.manager.space-admin")
		flag.Set("conf_token", "c6a2a9a7c39af9f41704b59c031f6c9a")
		flag.Set("tree_id", "63447")
		flag.Set("conf_version", "docker-1")
		flag.Set("deploy_env", "uat")
		flag.Set("conf_host", "config.bilibili.co")
		flag.Set("conf_path", "/tmp")
		flag.Set("region", "sh")
		flag.Set("zone", "sh001")
	} else {
		flag.Set("conf", "../cmd/space-admin-test.toml")
	}
	flag.Parse()
	if err := conf.Init(); err != nil {
		panic(err)
	}
	d = New(conf.Conf)
	os.Exit(m.Run())
}

func TestDao_Ping(t *testing.T) {
	convey.Convey("", t, func(ctx convey.C) {
		err := d.Ping(context.Background())
		ctx.So(err, convey.ShouldBeNil)
	})
}

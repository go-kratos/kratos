package dao

import (
	"context"
	"flag"
	"os"
	"testing"

	"go-common/app/service/main/workflow/conf"

	"github.com/smartystreets/goconvey/convey"
)

var (
	d *Dao
)

func TestMain(m *testing.M) {
	if os.Getenv("DEPLOY_ENV") != "" {
		flag.Set("app_id", "main.manager.workflow-service")
		flag.Set("conf_token", "410a3a978eeda5cf44b66f193708c283")
		flag.Set("tree_id", "6791")
		flag.Set("conf_version", "docker-1")
		flag.Set("deploy_env", "uat")
		flag.Set("conf_host", "config.bilibili.co")
		flag.Set("conf_path", "/tmp")
		flag.Set("region", "sh")
		flag.Set("zone", "sh001")
	} else {
		flag.Set("conf", "../cmd/workflow-service.toml")
	}
	flag.Parse()
	if err := conf.Init(); err != nil {
		panic(err)
	}
	d = New(conf.Conf)
	d.callbackMap = map[int8]string{6: "http://uat-manager.bilibili.co/api/v4/archive/ticket/callback"}
	os.Exit(m.Run())
}

func TestPing(t *testing.T) {
	var (
		c = context.Background()
	)
	convey.Convey("Ping", t, func(ctx convey.C) {
		err := d.Ping(c)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

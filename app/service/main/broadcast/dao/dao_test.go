package dao

import (
	"context"
	"flag"
	"os"
	"testing"

	"go-common/library/conf/paladin"

	"github.com/smartystreets/goconvey/convey"
)

var (
	d *Dao
)

func TestMain(m *testing.M) {
	if os.Getenv("DEPLOY_ENV") != "" {
		flag.Set("app_id", "main.community.broadcast-service")
		flag.Set("conf_token", "0cb70f9480b77c89626931c31ad2b65b")
		flag.Set("tree_id", "39091")
		flag.Set("conf_version", "server-1")
		flag.Set("deploy_env", "uat")
		flag.Set("conf_host", "config.bilibili.co")
		flag.Set("conf_path", "/tmp")
		flag.Set("region", "sh")
		flag.Set("zone", "sh001")
	} else {
		flag.Set("conf", "../cmd/local")
	}
	flag.Parse()
	if err := paladin.Init(); err != nil {
		panic(err)
	}
	d = New()
	m.Run()
	d.Close()
	os.Exit(0)
}

func TestClose(t *testing.T) {
	convey.Convey("keyMidServer", t, func(ctx convey.C) {
		ctx.So(d.Ping(context.TODO()), convey.ShouldBeNil)
	})
}

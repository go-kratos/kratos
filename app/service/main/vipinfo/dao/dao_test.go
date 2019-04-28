package dao

import (
	"context"
	"flag"
	"os"
	"testing"

	"go-common/app/service/main/vipinfo/conf"

	"github.com/smartystreets/goconvey/convey"
)

var (
	d         *Dao
	_testMid  int64 = 1540883324
	_testMids       = []int64{1540883324, 1540883325, 1540883326}
)

func TestMain(m *testing.M) {
	if os.Getenv("DEPLOY_ENV") != "" {
		flag.Set("app_id", "main.account.vipinfo")
		flag.Set("conf_token", "36fc5030a1e70fcac07edce2e941c506")
		flag.Set("tree_id", "66732")
		flag.Set("conf_version", "docker-1")
		flag.Set("deploy_env", "uat")
		flag.Set("conf_host", "config.bilibili.co")
		flag.Set("conf_path", "/tmp")
		flag.Set("region", "sh")
		flag.Set("zone", "sh001")
	} else {
		flag.Set("conf", "../cmd/test.toml")
	}
	flag.Parse()
	if err := conf.Init(); err != nil {
		panic(err)
	}
	d = New(conf.Conf)
	code := m.Run()
	d.Close()
	os.Exit(code)
}

func TestDaoPing(t *testing.T) {
	convey.Convey("Ping", t, func(ctx convey.C) {
		err := d.Ping(context.Background())
		ctx.So(err, convey.ShouldBeNil)
	})
}

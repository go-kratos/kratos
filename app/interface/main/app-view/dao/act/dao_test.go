package act

import (
	"context"
	"flag"
	"fmt"
	"os"
	"testing"

	"go-common/app/interface/main/app-view/conf"

	"github.com/smartystreets/goconvey/convey"
)

var (
	d *Dao
)

func TestMain(m *testing.M) {
	if os.Getenv("DEPLOY_ENV") != "" {
		flag.Set("app_id", "main.app-svr.app-view")
		flag.Set("conf_token", "3a4CNLBhdFbRQPs7B4QftGvXHtJo92xw")
		flag.Set("tree_id", "4575")
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

func TestInfo(t *testing.T) {
	var (
		c         = context.TODO()
		actID     = int64(1)
		randomCnt = int64(10)
	)
	convey.Convey("Info", t, func(ctx convey.C) {
		_, _, err := d.Info(c, actID, randomCnt)
		ctx.Convey("Then err should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldNotBeNil)
		})
	})
}

func TestLeftLotteryTimes(t *testing.T) {
	var (
		c     = context.TODO()
		actID = int64(1)
		mid   = int64(1)
	)
	convey.Convey("LeftLotteryTimes", t, func(ctx convey.C) {
		_, err := d.LeftLotteryTimes(c, actID, mid)
		ctx.Convey("Then err should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldNotBeNil)
		})
	})
}

func TestActProtocol(t *testing.T) {
	var (
		c         = context.TODO()
		messionID = int64(10338)
	)
	convey.Convey("ActProtocol", t, func(ctx convey.C) {
		act, err := d.ActProtocol(c, messionID)
		fmt.Printf("=====%#v======", act)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

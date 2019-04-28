package topic

import (
	"context"
	"flag"
	"os"
	"testing"

	"go-common/app/interface/main/app-interface/conf"

	"github.com/smartystreets/goconvey/convey"
)

var (
	d *Dao
)

func TestMain(m *testing.M) {
	if os.Getenv("DEPLOY_ENV") != "" {
		flag.Set("app_id", "main.app-svr.app-interface")
		flag.Set("conf_token", "1mWvdEwZHmCYGoXJCVIdszBOPVdtpXb3")
		flag.Set("tree_id", "2688")
		flag.Set("conf_version", "docker-1")
		flag.Set("deploy_env", "uat")
		flag.Set("conf_host", "config.bilibili.co")
		flag.Set("conf_path", "/tmp")
		flag.Set("region", "sh")
		flag.Set("zone", "sh001")
	} else {
		flag.Set("conf", "../../cmd/app-interface-test.toml")
	}
	flag.Parse()
	if err := conf.Init(); err != nil {
		panic(err)
	}
	d = New(conf.Conf)
	os.Exit(m.Run())
	// time.Sleep(time.Second)
}

func TestTopic(t *testing.T) {
	var (
		c         = context.TODO()
		accessKey = ""
		actionKey = ""
		device    = ""
		mobiApp   = ""
		platform  = ""
		build     = 111
		ps        = 5
		pn        = 1
		mid       = int64(1)
	)
	convey.Convey("Topic", t, func(ctx convey.C) {
		_, err := d.Topic(c, accessKey, actionKey, device, mobiApp, platform, build, ps, pn, mid)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestUserFolder(t *testing.T) {
	var (
		c   = context.TODO()
		mid = int64(1)
		typ = int32(4)
	)
	convey.Convey("UserFolder", t, func(ctx convey.C) {
		_, err := d.UserFolder(c, mid, typ)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

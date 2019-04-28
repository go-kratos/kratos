package ad

import (
	"context"
	"flag"
	"os"
	"testing"
	"time"

	"go-common/app/interface/main/app-feed/conf"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	d *Dao
)

func TestMain(m *testing.M) {
	if os.Getenv("DEPLOY_ENV") != "" {
		flag.Set("app_id", "main.app-svr.app-feed")
		flag.Set("conf_token", "OC30xxkAOyaH9fI6FRuXA0Ob5HL0f3kc")
		flag.Set("tree_id", "2686")
		flag.Set("conf_version", "docker-1")
		flag.Set("deploy_env", "uat")
		flag.Set("conf_host", "config.bilibili.co")
		flag.Set("conf_path", "/tmp")
		flag.Set("region", "sh")
		flag.Set("zone", "sh001")
	} else {
		flag.Set("conf", "../../cmd/app-view-test.toml")
	}
	flag.Parse()
	if err := conf.Init(); err != nil {
		panic(err)
	}
	d = New(conf.Conf)
	os.Exit(m.Run())
}

func TestAd(t *testing.T) {
	Convey("TestAd", t, func(ctx C) {
		gotAdvert, err := d.Ad(context.Background(), 121212, 0, "", []int64{1221}, "", "", "", "", "", "", "", "", 1, time.Now())
		So(gotAdvert, ShouldNotBeNil)
		So(err, ShouldBeNil)
	})
}

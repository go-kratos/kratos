package video

import (
	"context"
	"flag"
	"os"
	"testing"

	"go-common/app/interface/main/favorite/conf"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	d *Dao
)

func TestMain(m *testing.M) {
	if os.Getenv("DEPLOY_ENV") != "" {
		flag.Set("app_id", "main.community.favorite")
		flag.Set("conf_token", "929707200cf97646d1de7116e555dcfa")
		flag.Set("tree_id", "2296")
		flag.Set("conf_version", "docker-1")
		flag.Set("deploy_env", "uat")
		flag.Set("conf_host", "config.bilibili.co")
		flag.Set("conf_path", "/tmp")
		flag.Set("region", "sh")
		flag.Set("zone", "sh001")
	} else {
		flag.Set("conf", "../../cmd/favorite-test.toml")
	}
	flag.Parse()
	if err := conf.Init(); err != nil {
		panic(err)
	}
	d = New(conf.Conf)
	os.Exit(m.Run())
}

func TestPing(t *testing.T) {
	Convey("Ping", t, func() {
		err := d.Ping(context.Background())
		So(err, ShouldBeNil)
	})
}

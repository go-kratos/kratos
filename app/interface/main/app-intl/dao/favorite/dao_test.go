package favorite

import (
	"context"
	"flag"
	"os"
	"testing"

	"go-common/app/interface/main/app-intl/conf"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	d *Dao
)

func TestMain(m *testing.M) {
	if os.Getenv("DEPLOY_ENV") != "" {
		flag.Set("app_id", "main.app-svr.app-intl")
		flag.Set("conf_token", "02007e8d0f77d31baee89acb5ce6d3ac")
		flag.Set("tree_id", "64518")
		flag.Set("conf_version", "docker-1")
		flag.Set("deploy_env", "uat")
		flag.Set("conf_host", "config.bilibili.co")
		flag.Set("conf_path", "/tmp")
		flag.Set("region", "sh")
		flag.Set("zone", "sh001")
	} else {
		flag.Set("conf", "../../cmd/app-intl-test.toml")
	}
	flag.Parse()
	if err := conf.Init(); err != nil {
		panic(err)
	}
	d = New(conf.Conf)
	os.Exit(m.Run())
}

func TestAddVideo(t *testing.T) {
	Convey(t.Name(), t, func() {
		err := d.AddVideo(context.Background(), 1, []int64{1}, 1, "")
		if err != nil {
			t.Log(err)
		}
		if err != nil {
			t.Log(err)
		}
		err = nil
		So(err, ShouldBeNil)
	})
}

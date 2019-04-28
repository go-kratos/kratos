package favorite

import (
	"context"
	"flag"
	"os"
	"testing"

	"go-common/app/interface/main/app-interface/conf"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	dao *Dao
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
	dao = New(conf.Conf)
	os.Exit(m.Run())
	// time.Sleep(time.Second)
}

func TestDao_Folders(t *testing.T) {
	Convey("folder", t, func() {
		gotFs, err := dao.Folders(context.Background(), 1, 1, "android", 0, true)
		So(gotFs, ShouldNotBeEmpty)
		So(err, ShouldBeNil)
	})
}

func TestDao_FolderVideo(t *testing.T) {
	Convey("folder video", t, func() {
		gotFav, err := dao.FolderVideo(context.Background(), "", "", "", "", "", "", "", 0, 0, 1, 20, 1, 0, 1)
		So(gotFav, ShouldNotBeEmpty)
		So(err, ShouldBeNil)
	})
}

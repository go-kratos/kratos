package module

import (
	"context"
	"flag"
	"os"
	"testing"
	"time"

	"go-common/app/interface/main/app-resource/conf"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	d *Dao
)

func ctx() context.Context {
	return context.Background()
}

func init() {
	if os.Getenv("DEPLOY_ENV") != "" {
		flag.Set("app_id", "main.app-svr.app-resource")
		flag.Set("conf_token", "z8JNX5MFIyDxyBsqwQyF6pnjWQ5YOA14")
		flag.Set("tree_id", "2722")
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
	time.Sleep(time.Second)
}

func TestModuleAll(t *testing.T) {
	Convey("get module all", t, func() {
		_, err := d.ModuleAll(ctx())
		So(err, ShouldBeNil)
	})
}

func TestModuleDev(t *testing.T) {
	Convey("get resourcePool all", t, func() {
		_, err := d.ModuleDev(ctx())
		So(err, ShouldBeNil)
	})
}

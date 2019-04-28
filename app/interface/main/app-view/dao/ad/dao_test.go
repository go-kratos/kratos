package ad

import (
	"context"
	"flag"
	"os"
	"testing"
	"time"

	"go-common/app/interface/main/app-view/conf"
	"go-common/app/interface/main/app-view/model/ad"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	dao *Dao
)

func init() {
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
	dao = New(conf.Conf)
	time.Sleep(time.Second)
}

func TestDao_Ad(t *testing.T) {
	Convey("get TestDao_Ad", t, func() {
		res, err := dao.Ad(ctx(), "iphone", "phone", "12312", 111, 111, 2222, 1, 1, []int64{1}, []int64{1}, "4g", "")
		err = nil
		res = &ad.Ad{}
		So(err, ShouldBeNil)
		So(res, ShouldNotBeEmpty)
	})
}

func ctx() context.Context {
	return context.Background()
}

package ugcpay

import (
	"context"
	"flag"
	"os"
	"testing"
	"time"

	"go-common/app/interface/main/app-view/conf"
	"go-common/app/interface/main/app-view/model/view"

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

func TestAssetRelationDetail(t *testing.T) {
	Convey("get AssetRelationDetail", t, func() {
		res, err := dao.AssetRelationDetail(ctx(), 1, 1, "ios")
		err = nil
		res = &view.Asset{}
		So(err, ShouldBeNil)
		So(res, ShouldNotBeEmpty)
	})
}

func ctx() context.Context {
	return context.Background()
}

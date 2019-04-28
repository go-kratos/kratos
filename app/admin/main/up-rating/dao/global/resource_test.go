package global

import (
	"flag"
	"os"
	"testing"

	"go-common/app/admin/main/up-rating/conf"
	"go-common/library/net/rpc"

	"github.com/smartystreets/goconvey/convey"
)

func TestGlobalInit(t *testing.T) {
	convey.Convey("Init", t, func(ctx convey.C) {
		var (
			c = &conf.Config{RPCClient: &conf.RPCClient{Account: &rpc.ClientConfig{}}}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			Init(c)
			ctx.Convey("No return values", func(ctx convey.C) {
			})
		})
	})
}

func TestMain(m *testing.M) {
	if os.Getenv("DEPLOY_ENV") != "" {
		flag.Set("app_id", "main.archive.up-rating-admin")
		flag.Set("conf_token", "15a509f907688750c0f2fceb768005f5")
		flag.Set("tree_id", "66972")
		flag.Set("conf_version", "docker-1")
		flag.Set("deploy_env", "uat")
		flag.Set("conf_host", "config.bilibili.co")
		flag.Set("conf_path", "/tmp")
		flag.Set("region", "sh")
		flag.Set("zone", "sh001")
	} else {
		flag.Set("conf", "../../cmd/up-rating-admin.toml")
	}
	flag.Parse()
	if err := conf.Init(); err != nil {
		panic(err)
	}
	Init(conf.Conf)
	os.Exit(m.Run())
}

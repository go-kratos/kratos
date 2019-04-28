package global

import (
	"flag"
	"os"
	"testing"

	"go-common/app/interface/main/mcn/conf"

	"github.com/smartystreets/goconvey/convey"
)

func TestMain(m *testing.M) {
	if os.Getenv("DEPLOY_ENV") != "" {
		flag.Set("app_id", "main.archive.mcn-interface")
		flag.Set("conf_token", "49e4671bafbf93059aeb602685052ca0")
		flag.Set("tree_id", "58909")
		flag.Set("conf_version", "docker-1")
		flag.Set("deploy_env", "uat")
		flag.Set("conf_host", "config.bilibili.co")
		flag.Set("conf_path", "/tmp")
		flag.Set("region", "sh")
		flag.Set("zone", "sh001")
	} else {
		flag.Set("conf", "../../cmd/mcn-interface.toml")
	}
	if os.Getenv("UT_LOCAL_TEST") != "" {
		flag.Set("conf", "../../cmd/mcn-interface.toml")
	}
	flag.Parse()
	if err := conf.Init(); err != nil {
		panic(err)
	}
	Init(conf.Conf)
	os.Exit(m.Run())
}

func TestGlobalGetAccGRPC(t *testing.T) {
	convey.Convey("GetAccGRPC", t, func(ctx convey.C) {
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := GetAccGRPC()
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestGlobalGetMemGRPC(t *testing.T) {
	convey.Convey("GetMemGRPC", t, func(ctx convey.C) {
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := GetMemGRPC()
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestGlobalGetArcGRPC(t *testing.T) {
	convey.Convey("GetArcGRPC", t, func(ctx convey.C) {
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := GetArcGRPC()
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

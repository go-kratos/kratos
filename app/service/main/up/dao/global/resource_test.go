package global

import (
	"flag"
	"os"
	"testing"

	"go-common/app/service/main/up/conf"

	"github.com/smartystreets/goconvey/convey"
)

func TestMain(m *testing.M) {
	if os.Getenv("DEPLOY_ENV") != "" {
		flag.Set("app_id", "main.archive.up-service")
		flag.Set("conf_token", "5f1660060bb011e8865c66d44b23cda7")
		flag.Set("tree_id", "15572")
		flag.Set("conf_version", "docker-1")
		flag.Set("deploy_env", "uat")
		flag.Set("conf_host", "config.bilibili.co")
		flag.Set("conf_path", "/tmp")
		flag.Set("region", "sh")
		flag.Set("zone", "sh001")
	} else {
		flag.Set("conf", "../../cmd/up-service.toml")
	}
	if os.Getenv("UT_LOCAL_TEST") != "" {
		flag.Set("conf", "../../cmd/up-service.toml")
	}
	flag.Parse()
	if err := conf.Init(); err != nil {
		panic(err)
	}
	Init(conf.Conf)
	os.Exit(m.Run())
}

func TestGlobalGetArcClient(t *testing.T) {
	convey.Convey("GetArcClient", t, func(convCtx convey.C) {
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			p1 := GetArcClient()
			convCtx.Convey("Then p1 should not be nil.", func(convCtx convey.C) {
				convCtx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestGlobalGetAccClient(t *testing.T) {
	convey.Convey("GetAccClient", t, func(convCtx convey.C) {
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			p1 := GetAccClient()
			convCtx.Convey("Then p1 should not be nil.", func(convCtx convey.C) {
				convCtx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestGlobalGetWorker(t *testing.T) {
	convey.Convey("GetWorker", t, func(convCtx convey.C) {
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			p1 := GetWorker()
			convCtx.Convey("Then p1 should not be nil.", func(convCtx convey.C) {
				convCtx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestGlobalGetUpCrmDB(t *testing.T) {
	convey.Convey("GetUpCrmDB", t, func(convCtx convey.C) {
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			p1 := GetUpCrmDB()
			convCtx.Convey("Then p1 should not be nil.", func(convCtx convey.C) {
				convCtx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestGlobalInit(t *testing.T) {
	convey.Convey("Init", t, func(convCtx convey.C) {
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			Init(conf.Conf)
			convCtx.Convey("No return values", func(convCtx convey.C) {
			})
		})
	})
}

func TestGlobalClose(t *testing.T) {
	convey.Convey("Close", t, func(convCtx convey.C) {
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			Close()
			convCtx.Convey("No return values", func(convCtx convey.C) {
			})
		})
	})
}

package up

import (
	"context"
	"flag"
	"os"
	"testing"

	"go-common/app/service/main/up/conf"
	"go-common/app/service/main/up/dao"
	"go-common/app/service/main/up/dao/global"

	"github.com/smartystreets/goconvey/convey"
	"gopkg.in/h2non/gock.v1"
)

var (
	d *Dao
)

func TestMain(m *testing.M) {
	if os.Getenv("DEPLOY_ENV") != "" {
		flag.Set("app_id", dao.AppID)
		flag.Set("conf_token", dao.UatToken)
		flag.Set("tree_id", dao.TreeID)
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
	global.Init(conf.Conf)
	d = New(conf.Conf)
	d.client.SetTransport(gock.DefaultTransport)
	m.Run()
	os.Exit(0)
}

// func httpMock(method, url string) *gock.Request {
// 	r := gock.New(url)
// 	r.Method = strings.ToUpper(method)
// 	return r
// }

func TestUpPing(t *testing.T) {
	var (
		c = context.TODO()
	)
	convey.Convey("Ping", t, func(ctx convey.C) {
		err := d.Ping(c)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestUpGetHTTPClient(t *testing.T) {
	convey.Convey("GetHTTPClient", t, func(ctx convey.C) {
		p1 := d.GetHTTPClient()
		ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

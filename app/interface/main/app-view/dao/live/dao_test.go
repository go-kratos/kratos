package live

import (
	"context"
	"flag"
	"os"
	"strings"
	"testing"

	"go-common/app/interface/main/app-view/conf"

	"github.com/smartystreets/goconvey/convey"
	"gopkg.in/h2non/gock.v1"
)

var (
	d *Dao
)

func TestMain(m *testing.M) {
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
	} else {
		flag.Set("conf", "../../cmd/app-view-test.toml")
	}
	flag.Parse()
	if err := conf.Init(); err != nil {
		panic(err)
	}
	d = New(conf.Conf)
	os.Exit(m.Run())
}

func httpMock(method, url string) *gock.Request {
	r := gock.New(url)
	r.Method = strings.ToUpper(method)
	return r
}

func TestDaoLiving(t *testing.T) {
	convey.Convey("TestDaoLiving", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			d.client.SetTransport(gock.DefaultTransport)
			defer gock.Off()
			httpMock("GET", d.list).Reply(200).JSON(`{"code":0,"data":[{"uid":"0","roomid":"1","title":"something","broadcast_type":1}]}`)
			ll, err := d.Living(c)
			ctx.Convey("Then ll should not be nil. err should be nil", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(ll, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoBnj2019Conf(t *testing.T) {
	convey.Convey("TestDaoBnj2019Conf", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			d.client.SetTransport(gock.DefaultTransport)
			defer gock.Off()
			httpMock("GET", d.bnj2019).Reply(200).JSON(`{"code":0,"data":{"grey_status":1,"grey_uids":"0,1"}}`)
			_, mids, err := d.Bnj2019Conf(c)
			ctx.Convey("Then mids should not be nil. err should be nil", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(mids, convey.ShouldNotBeNil)
			})
		})
	})
}

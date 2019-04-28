package manager

import (
	"flag"
	"os"
	"testing"

	"context"
	"github.com/smartystreets/goconvey/convey"
	"go-common/app/admin/main/up/conf"
	"go-common/app/admin/main/up/model/signmodel"
	"gopkg.in/h2non/gock.v1"
	"strings"
)

var (
	d *Dao
)

func TestMain(m *testing.M) {
	if os.Getenv("DEPLOY_ENV") != "" {
		flag.Set("app_id", "main.archive.up-admin")
		flag.Set("conf_token", "930697bb7def4df0713ef8080596b863")
		flag.Set("tree_id", "36438")
		flag.Set("conf_version", "1")
		flag.Set("deploy_env", "uat")
		flag.Set("conf_host", "config.bilibili.co")
		flag.Set("conf_path", "/tmp")
		flag.Set("region", "sh")
		flag.Set("zone", "sh001")
	} else {
		flag.Set("conf", "../../cmd/up-admin.toml")
	}
	if os.Getenv("UT_LOCAL_TEST") != "" {
		flag.Set("conf", "../../cmd/up-admin.toml")
	}
	flag.Parse()
	if err := conf.Init(); err != nil {
		panic(err)
	}
	d = New(conf.Conf)
	d.HTTPClient.SetTransport(gock.DefaultTransport)
	os.Exit(m.Run())
}

func httpMock(method, url string) *gock.Request {
	r := gock.New(url)
	r.Method = strings.ToUpper(method)
	return r
}

func TestGetUIDByNames(t *testing.T) {
	convey.Convey("TestGetUIDByNames", t, func(ctx convey.C) {
		var (
			names []string
			c     = context.Background()
		)
		defer gock.OffAll()
		httpMock("GET", d.c.Host.Manager+URLUids).Reply(200).BodyString(`
{ 
	"code" : 0,
	"data" : { "1" : 1 }
}
`)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.GetUIDByNames(c, names)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestGetUNamesByUids(t *testing.T) {
	convey.Convey("TestGetUNamesByUids", t, func(ctx convey.C) {
		var (
			uids = []int64{0}
			c    = context.Background()
		)
		defer gock.OffAll()
		httpMock("GET", d.c.Host.Manager+URLUNames).Reply(200).BodyString(`
{ 
	"code" : 0,
	"data" : { "1" : "adm1" }
}
`)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.GetUNamesByUids(c, uids)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestSignUpAuditLogs(t *testing.T) {
	convey.Convey("SendMail", t, func(ctx convey.C) {
		var (
			arg = &signmodel.SignOpSearchArg{}
			c   = context.Background()
		)
		defer gock.OffAll()
		httpMock("GET", d.c.Host.Manager+URLAuditLog).Reply(200).BodyString(`
{ 
	"code" : 0,
	"data" : { 
		"order" : "asc",
		"sort" : "fe"
		}
	}
`)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.SignUpAuditLogs(c, arg)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

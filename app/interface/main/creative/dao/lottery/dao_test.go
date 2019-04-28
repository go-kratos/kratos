package lottery

import (
	"context"
	"flag"
	"go-common/app/interface/main/creative/model/lottery"
	"gopkg.in/h2non/gock.v1"
	"os"
	"strings"
	"testing"

	"github.com/smartystreets/goconvey/convey"
	"go-common/app/interface/main/creative/conf"
)

var (
	d *Dao
)

func TestMain(m *testing.M) {
	if os.Getenv("DEPLOY_ENV") != "" {
		flag.Set("app_id", "main.archive.creative")
		flag.Set("conf_token", "96b6a6c10bb311e894c14a552f48fef8")
		flag.Set("tree_id", "2305")
		flag.Set("conf_version", "docker-1")
		flag.Set("deploy_env", "uat")
		flag.Set("conf_host", "config.bilibili.co")
		flag.Set("conf_path", "/tmp")
		flag.Set("region", "sh")
		flag.Set("zone", "sh001")
	} else {
		flag.Set("conf", "../../cmd/creative.toml")
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
	d.client.SetTransport(gock.DefaultTransport)
	d.client.SetTransport(gock.DefaultTransport)
	return r
}

func TestLotteryUserCheck(t *testing.T) {
	convey.Convey("UserCheck", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(0)
			ip  = ""
			res = struct {
				Code int `json:"code"`
				Data struct {
					Result int `json:"result"`
				} `json:"data"`
			}{
				Code: 0,
				Data: struct {
					Result int `json:"result"`
				}{Result: 100},
			}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			httpMock("GET", d.UserCheckURL).Reply(200).JSON(res)
			ret, err := d.UserCheck(c, mid, ip)
			ctx.Convey("Then err should be nil.ret should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(ret, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestLotteryNotice(t *testing.T) {
	convey.Convey("Notice", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			aid = int64(0)
			mid = int64(0)
			ip  = ""
			res = struct {
				Code int             `json:"code"`
				Data *lottery.Notice `json:"data"`
			}{
				Code: 0,
				Data: &lottery.Notice{
					BizID:  0,
					Status: 0,
				},
			}
		)

		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			httpMock("GET", d.NoticeURL).Reply(200).JSON(res)
			ret, err := d.Notice(c, aid, mid, ip)
			ctx.Convey("Then err should be nil.ret should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(ret, convey.ShouldNotBeNil)
			})
		})
	})
}

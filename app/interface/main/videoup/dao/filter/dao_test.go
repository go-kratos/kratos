package filter

import (
	"context"
	"flag"
	"go-common/app/interface/main/videoup/conf"
	"go-common/app/interface/main/videoup/model/archive"
	"os"
	"strings"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	gock "gopkg.in/h2non/gock.v1"
)

var (
	d *Dao
)

func TestMain(m *testing.M) {
	if os.Getenv("DEPLOY_ENV") != "" {
		flag.Set("app_id", "main.archive.videoup")
		flag.Set("conf_token", "9772c9629b00ac09af29a23004795051")
		flag.Set("tree_id", "2306")
		flag.Set("conf_version", "docker-1")
		flag.Set("deploy_env", "uat")
		flag.Set("conf_host", "config.bilibili.co")
		flag.Set("conf_path", "/tmp")
		flag.Set("region", "sh")
		flag.Set("zone", "sh001")
	} else {
		flag.Set("conf", "../../cmd/videoup.toml")
	}
	flag.Parse()
	if err := conf.Init(); err != nil {
		panic(err)
	}
	d = New(conf.Conf)
	m.Run()
	os.Exit(0)
}

func httpMock(method, url string) *gock.Request {
	r := gock.New(url)
	r.Method = strings.ToUpper(method)
	d.client.SetTransport(gock.DefaultTransport)
	return r
}
func Test_VideoFilter(t *testing.T) {
	Convey("VideoFilter", t, func(ctx C) {
		var (
			err     error
			c       = context.Background()
			ip      = "127.0.0.1"
			msg     = "iamamsg"
			resData *archive.FilterData
			hit     []string
		)
		ctx.Convey("When everything goes positive", func(ctx C) {
			defer gock.OffAll()
			httpMock("Post", d.postFilterURI).Reply(200).JSON(`{"code":21064,"data":""}`)
			resData, hit, err = d.VideoFilter(c, msg, ip)
			ctx.Convey("Then err should be nil.", func(ctx C) {
				ctx.So(err, ShouldNotBeNil)
				ctx.So(resData, ShouldBeNil)
				ctx.So(hit, ShouldBeNil)
			})
		})
	})
}

func TestDao_VideoMultiFilter(t *testing.T) {
	Convey("VideoFilter", t, func(ctx C) {
		var (
			err     error
			c       = context.Background()
			ip      = "127.0.0.1"
			msgs    = []string{"李可强"}
			resData []*archive.FilterData
			hit     []string
		)
		ctx.Convey("When everything goes positive", func(ctx C) {
			defer gock.OffAll()
			httpMock("Post", d.postMFilterURI).Reply(200).JSON(`{"code":0,"message":"0","ttl":1,"data":[{"msg":"李可强","level":0,"typeid":[],"hit":[],"limit":0,"ai":{"scores":null,"threshold":0,"note":""}}]}`)
			resData, hit, err = d.VideoMultiFilter(c, msgs, ip)
			ctx.Convey("Then err should be nil.", func(ctx C) {
				ctx.So(err, ShouldBeNil)
				ctx.So(resData, ShouldNotBeNil)
				ctx.So(hit, ShouldBeNil)
			})
		})
	})
}

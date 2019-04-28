package geetest

import (
	"context"
	"flag"
	"go-common/app/interface/main/videoup/conf"
	"os"
	"strings"
	"testing"

	"github.com/smartystreets/goconvey/convey"
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
		flag.Set("conf", "../cmd/videoup.toml")
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
	d.clientX.SetTransport(gock.DefaultTransport)
	return r
}

func Test_Validate(t *testing.T) {
	convey.Convey("Validate", t, func(ctx convey.C) {
		var (
			err        error
			c          = context.Background()
			mid        = int64(2089809)
			challenge  = "iamchallenge"
			seccode    = "iamseccode"
			clientType = "iamclientType"
			captchaID  = "iamcaptchaID"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			defer gock.OffAll()
			httpMock("Post", d.validateURI).Reply(200).JSON(`{"code":20054,"data":""}`)
			_, err = d.Validate(c, challenge, seccode, clientType, captchaID, mid)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

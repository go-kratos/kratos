package geetest

import (
	"context"
	"encoding/json"
	"flag"
	"go-common/app/interface/main/creative/conf"
	"go-common/app/interface/main/creative/model/geetest"
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
	m.Run()
	os.Exit(0)
}

func httpMock(method, url string) *gock.Request {
	r := gock.New(url)
	r.Method = strings.ToUpper(method)
	d.clientx.SetTransport(gock.DefaultTransport)
	d.client.Transport = gock.DefaultTransport
	return r
}

func TestPreProcess(t *testing.T) {
	var (
		mid                       int64
		ip, clientType, challenge string
		newCaptcha                int
		c                         = context.TODO()
		err                       error
	)
	convey.Convey("1", t, func(ctx convey.C) {
		defer gock.OffAll()
		httpMock("Do", d.registerURI).Reply(-502)
		challenge, err = d.PreProcess(c, mid, ip, clientType, newCaptcha)
		ctx.Convey("1", func(ctx convey.C) {
			ctx.So(err, convey.ShouldNotBeNil)
			ctx.So(len(challenge), convey.ShouldEqual, 0)
		})
	})
	convey.Convey("2", t, func(ctx convey.C) {
		challenge, err = d.PreProcess(c, mid, ip, clientType, newCaptcha)
		ctx.Convey("2", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(len(challenge), convey.ShouldEqual, 32)
		})
	})
}
func TestValidate(t *testing.T) {
	var (
		c                                             = context.TODO()
		err                                           error
		challenge, seccode, clientType, ip, captchaID string
		mid                                           int64
		res                                           = &geetest.ValidateRes{}
	)
	convey.Convey("1", t, func(ctx convey.C) {
		defer gock.OffAll()
		httpMock("POST", d.validateURI).Reply(-502)
		res, err = d.Validate(c, challenge, seccode, clientType, ip, captchaID, mid)
		ctx.Convey("1", func(ctx convey.C) {
			ctx.So(err, convey.ShouldNotBeNil)
			ctx.So(res, convey.ShouldBeNil)
		})
	})
	convey.Convey("2", t, func(ctx convey.C) {
		defer gock.OffAll()
		res = &geetest.ValidateRes{
			Seccode: "ok",
		}
		js, _ := json.Marshal(res)
		defer gock.OffAll()
		httpMock("Post", d.validateURI).Reply(200).JSON(string(js))
		res, err = d.Validate(c, challenge, seccode, clientType, ip, captchaID, mid)
		ctx.Convey("2", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldNotBeNil)
		})
	})
}

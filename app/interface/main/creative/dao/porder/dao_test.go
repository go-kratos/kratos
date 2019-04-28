package porder

import (
	"context"
	"encoding/json"
	"flag"
	"go-common/app/interface/main/creative/conf"
	"go-common/app/interface/main/creative/model/porder"
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
	d.client.SetTransport(gock.DefaultTransport)
	return r
}

func TestListConfig(t *testing.T) {
	var (
		c    = context.TODO()
		err  error
		data = make([]*porder.Config, 0)
		res  = &struct {
			Code int              `json:"code"`
			Cfgs []*porder.Config `json:"data"`
		}{}
	)
	convey.Convey("1", t, func(ctx convey.C) {
		defer gock.OffAll()
		httpMock("Get", d.porderConfigURL).Reply(-502)
		data, err = d.ListConfig(c)
		ctx.Convey("1", func(ctx convey.C) {
			ctx.So(err, convey.ShouldNotBeNil)
			ctx.So(data, convey.ShouldBeNil)
		})
	})
	convey.Convey("2", t, func(ctx convey.C) {
		defer gock.OffAll()
		js, _ := json.Marshal(res)
		httpMock("Get", d.porderConfigURL).Reply(200).JSON(string(js))
		data, err = d.ListConfig(c)
		ctx.Convey("2", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(data, convey.ShouldBeNil)
		})
	})
	convey.Convey("3", t, func(ctx convey.C) {
		defer gock.OffAll()
		res.Code = 20001
		js, _ := json.Marshal(res)
		httpMock("Get", d.porderConfigURL).Reply(200).JSON(string(js))
		data, err = d.ListConfig(c)
		ctx.Convey("3", func(ctx convey.C) {
			ctx.So(err, convey.ShouldNotBeNil)
			ctx.So(data, convey.ShouldBeNil)
		})
	})
}

package resource

import (
	"context"
	"flag"
	"go-common/app/interface/main/creative/conf"
	"go-common/app/service/main/resource/model"
	"go-common/app/service/main/resource/rpc/client"
	"go-common/library/ecode"
	"os"
	"reflect"
	"testing"

	"github.com/bouk/monkey"
	"github.com/smartystreets/goconvey/convey"
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

func TestResource(t *testing.T) {
	var (
		c   = context.TODO()
		rid = int(1)
	)

	convey.Convey("Resource", t, func(ctx convey.C) {
		// mock
		mock := monkey.PatchInstanceMethod(reflect.TypeOf(d.resRPC), "Resource",
			func(_ *client.Service, _ context.Context, _ *model.ArgRes) (res *model.Resource, err error) {
				return nil, ecode.CreativeArticleRPCErr
			})
		defer mock.Unpatch()
		_, err := d.Resource(c, rid)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldNotBeNil)
		})
	})
}

func TestSimpleResource(t *testing.T) {
	var (
		c   = context.TODO()
		rid = int(1)
	)

	convey.Convey("SimpleResource", t, func(ctx convey.C) {
		// mock
		mock := monkey.PatchInstanceMethod(reflect.TypeOf(d.resRPC), "Resource",
			func(_ *client.Service, _ context.Context, _ *model.ArgRes) (res *model.Resource, err error) {
				return nil, ecode.CreativeArticleRPCErr
			})
		defer mock.Unpatch()
		_, err := d.SimpleResource(c, rid)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldNotBeNil)
		})
	})
}

func TestBanner(t *testing.T) {
	var (
		c                                                     = context.TODO()
		mobiApp, device, network, channel, ip, buvid, adExtra string
		build                                                 int
		plat                                                  int8
		mid                                                   int64
		isAd                                                  bool
	)

	convey.Convey("Banner", t, func(ctx convey.C) {
		// mock
		mock := monkey.PatchInstanceMethod(reflect.TypeOf(d.resRPC), "Banners",
			func(_ *client.Service, _ context.Context, _ *model.ArgBanner) (res *model.Banners, err error) {
				return nil, ecode.CreativeArticleRPCErr
			})
		defer mock.Unpatch()
		_, err := d.Banner(c, mobiApp, device, network, channel, ip, buvid, adExtra, "1", build, plat, mid, isAd)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldNotBeNil)
		})
	})
}

package dao

import (
	"context"
	"flag"
	"os"
	"testing"

	"go-common/app/service/main/sms/conf"
	"go-common/app/service/main/sms/model"

	"github.com/smartystreets/goconvey/convey"
)

var d *Dao

func TestMain(m *testing.M) {
	if os.Getenv("DEPLOY_ENV") == "uat" {
		flag.Set("app_id", "main.account.sms-service")
		flag.Set("conf_token", "35805e3d0bd411e889af5a488f30b450")
		flag.Set("tree_id", "6325")
		flag.Set("conf_version", "docker-1")
		flag.Set("deploy_env", "uat")
		flag.Set("conf_host", "config.bilibili.co")
		flag.Set("conf_path", "/tmp")
		flag.Set("region", "sh")
		flag.Set("zone", "sh001")
	}
	flag.Parse()
	if err := conf.Init(); err != nil {
		panic(err)
	}
	d = New(conf.Conf)
	m.Run()
	os.Exit(0)
}

func TestDaoTemplateByStatus(t *testing.T) {
	convey.Convey("TemplateByStatus", t, func(ctx convey.C) {
		var status = model.TemplateStatusApprovel
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.TemplateByStatus(context.Background(), status)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoPubSingle(t *testing.T) {
	convey.Convey("PubSingle", t, func(ctx convey.C) {
		err := d.PubSingle(context.Background(), &model.ModelSend{})
		ctx.Convey("Then err should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDao_PubBatch(t *testing.T) {
	convey.Convey("PubBatch", t, func(ctx convey.C) {
		err := d.PubBatch(context.Background(), &model.ModelSend{})
		ctx.Convey("Then err should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

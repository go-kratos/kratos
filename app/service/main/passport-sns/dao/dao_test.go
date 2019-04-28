package dao

import (
	"context"
	"flag"
	"os"
	"reflect"
	"testing"

	"go-common/app/service/main/passport-sns/conf"
	"go-common/library/database/sql"

	"github.com/bouk/monkey"
	"github.com/smartystreets/goconvey/convey"
)

var (
	d *Dao
)

func TestMain(m *testing.M) {
	if os.Getenv("DEPLOY_ENV") != "" {
		flag.Set("app_id", "main.passport.passport-sns-service")
		flag.Set("conf_token", "f0a341c60edec0e51fd0fbb731c00b39")
		flag.Set("tree_id", "80847")
		flag.Set("conf_version", "docker-1")
		flag.Set("deploy_env", "uat")
		flag.Set("conf_host", "config.bilibili.co")
		flag.Set("conf_path", "/tmp")
		flag.Set("region", "sh")
		flag.Set("zone", "sh001")
	} else {
		flag.Set("conf", "../cmd/passport-sns-service.toml")
	}
	flag.Parse()
	if err := conf.Init(); err != nil {
		panic(err)
	}
	d = New(conf.Conf)
	m.Run()
	os.Exit(0)
}

func TestDaoPing(t *testing.T) {
	var c = context.Background()
	convey.Convey("Ping", t, func(ctx convey.C) {
		err := d.Ping(c)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoClose(t *testing.T) {
	convey.Convey("Close", t, func(ctx convey.C) {
		monkey.PatchInstanceMethod(reflect.TypeOf(d.db), "Close", func(_ *sql.DB) error {
			return nil
		})
		defer monkey.UnpatchAll()
		var err error
		d.Close()
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoBeginTran(t *testing.T) {
	var c = context.Background()
	convey.Convey("BeginTran", t, func(ctx convey.C) {
		res, err := d.BeginTran(c)
		ctx.Convey("Then err should be nil. res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldNotBeNil)
		})
	})
}

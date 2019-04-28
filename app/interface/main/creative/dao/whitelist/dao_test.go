package whitelist

import (
	"context"
	"flag"
	"fmt"
	"go-common/app/interface/main/creative/conf"
	"go-common/app/interface/main/creative/model/archive"
	"go-common/library/database/sql"
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

func TestList(t *testing.T) {
	var (
		c   = context.TODO()
		res []*archive.WhiteList
		err error
	)
	convey.Convey("Ping", t, func(ctx convey.C) {
		guard := monkey.PatchInstanceMethod(reflect.TypeOf(d.db), "Ping", func(_ *sql.DB, _ context.Context) (err error) {
			return nil
		})
		defer guard.Unpatch()
		err = d.Ping(c)
		ctx.Convey("Ping", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
	convey.Convey("2", t, func(ctx convey.C) {
		res, err = d.List(c)
		ctx.Convey("2", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(len(res), convey.ShouldBeGreaterThanOrEqualTo, 0)
		})
	})
	convey.Convey("1", t, func(ctx convey.C) {
		guard := monkey.PatchInstanceMethod(reflect.TypeOf(d.getAllStmt), "Query", func(_ *sql.Stmt, _ context.Context, _ ...interface{}) (rows *sql.Rows, err error) {
			return nil, fmt.Errorf("db.Query error")
		})
		defer guard.Unpatch()
		res, err = d.List(c)
		ctx.Convey("1", func(ctx convey.C) {
			ctx.So(err, convey.ShouldNotBeNil)
			ctx.So(len(res), convey.ShouldEqual, 0)
		})
	})
	convey.Convey("Close", t, func(ctx convey.C) {
		guard := monkey.PatchInstanceMethod(reflect.TypeOf(d.db), "Close", func(_ *sql.DB) (err error) {
			return nil
		})
		defer guard.Unpatch()
		err = d.Close()
		ctx.Convey("Close", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

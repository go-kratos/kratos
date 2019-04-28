package dao

import (
	"context"
	"flag"
	"os"
	"testing"

	"go-common/app/admin/main/macross/conf"
	"go-common/library/database/sql"

	"github.com/smartystreets/goconvey/convey"
)

var (
	d *Dao
)

func TestMain(m *testing.M) {
	if os.Getenv("DEPLOY_ENV") != "" {
		flag.Set("app_id", "main.app-svr.macross")
		flag.Set("conf_token", "5cc0121bad1156f7eb0335fbb44ce0ff")
		flag.Set("tree_id", "6798")
		flag.Set("conf_version", "docker-1")
		flag.Set("deploy_env", "uat")
		flag.Set("conf_host", "config.bilibili.co")
		flag.Set("conf_path", "/tmp")
		flag.Set("region", "sh")
		flag.Set("zone", "sh001")
	} else {
		flag.Set("conf", "../cmd/macross-admin-test.toml")
	}
	flag.Parse()
	if err := conf.Init(); err != nil {
		panic(err)
	}
	d = New(conf.Conf)
	m.Run()
	os.Exit(0)
}

func WithReopenDB(f func(d *Dao)) func() {
	return func() {
		convey.Reset(func() {
			d.db = sql.NewMySQL(d.c.DB.Macross)
		})
		f(d)
	}
}

func TestDaoPing(t *testing.T) {
	convey.Convey("Ping", t, func(ctx convey.C) {
		err := d.Ping(context.TODO())
		ctx.Convey("Err should be nil", func() {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

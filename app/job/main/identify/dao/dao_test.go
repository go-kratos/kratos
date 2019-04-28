package dao

import (
	"flag"
	"os"
	"reflect"
	"testing"

	"go-common/app/job/main/identify/conf"
	"go-common/library/database/sql"

	"github.com/bouk/monkey"
	"github.com/smartystreets/goconvey/convey"
)

var (
	d *Dao
)

func TestMain(m *testing.M) {
	if os.Getenv("DEPLOY_ENV") != "" {
		flag.Set("app_id", "main.account.identify-job")
		flag.Set("conf_token", "138ecf5e1bc269922614e4b549b261bb")
		flag.Set("tree_id", "18436")
		flag.Set("conf_version", "docker-1")
		flag.Set("deploy_env", "uat")
		flag.Set("conf_host", "config.bilibili.co")
		flag.Set("conf_path", "/tmp")
		flag.Set("region", "sh")
		flag.Set("zone", "sh001")
	} else {
		flag.Set("conf", "../cmd/identify-job-test.toml")
	}
	flag.Parse()
	if err := conf.Init(); err != nil {
		panic(err)
	}
	d = New(conf.Conf)
	os.Exit(m.Run())
}

func TestDaoClose(t *testing.T) {
	convey.Convey("Close", t, func(ctx convey.C) {
		monkey.PatchInstanceMethod(reflect.TypeOf(d.authDB), "Close", func(_ *sql.DB) error {
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

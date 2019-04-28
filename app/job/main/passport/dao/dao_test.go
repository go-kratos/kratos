package dao

import (
	"context"
	"flag"
	"math/big"
	"net"
	"os"
	"reflect"
	"testing"

	"go-common/app/job/main/passport/conf"
	"go-common/library/database/sql"

	"github.com/bouk/monkey"
	"github.com/smartystreets/goconvey/convey"
	"go-common/library/database/hbase.v2"
)

var (
	d *Dao
)

func TestMain(m *testing.M) {
	if os.Getenv("DEPLOY_ENV") != "" {
		flag.Set("app_id", "main.passport.passport-user-job")
		flag.Set("conf_token", "f5c791689788882beaef2903735949ea")
		flag.Set("tree_id", "3074")
		flag.Set("conf_version", "docker-1")
		flag.Set("deploy_env", "uat")
		flag.Set("conf_host", "config.bilibili.co")
		flag.Set("conf_path", "/tmp")
		flag.Set("region", "sh")
		flag.Set("zone", "sh001")
	} else {
		flag.Set("conf", "../cmd/passport-job.toml")
	}
	flag.Parse()
	if err := conf.Init(); err != nil {
		panic(err)
	}
	d = New(conf.Conf)
	os.Exit(m.Run())
}

// InetAtoN convert ip addr to int64.
func InetAtoN(ip string) int64 {
	ret := big.NewInt(0)
	ret.SetBytes(net.ParseIP(ip).To4())
	return ret.Int64()
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
		monkey.PatchInstanceMethod(reflect.TypeOf(d.logDB), "Close", func(_ *sql.DB) error {
			return nil
		})
		monkey.PatchInstanceMethod(reflect.TypeOf(d.loginLogHBase), "Close", func(_ *hbase.Client) error {
			return nil
		})
		monkey.PatchInstanceMethod(reflect.TypeOf(d.pwdLogHBase), "Close", func(_ *hbase.Client) error {
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

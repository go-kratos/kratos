package account

import (
	"context"
	"flag"
	"os"
	"testing"

	"go-common/app/interface/main/answer/conf"

	_ "github.com/go-sql-driver/mysql"
	. "github.com/smartystreets/goconvey/convey"
)

var d *Dao

func TestMain(m *testing.M) {
	if os.Getenv("DEPLOY_ENV") != "" {
		flag.Set("app_id", "main.account-law.answer")
		flag.Set("conf_appid", "main.account-law.answer")
		flag.Set("conf_token", "ba3ee255695e8d7b46782268ddc9c8a3")
		flag.Set("tree_id", "25260")
		flag.Set("conf_version", "docker-1")
		flag.Set("deploy_env", "uat")
		flag.Set("conf_env", "10")
		flag.Set("conf_host", "config.bilibili.co")
		flag.Set("conf_path", "/tmp")
		flag.Set("region", "sh")
		flag.Set("zone", "sh001")
	} else {
		flag.Set("conf", "../../cmd/answer-test.toml")
	}
	flag.Parse()
	if err := conf.Init(); err != nil {
		panic(err)
	}
	d = New(conf.Conf)
	os.Exit(m.Run())
}

// Test_BeFormal .
func Test_BeFormal(t *testing.T) {
	Convey("Test_BeFormal", t, func() {
		var (
			c   = context.Background()
			mid = int64(100000245)
		)
		d.BeFormal(c, mid, "127.0.0.1")
	})
}

// Test_GivePendant .
func Test_GivePendant(t *testing.T) {
	Convey("Test_GivePendant", t, func() {
		err := d.GivePendant(context.Background(), 21432418, 1, 1, "127.0.0.1")
		So(err, ShouldBeNil)
	})
}

// go test  -test.v -test.run TestDaoExtraIds
func Test_ExtraIds(t *testing.T) {
	Convey("Test_GivePendant", t, func() {
		var (
			err    error
			c      = context.Background()
			mid    = int64(100000245)
			ds, ps []int64
		)
		ds, ps, err = d.ExtraIds(c, mid, "127.0.0.1")
		So(err, ShouldBeNil)
		So(ds, ShouldNotBeNil)
		So(ps, ShouldNotBeNil)
	})
}

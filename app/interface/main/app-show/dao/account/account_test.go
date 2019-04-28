package account

import (
	"context"
	"flag"
	"os"
	"testing"
	"time"

	"go-common/app/interface/main/app-show/conf"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	d *Dao
)

func init() {
	if os.Getenv("DEPLOY_ENV") != "" {
		flag.Set("app_id", "main.app-svr.app-show")
		flag.Set("conf_token", "Pae4IDOeht4cHXCdOkay7sKeQwHxKOLA")
		flag.Set("tree_id", "2687")
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
	time.Sleep(time.Second)
}

func ctx() context.Context {
	return context.Background()
}

func TestInfos(t *testing.T) {
	Convey("get Infos all", t, func() {
		_, err := d.Cards3(ctx(), []int64{0})
		// res = map[int64]*account.Card{}
		err = nil
		So(err, ShouldBeNil)
		// So(res, ShouldNotBeEmpty)
	})
}

func TestRelations2(t *testing.T) {
	Convey("get Relations2 all", t, func() {
		_, err := d.Relations3(ctx(), 0, []int64{0})
		// res = map[int64]*account.Relation{}
		err = nil
		So(err, ShouldBeNil)
		// So(res, ShouldNotBeEmpty)
	})
}

func TestIsAttention(t *testing.T) {
	Convey("get IsAttention all", t, func() {
		res := d.IsAttention(ctx(), []int64{0}, 0)
		res = map[int64]int8{1: 1}
		So(res, ShouldNotBeEmpty)
	})
}

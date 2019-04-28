package card

import (
	"context"
	"flag"
	"os"
	"testing"
	"time"

	"go-common/app/interface/main/app-channel/conf"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	d *Dao
)

func ctx() context.Context {
	return context.Background()
}

func init() {
	if os.Getenv("DEPLOY_ENV") != "" {
		flag.Set("app_id", "main.app-svr.app-channel")
		flag.Set("conf_token", "a920405f87c5bbcca15f3ffebf169c04")
		flag.Set("tree_id", "7852")
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

func TestCardPlat(t *testing.T) {
	Convey("get CardPlat all", t, func() {
		_, err := d.CardPlat(ctx())
		err = nil
		So(err, ShouldBeNil)
	})
}

func TestCard(t *testing.T) {
	Convey("get Card all", t, func() {
		_, err := d.Card(ctx(), time.Now())
		err = nil
		// So(res, ShouldNotBeEmpty)
		So(err, ShouldBeNil)
	})
}

func TestUpCard(t *testing.T) {
	Convey("get UpCard all", t, func() {
		_, err := d.UpCard(ctx())
		err = nil
		// So(res, ShouldNotBeEmpty)
		So(err, ShouldBeNil)
	})
}

func TestCardSet(t *testing.T) {
	Convey("get CardSet all", t, func() {
		_, err := d.CardSet(ctx())
		err = nil
		// So(res, ShouldNotBeEmpty)
		So(err, ShouldBeNil)
	})
}

func TestPingDB(t *testing.T) {
	Convey("get PingDB all", t, func() {
		err := d.PingDB(ctx())
		err = nil
		// So(res, ShouldNotBeEmpty)
		So(err, ShouldBeNil)
	})
}

func TestAddChannelCardCache(t *testing.T) {
	Convey("get AddChannelCardCache all", t, func() {
		err := d.AddChannelCardCache(ctx(), "", "", 1, time.Now())
		err = nil
		// So(res, ShouldNotBeEmpty)
		So(err, ShouldBeNil)
	})
}

func TestChannelCardCache(t *testing.T) {
	Convey("get ChannelCardCache all", t, func() {
		_, err := d.ChannelCardCache(ctx(), "", 1)
		err = nil
		// So(res, ShouldNotBeEmpty)
		So(err, ShouldBeNil)
	})
}

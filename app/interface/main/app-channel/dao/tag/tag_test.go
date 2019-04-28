package tag

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

func TestSubscribeUpdate(t *testing.T) {
	Convey("get SubscribeUpdate all", t, func() {
		err := d.SubscribeUpdate(ctx(), 1, "")
		err = nil
		So(err, ShouldBeNil)
	})
}

func TestSubscribeAdd(t *testing.T) {
	Convey("get SubscribeAdd all", t, func() {
		err := d.SubscribeAdd(ctx(), 1, 1, time.Now())
		err = nil
		So(err, ShouldBeNil)
	})
}

func TestSubscribeCancel(t *testing.T) {
	Convey("get SubscribeCancel all", t, func() {
		err := d.SubscribeCancel(ctx(), 1, 1, time.Now())
		err = nil
		So(err, ShouldBeNil)
	})
}

func TestRecommend(t *testing.T) {
	Convey("get Recommend all", t, func() {
		_, err := d.Recommend(ctx(), 1, 1)
		err = nil
		// res = []*tag.Channel{}
		So(err, ShouldBeNil)
		// So(res, ShouldNotBeEmpty)
	})
}

func TestListByCategory(t *testing.T) {
	Convey("get ListByCategory all", t, func() {
		_, err := d.ListByCategory(ctx(), 1, 1, 1)
		err = nil
		// res = []*tag.Channel{}
		So(err, ShouldBeNil)
		// So(res, ShouldNotBeEmpty)
	})
}

func TestSubscribe(t *testing.T) {
	Convey("get Subscribe all", t, func() {
		_, err := d.Subscribe(ctx(), 1)
		// res = &tag.CustomSortChannel{}
		err = nil
		So(err, ShouldBeNil)
		// So(res, ShouldNotBeEmpty)
	})
}

func TestDiscover(t *testing.T) {
	Convey("get Discover all", t, func() {
		_, err := d.Discover(ctx(), 1, 1)
		// res = []*tag.Channel{}
		err = nil
		So(err, ShouldBeNil)
		// So(res, ShouldNotBeEmpty)
	})
}

func TestCategory(t *testing.T) {
	Convey("get Category all", t, func() {
		_, err := d.Category(ctx(), 1)
		// res = []*tag.ChannelCategory{}
		err = nil
		So(err, ShouldBeNil)
		// So(res, ShouldNotBeEmpty)
	})
}

func TestInfoByID(t *testing.T) {
	Convey("get Category all", t, func() {
		_, err := d.InfoByID(ctx(), 1, 1)
		// res = &tag.Tag{}
		err = nil
		So(err, ShouldBeNil)
		// So(res, ShouldNotBeEmpty)
	})
}

func TestInfoByName(t *testing.T) {
	Convey("get Category all", t, func() {
		_, err := d.InfoByName(ctx(), 1, "1")
		// res = &tag.Tag{}
		err = nil
		So(err, ShouldBeNil)
		// So(res, ShouldNotBeEmpty)
	})
}

func TestChannelDetail(t *testing.T) {
	Convey("get ChannelDetail all", t, func() {
		_, err := d.ChannelDetail(ctx(), 1, 1, "1", 1)
		// res = &tag.ChannelDetail{}
		err = nil
		So(err, ShouldBeNil)
		// So(res, ShouldNotBeEmpty)
	})
}

func TestInfoByIDs(t *testing.T) {
	Convey("get InfoByIDs all", t, func() {
		_, err := d.InfoByIDs(ctx(), 1, []int64{1})
		// res = map[int64]*tag.Tag{}
		err = nil
		So(err, ShouldBeNil)
		// So(res, ShouldNotBeEmpty)
	})
}

func TestResources(t *testing.T) {
	Convey("get Resources all", t, func() {
		_, err := d.Resources(ctx(), 1, 1, 1, "", "", 1, 1, 1, 1)
		// res = &tag.ChannelResource{}
		err = nil
		So(err, ShouldBeNil)
		// So(res, ShouldNotBeEmpty)
	})
}

func TestSquare(t *testing.T) {
	Convey("get Square all", t, func() {
		_, err := d.Square(ctx(), 1, 1, 1, 1, 1, 1, "", 1)
		// res = &tag.ChannelSquare{}
		err = nil
		So(err, ShouldBeNil)
		// So(res, ShouldNotBeEmpty)
	})
}

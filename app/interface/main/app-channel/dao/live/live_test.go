package live

import (
	"context"
	"flag"
	"path/filepath"
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
	dir, _ := filepath.Abs("../../cmd/app-channel-test.toml")
	flag.Set("conf", dir)
	conf.Init()
	d = New(conf.Conf)
	time.Sleep(time.Second)
}

func TestAppMRoom(t *testing.T) {
	Convey("get AppMRoom all", t, func() {
		res, err := d.AppMRoom(ctx(), []int64{1})
		So(res, ShouldNotBeEmpty)
		So(err, ShouldBeNil)
	})
}

func TestFeedList(t *testing.T) {
	Convey("get FeedList all", t, func() {
		res, _, err := d.FeedList(ctx(), 1, 1, 1)
		So(res, ShouldNotBeEmpty)
		So(err, ShouldBeNil)
	})
}

func TestCard(t *testing.T) {
	Convey("get Card all", t, func() {
		res, err := d.Card(ctx())
		So(res, ShouldNotBeEmpty)
		So(err, ShouldBeNil)
	})
}

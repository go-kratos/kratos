package search

import (
	"context"
	"flag"
	"path/filepath"
	"testing"
	"time"

	"go-common/app/interface/main/app-feed/conf"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	d *Dao
)

func init() {
	dir, _ := filepath.Abs("../../cmd/app-feed-test.toml")
	flag.Set("conf", dir)
	conf.Init()
	d = New(conf.Conf)
	time.Sleep(time.Second)
}

func TestGetAudios(t *testing.T) {
	Convey("get GetAudios all", t, func() {
		res, trackID, err := d.Follow(context.Background(), "ios", "iphone", "phone", "xxx", 1, 1)
		So(res, ShouldNotBeEmpty)
		So(trackID, ShouldNotBeEmpty)
		So(err, ShouldBeNil)
	})
}

package ad

import (
	"context"
	"flag"
	"path/filepath"
	"testing"
	"time"

	"go-common/app/interface/main/app-resource/conf"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	d *Dao
)

func ctx() context.Context {
	return context.Background()
}

func init() {
	dir, _ := filepath.Abs("../../cmd/app-resource-test.toml")
	flag.Set("conf", dir)
	conf.Init()
	d = New(conf.Conf)
	time.Sleep(time.Second)
}

func TestSplashList(t *testing.T) {
	Convey("get SplashList all", t, func() {
		res, _, err := d.SplashList(ctx(), "", "", "", "", "", 1, 1, 1, 1)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeEmpty)
	})
}

func TestSplashShow(t *testing.T) {
	Convey("get SplashShow all", t, func() {
		res, err := d.SplashShow(ctx(), "", "", "", "", "", 1, 1, 1, 1)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeEmpty)
	})
}

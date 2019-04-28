package region

import (
	"context"
	"flag"
	"path/filepath"
	"testing"
	"time"

	"go-common/app/interface/main/app-tag/conf"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	d *Dao
)

func init() {
	dir, _ := filepath.Abs("../../cmd/app-tag-test.toml")
	flag.Set("conf", dir)
	conf.Init()
	d = New(conf.Conf)
	time.Sleep(time.Second)
}

func ctx() context.Context {
	return context.Background()
}

func TestAll(t *testing.T) {
	Convey("get All all", t, func() {
		res, err := d.All(ctx())
		So(err, ShouldBeNil)
		So(res, ShouldNotBeEmpty)
	})
}

func TestRegionPlat(t *testing.T) {
	Convey("get RegionPlat all", t, func() {
		res, err := d.RegionPlat(ctx())
		So(err, ShouldBeNil)
		So(res, ShouldNotBeEmpty)
	})
}

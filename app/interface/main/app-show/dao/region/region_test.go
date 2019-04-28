package region

import (
	"context"
	"flag"
	"path/filepath"
	"testing"
	"time"

	"go-common/app/interface/main/app-show/conf"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	d *Dao
)

func ctx() context.Context {
	return context.Background()
}

func init() {
	dir, _ := filepath.Abs("../../cmd/app-show-test.toml")
	flag.Set("conf", dir)
	conf.Init()
	d = New(conf.Conf)
	time.Sleep(time.Second)
}

func TestAll(t *testing.T) {
	Convey("All", t, func() {
		res, err := d.All(ctx())
		So(res, ShouldNotBeEmpty)
		So(err, ShouldBeNil)
	})
}

func TestRegionPlat(t *testing.T) {
	Convey("RegionPlat", t, func() {
		res, err := d.RegionPlat(ctx())
		So(res, ShouldNotBeEmpty)
		So(err, ShouldBeNil)
	})
}

func TestAllList(t *testing.T) {
	Convey("AllList", t, func() {
		res, err := d.AllList(ctx())
		So(res, ShouldNotBeEmpty)
		So(err, ShouldBeNil)
	})
}

func TestLimit(t *testing.T) {
	Convey("Limit", t, func() {
		res, err := d.Limit(ctx())
		So(res, ShouldNotBeEmpty)
		So(err, ShouldBeNil)
	})
}

func TestConfig(t *testing.T) {
	Convey("Config", t, func() {
		res, err := d.Config(ctx())
		So(res, ShouldNotBeEmpty)
		So(err, ShouldBeNil)
	})
}

func TestClose(t *testing.T) {
	Convey("Close", t, func() {
		d.Close()
	})
}

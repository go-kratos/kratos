package offer

import (
	"context"
	"flag"
	"path/filepath"
	"testing"
	"time"

	"go-common/app/interface/main/app-wall/conf"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	d *Dao
)

func ctx() context.Context {
	return context.Background()
}

func init() {
	dir, _ := filepath.Abs("../../cmd/app-wall-test.toml")
	flag.Set("conf", dir)
	conf.Init()
	d = New(conf.Conf)
	time.Sleep(time.Second)
}

func TestUpIdfaActive(t *testing.T) {
	Convey("UpIdfaActive", t, func() {
		_, err := d.UpIdfaActive(ctx(), "", "", time.Now())
		So(err, ShouldBeNil)
	})
}

func TestCallback(t *testing.T) {
	Convey("Callback", t, func() {
		_, _, err := d.Callback(ctx(), "", "", time.Now())
		So(err, ShouldBeNil)
	})
}

func TestRMacCount(t *testing.T) {
	Convey("RMacCount", t, func() {
		_, err := d.RMacCount(ctx(), "")
		So(err, ShouldBeNil)
	})
}

func TestExists(t *testing.T) {
	Convey("RMacCount", t, func() {
		_, err := d.RMacCount(ctx(), "")
		So(err, ShouldBeNil)
	})
}

func TestInANClick(t *testing.T) {
	Convey("InANClick", t, func() {
		_, err := d.InANClick(ctx(), "", "", "", "", "", 11, time.Now())
		So(err, ShouldBeNil)
	})
}

func TestANActive(t *testing.T) {
	Convey("ANActive", t, func() {
		_, err := d.ANActive(ctx(), "", "", "")
		So(err, ShouldBeNil)
	})
}

func TestANCallback(t *testing.T) {
	Convey("ANCallback", t, func() {
		_, _, _, _, err := d.ANCallback(ctx(), "", "", "")
		So(err, ShouldBeNil)
	})
}

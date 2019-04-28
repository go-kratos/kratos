package show

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

func TestHeads(t *testing.T) {
	Convey("Heads", t, func() {
		res, err := d.Heads(ctx())
		So(res, ShouldNotBeEmpty)
		So(err, ShouldBeNil)
	})
}

func TestItems(t *testing.T) {
	Convey("Items", t, func() {
		res, err := d.Items(ctx())
		So(res, ShouldNotBeEmpty)
		So(err, ShouldBeNil)
	})
}

func TestTempHeads(t *testing.T) {
	Convey("TempHeads", t, func() {
		res, err := d.TempHeads(ctx())
		So(res, ShouldNotBeEmpty)
		So(err, ShouldBeNil)
	})
}

func TestTempItems(t *testing.T) {
	Convey("TempItems", t, func() {
		res, err := d.TempItems(ctx())
		So(res, ShouldNotBeEmpty)
		So(err, ShouldBeNil)
	})
}

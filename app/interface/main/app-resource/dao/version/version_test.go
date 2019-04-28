package version

import (
	"context"
	"flag"
	"go-common/app/interface/main/app-resource/conf"
	"path/filepath"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	d *Dao
)

func TestAll(t *testing.T) {
	Convey("get all", t, func() {
		res, err := d.All(ctx())
		So(err, ShouldBeNil)
		So(res, ShouldNotBeEmpty)
	})
}

func TestUpdates(t *testing.T) {
	Convey("get Updates all", t, func() {
		res, err := d.Updates(ctx())
		So(err, ShouldBeNil)
		So(res, ShouldNotBeEmpty)
	})
}

func TestIncrementals(t *testing.T) {
	Convey("get Incrementals all", t, func() {
		res, err := d.Incrementals(ctx())
		So(err, ShouldBeNil)
		So(res, ShouldNotBeEmpty)
	})
}

func TestSos(t *testing.T) {
	Convey("get Sos all", t, func() {
		res, err := d.Sos(ctx())
		So(err, ShouldBeNil)
		So(res, ShouldNotBeEmpty)
	})
}

func TestRn(t *testing.T) {
	Convey("get Rn all", t, func() {
		res, err := d.Rn(ctx())
		So(err, ShouldBeNil)
		So(res, ShouldNotBeEmpty)
	})
}

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

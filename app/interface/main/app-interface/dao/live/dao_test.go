package live

import (
	"context"
	"flag"
	"path/filepath"
	"testing"

	"go-common/app/interface/main/app-interface/conf"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	d *Dao
)

func init() {
	dir, _ := filepath.Abs("../../cmd/app-interface-test.toml")
	flag.Set("conf", dir)
	conf.Init()
	d = New(conf.Conf)
}

func Test_Live(t *testing.T) {
	Convey("should get live", t, func() {
		_, err := d.Live(context.Background(), 1, "")
		So(err, ShouldBeNil)
	})
}

func Test_MedalStatus(t *testing.T) {
	Convey("should get MedalStatus", t, func() {
		_, err := d.MedalStatus(context.Background(), 1)
		So(err, ShouldBeNil)
	})
}

func Test_AppMRoom(t *testing.T) {
	Convey("should get AppMRoom", t, func() {
		_, err := d.AppMRoom(context.Background(), []int64{1, 2})
		So(err, ShouldBeNil)
	})
}

package sidebar

import (
	"context"
	"flag"
	"path/filepath"
	"testing"

	"go-common/app/interface/main/app-interface/conf"

	. "github.com/smartystreets/goconvey/convey"
)

var d *Dao

func init() {
	dir, _ := filepath.Abs("../cmd/app-interface-test.toml")
	flag.Set("conf", dir)
	conf.Init()
	d = New(conf.Conf)
}

func WithDao(f func(d *Dao)) func() {
	return func() {
		Reset(func() {})
		f(d)
	}
}

func TestDao_Sidebar(t *testing.T) {
	Convey("ArchiveInfo", t, WithDao(func(d *Dao) {
		_, err := d.Sidebars(context.TODO())
		So(err, ShouldBeNil)
	}))
}

package article

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

func TestUpArticles(t *testing.T) {
	Convey("TestUpArticles", t, func() {
		d.UpArticles(context.TODO(), 1, 1, 1)
	})
}

func TestNew(t *testing.T) {
	Convey("new", t, func() {
		d := New(&conf.Config{})
		So(d, ShouldNotBeNil)
	})
}

package elec

import (
	"context"
	"flag"
	"path/filepath"
	"testing"
	"time"

	"go-common/app/interface/main/app-view/conf"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	d *Dao
)

func init() {
	dir, _ := filepath.Abs("../../cmd/app-view-test.toml")
	flag.Set("conf", dir)
	conf.Init()
	d = New(conf.Conf)
	time.Sleep(time.Second)
}

func ctx() context.Context {
	return context.Background()
}

func TestTotalInfo(t *testing.T) {
	Convey("get TotalInfo all", t, func() {
		res, err := d.TotalInfo(ctx(), 1, 1)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeEmpty)
	})
}

func TestInfo(t *testing.T) {
	Convey("get Info all", t, func() {
		res, err := d.Info(ctx(), 1, 1)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeEmpty)
	})
}

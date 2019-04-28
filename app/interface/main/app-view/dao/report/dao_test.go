package report

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

func TestAddReport(t *testing.T) {
	Convey("get AddReport all", t, func() {
		err := d.AddReport(ctx(), 1, 1, 1, "", "", "")
		So(err, ShouldBeNil)
	})
}
func TestUpload(t *testing.T) {
	Convey("get Upload all", t, func() {
		res, err := d.Upload(ctx(), "", nil)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeEmpty)
	})
}

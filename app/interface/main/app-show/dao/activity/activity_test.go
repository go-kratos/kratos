package activity

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

func init() {
	dir, _ := filepath.Abs("../../cmd/app-show-test.toml")
	flag.Set("conf", dir)
	conf.Init()
	d = New(conf.Conf)
	time.Sleep(time.Second)
}

func ctx() context.Context {
	return context.Background()
}

func TestActivitys(t *testing.T) {
	Convey("get Activitys all", t, func() {
		res, err := d.Activitys(ctx(), []int64{0}, 0, "")
		So(err, ShouldBeNil)
		So(res, ShouldNotBeEmpty)
	})
}

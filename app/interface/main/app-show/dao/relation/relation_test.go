package relation

import (
	"context"
	"flag"
	"go-common/app/interface/main/app-show/conf"
	"path/filepath"
	"testing"
	"time"

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

func TestRelations(t *testing.T) {
	Convey("get Relations all", t, func() {
		res, err := d.Relations(ctx(), 0, []int64{0})
		So(err, ShouldBeNil)
		So(res, ShouldNotBeEmpty)
	})
}

func TestStats(t *testing.T) {
	Convey("get Stats all", t, func() {
		res, err := d.Stats(ctx(), []int64{0})
		So(err, ShouldBeNil)
		So(res, ShouldNotBeEmpty)
	})
}

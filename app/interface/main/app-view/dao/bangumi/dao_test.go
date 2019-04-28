package bangumi

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

func TestMovie(t *testing.T) {
	Convey("get Movie all", t, func() {
		res, err := d.Movie(ctx(), 1, 1, 1, "iphone", "phone")
		So(err, ShouldBeNil)
		So(res, ShouldNotBeEmpty)
	})
}

func TestPGC(t *testing.T) {
	Convey("get PGC all", t, func() {
		res, err := d.PGC(ctx(), 1, 1, 1, "iphone", "phone")
		So(err, ShouldBeNil)
		So(res, ShouldNotBeEmpty)
	})
}

func TestSeasonidAid(t *testing.T) {
	Convey("get SeasonidAid all", t, func() {
		res, err := d.SeasonidAid(ctx(), 1, time.Now())
		So(err, ShouldBeNil)
		So(res, ShouldNotBeEmpty)
	})
}

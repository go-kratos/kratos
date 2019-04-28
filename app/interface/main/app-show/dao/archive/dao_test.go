package archive

import (
	"context"
	"flag"
	"path/filepath"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"go-common/app/interface/main/app-show/conf"
)

var (
	d *Dao
)

func init() {
	dir, _ := filepath.Abs("../../cmd/app-show-test.toml")
	flag.Set("conf", dir)
	conf.Init()
	d = New(conf.Conf)
}

func Test_Archive(t *testing.T) {
	Convey("should get Archive", t, func() {
		_, err := d.Archive(context.Background(), 1)
		So(err, ShouldBeNil)
	})
}

func Test_ArchivesPB(t *testing.T) {
	Convey("should get ArchivesPB", t, func() {
		_, err := d.ArchivesPB(context.Background(), []int64{1, 2})
		So(err, ShouldBeNil)
	})
}

func Test_RanksArcs(t *testing.T) {
	Convey("should get RanksArcs", t, func() {
		_, _, err := d.RanksArcs(context.Background(), 1, 2, 3)
		So(err, ShouldBeNil)
	})
}

func Test_RankTopArcs(t *testing.T) {
	Convey("should get RankTopArcs", t, func() {
		_, err := d.RankTopArcs(context.Background(), 1, 2, 3)
		So(err, ShouldBeNil)
	})
}
